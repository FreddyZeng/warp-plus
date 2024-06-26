// automatically generated by stateify.

//go:build !false
// +build !false

package kernel

import (
	"context"

	"gvisor.dev/gvisor/pkg/state"
)

func (s *SpecialOpts) StateTypeName() string {
	return "pkg/sentry/kernel.SpecialOpts"
}

func (s *SpecialOpts) StateFields() []string {
	return []string{}
}

func (s *SpecialOpts) beforeSave() {}

// +checklocksignore
func (s *SpecialOpts) StateSave(stateSinkObject state.Sink) {
	s.beforeSave()
}

func (s *SpecialOpts) afterLoad(context.Context) {}

// +checklocksignore
func (s *SpecialOpts) StateLoad(ctx context.Context, stateSourceObject state.Source) {
}

func init() {
	state.Register((*SpecialOpts)(nil))
}
