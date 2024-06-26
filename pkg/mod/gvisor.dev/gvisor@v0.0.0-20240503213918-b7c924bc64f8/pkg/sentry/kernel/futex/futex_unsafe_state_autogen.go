// automatically generated by stateify.

package futex

import (
	"context"

	"gvisor.dev/gvisor/pkg/state"
)

func (p *AtomicPtrBucket) StateTypeName() string {
	return "pkg/sentry/kernel/futex.AtomicPtrBucket"
}

func (p *AtomicPtrBucket) StateFields() []string {
	return []string{
		"ptr",
	}
}

func (p *AtomicPtrBucket) beforeSave() {}

// +checklocksignore
func (p *AtomicPtrBucket) StateSave(stateSinkObject state.Sink) {
	p.beforeSave()
	var ptrValue *bucket
	ptrValue = p.savePtr()
	stateSinkObject.SaveValue(0, ptrValue)
}

func (p *AtomicPtrBucket) afterLoad(context.Context) {}

// +checklocksignore
func (p *AtomicPtrBucket) StateLoad(ctx context.Context, stateSourceObject state.Source) {
	stateSourceObject.LoadValue(0, new(*bucket), func(y any) { p.loadPtr(ctx, y.(*bucket)) })
}

func init() {
	state.Register((*AtomicPtrBucket)(nil))
}
