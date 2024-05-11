// automatically generated by stateify.

package unix

import (
	"context"

	"gvisor.dev/gvisor/pkg/state"
)

func (r *socketRefs) StateTypeName() string {
	return "pkg/sentry/socket/unix.socketRefs"
}

func (r *socketRefs) StateFields() []string {
	return []string{
		"refCount",
	}
}

func (r *socketRefs) beforeSave() {}

// +checklocksignore
func (r *socketRefs) StateSave(stateSinkObject state.Sink) {
	r.beforeSave()
	stateSinkObject.Save(0, &r.refCount)
}

// +checklocksignore
func (r *socketRefs) StateLoad(ctx context.Context, stateSourceObject state.Source) {
	stateSourceObject.Load(0, &r.refCount)
	stateSourceObject.AfterLoad(func() { r.afterLoad(ctx) })
}

func (s *Socket) StateTypeName() string {
	return "pkg/sentry/socket/unix.Socket"
}

func (s *Socket) StateFields() []string {
	return []string{
		"vfsfd",
		"FileDescriptionDefaultImpl",
		"DentryMetadataFileDescriptionImpl",
		"LockFD",
		"SendReceiveTimeout",
		"socketRefs",
		"namespace",
		"ep",
		"stype",
		"abstractName",
		"abstractBound",
	}
}

func (s *Socket) beforeSave() {}

// +checklocksignore
func (s *Socket) StateSave(stateSinkObject state.Sink) {
	s.beforeSave()
	stateSinkObject.Save(0, &s.vfsfd)
	stateSinkObject.Save(1, &s.FileDescriptionDefaultImpl)
	stateSinkObject.Save(2, &s.DentryMetadataFileDescriptionImpl)
	stateSinkObject.Save(3, &s.LockFD)
	stateSinkObject.Save(4, &s.SendReceiveTimeout)
	stateSinkObject.Save(5, &s.socketRefs)
	stateSinkObject.Save(6, &s.namespace)
	stateSinkObject.Save(7, &s.ep)
	stateSinkObject.Save(8, &s.stype)
	stateSinkObject.Save(9, &s.abstractName)
	stateSinkObject.Save(10, &s.abstractBound)
}

func (s *Socket) afterLoad(context.Context) {}

// +checklocksignore
func (s *Socket) StateLoad(ctx context.Context, stateSourceObject state.Source) {
	stateSourceObject.Load(0, &s.vfsfd)
	stateSourceObject.Load(1, &s.FileDescriptionDefaultImpl)
	stateSourceObject.Load(2, &s.DentryMetadataFileDescriptionImpl)
	stateSourceObject.Load(3, &s.LockFD)
	stateSourceObject.Load(4, &s.SendReceiveTimeout)
	stateSourceObject.Load(5, &s.socketRefs)
	stateSourceObject.Load(6, &s.namespace)
	stateSourceObject.Load(7, &s.ep)
	stateSourceObject.Load(8, &s.stype)
	stateSourceObject.Load(9, &s.abstractName)
	stateSourceObject.Load(10, &s.abstractBound)
}

func init() {
	state.Register((*socketRefs)(nil))
	state.Register((*Socket)(nil))
}
