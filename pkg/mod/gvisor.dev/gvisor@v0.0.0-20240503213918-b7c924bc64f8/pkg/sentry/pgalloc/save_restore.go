// Copyright 2018 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pgalloc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"runtime"

	"golang.org/x/sys/unix"
	"gvisor.dev/gvisor/pkg/atomicbitops"
	"gvisor.dev/gvisor/pkg/hostarch"
	"gvisor.dev/gvisor/pkg/log"
	"gvisor.dev/gvisor/pkg/sentry/usage"
	"gvisor.dev/gvisor/pkg/state"
	"gvisor.dev/gvisor/pkg/state/statefile"
)

// SaveTo writes f's state to the given stream.
func (f *MemoryFile) SaveTo(ctx context.Context, w io.Writer, pw io.Writer) error {
	// Wait for reclaim.
	f.mu.Lock()
	defer f.mu.Unlock()
	for f.reclaimable {
		f.reclaimCond.Signal()
		f.mu.Unlock()
		runtime.Gosched()
		f.mu.Lock()
	}

	// Ensure that there are no pending evictions.
	if len(f.evictable) != 0 {
		panic(fmt.Sprintf("evictions still pending for %d users; call StartEvictions and WaitForEvictions before SaveTo", len(f.evictable)))
	}

	// Ensure that all pages that contain data have knownCommitted set, since
	// we only store knownCommitted pages below.
	zeroPage := make([]byte, hostarch.PageSize)
	err := f.updateUsageLocked(0, nil, func(bs []byte, committed []byte) error {
		for pgoff := 0; pgoff < len(bs); pgoff += hostarch.PageSize {
			i := pgoff / hostarch.PageSize
			pg := bs[pgoff : pgoff+hostarch.PageSize]
			if !bytes.Equal(pg, zeroPage) {
				committed[i] = 1
				continue
			}
			committed[i] = 0
			// Reading the page caused it to be committed; decommit it to
			// reduce memory usage.
			//
			// "MADV_REMOVE [...] Free up a given range of pages and its
			// associated backing store. This is equivalent to punching a hole
			// in the corresponding byte range of the backing store (see
			// fallocate(2))." - madvise(2)
			if err := unix.Madvise(pg, unix.MADV_REMOVE); err != nil {
				// This doesn't impact the correctness of saved memory, it
				// just means that we're incrementally more likely to OOM.
				// Complain, but don't abort saving.
				log.Warningf("Decommitting page %p while saving failed: %v", pg, err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Save metadata.
	if _, err := state.Save(ctx, w, &f.fileSize); err != nil {
		return err
	}
	if _, err := state.Save(ctx, w, &f.usage); err != nil {
		return err
	}

	// Dump out committed pages.
	for seg := f.usage.FirstSegment(); seg.Ok(); seg = seg.NextSegment() {
		if !seg.Value().knownCommitted {
			continue
		}
		// Write a header to distinguish from objects.
		if err := state.WriteHeader(w, uint64(seg.Range().Length()), false); err != nil {
			return err
		}
		// Write out data.
		var ioErr error
		err := f.forEachMappingSlice(seg.Range(), func(s []byte) {
			if ioErr != nil {
				return
			}
			_, ioErr = pw.Write(s)
		})
		if ioErr != nil {
			return ioErr
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// MarkSavable marks f as savable.
func (f *MemoryFile) MarkSavable() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.savable = true
}

// IsSavable returns true if f is savable.
func (f *MemoryFile) IsSavable() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.savable
}

// RestoreID returns the restore ID for f.
func (f *MemoryFile) RestoreID() string {
	return f.opts.RestoreID
}

// LoadFrom loads MemoryFile state from the given stream.
func (f *MemoryFile) LoadFrom(ctx context.Context, r io.Reader, pr *statefile.AsyncReader) error {
	// Load metadata.
	if _, err := state.Load(ctx, r, &f.fileSize); err != nil {
		return err
	}
	if err := f.file.Truncate(f.fileSize); err != nil {
		return err
	}
	newMappings := make([]uintptr, f.fileSize>>chunkShift)
	f.mappings.Store(&newMappings)
	if _, err := state.Load(ctx, r, &f.usage); err != nil {
		return err
	}

	// Try to map committed chunks concurrently: For any given chunk, either
	// this loop or the following one will mmap the chunk first and cache it in
	// f.mappings for the other, but this loop is likely to run ahead of the
	// other since it doesn't do any work between mmaps. The rest of this
	// function doesn't mutate f.usage, so it's safe to iterate concurrently.
	mapperDone := make(chan struct{})
	mapperCanceled := atomicbitops.FromInt32(0)
	go func() { // S/R-SAFE: see comment
		defer func() { close(mapperDone) }()
		for seg := f.usage.FirstSegment(); seg.Ok(); seg = seg.NextSegment() {
			if mapperCanceled.Load() != 0 {
				return
			}
			if seg.Value().knownCommitted {
				f.forEachMappingSlice(seg.Range(), func(s []byte) {})
			}
		}
	}()
	defer func() {
		mapperCanceled.Store(1)
		<-mapperDone
	}()

	// Load committed pages.
	for seg := f.usage.FirstSegment(); seg.Ok(); seg = seg.NextSegment() {
		if !seg.Value().knownCommitted {
			continue
		}
		// Verify header.
		length, object, err := state.ReadHeader(r)
		if err != nil {
			return err
		}
		if object {
			// Not expected.
			return fmt.Errorf("unexpected object")
		}
		if expected := uint64(seg.Range().Length()); length != expected {
			// Size mismatch.
			return fmt.Errorf("mismatched segment: expected %d, got %d", expected, length)
		}
		// Read data.
		var ioErr error
		err = f.forEachMappingSlice(seg.Range(), func(s []byte) {
			if ioErr != nil {
				return
			}
			if pr != nil {
				pr.ReadAsync(s)
			} else {
				_, ioErr = io.ReadFull(r, s)
			}
		})
		if ioErr != nil {
			return ioErr
		}
		if err != nil {
			return err
		}

		// Update accounting for restored pages. We need to do this here since
		// these segments are marked as "known committed", and will be skipped
		// over on accounting scans.
		usage.MemoryAccounting.Inc(seg.End()-seg.Start(), seg.Value().kind, seg.Value().memCgID)
	}

	return nil
}
