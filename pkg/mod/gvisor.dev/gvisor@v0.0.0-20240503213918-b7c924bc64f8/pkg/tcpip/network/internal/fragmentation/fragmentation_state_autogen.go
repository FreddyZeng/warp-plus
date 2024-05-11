// automatically generated by stateify.

package fragmentation

import (
	"context"

	"gvisor.dev/gvisor/pkg/state"
)

func (l *reassemblerList) StateTypeName() string {
	return "pkg/tcpip/network/internal/fragmentation.reassemblerList"
}

func (l *reassemblerList) StateFields() []string {
	return []string{
		"head",
		"tail",
	}
}

func (l *reassemblerList) beforeSave() {}

// +checklocksignore
func (l *reassemblerList) StateSave(stateSinkObject state.Sink) {
	l.beforeSave()
	stateSinkObject.Save(0, &l.head)
	stateSinkObject.Save(1, &l.tail)
}

func (l *reassemblerList) afterLoad(context.Context) {}

// +checklocksignore
func (l *reassemblerList) StateLoad(ctx context.Context, stateSourceObject state.Source) {
	stateSourceObject.Load(0, &l.head)
	stateSourceObject.Load(1, &l.tail)
}

func (e *reassemblerEntry) StateTypeName() string {
	return "pkg/tcpip/network/internal/fragmentation.reassemblerEntry"
}

func (e *reassemblerEntry) StateFields() []string {
	return []string{
		"next",
		"prev",
	}
}

func (e *reassemblerEntry) beforeSave() {}

// +checklocksignore
func (e *reassemblerEntry) StateSave(stateSinkObject state.Sink) {
	e.beforeSave()
	stateSinkObject.Save(0, &e.next)
	stateSinkObject.Save(1, &e.prev)
}

func (e *reassemblerEntry) afterLoad(context.Context) {}

// +checklocksignore
func (e *reassemblerEntry) StateLoad(ctx context.Context, stateSourceObject state.Source) {
	stateSourceObject.Load(0, &e.next)
	stateSourceObject.Load(1, &e.prev)
}

func init() {
	state.Register((*reassemblerList)(nil))
	state.Register((*reassemblerEntry)(nil))
}