package neslink

import "github.com/vishvananda/netns"

// NsProvider is a function that retruns a network namespace (netns) handle
// (file descriptor) when it is called. If the provider fails to correctly
// obtain the netns, an error is returned.
type NsProvider func() (NsFd, error)

// NPNow returns a netns provider that simply provides the netns fd that was
// found at the time that the provider itself was generated (the time that this
// function is called).
func NPNow() NsProvider {
	ns, err := netns.Get()
	return func() (NsFd, error) {
		return NsFd(ns), err
	}
}

// NPCurrent returns a netns provider that will return the current netns fd as
// of the time when the provider itself is called.
func NPCurrent() NsProvider {
	return func() (NsFd, error) {
		ns, err := netns.Get()
		return NsFd(ns), err
	}
}

// NPName returns a netns provider that when called will return the netns fd
// associated with the given name.
func NPName(name string) NsProvider {
	return func() (NsFd, error) {
		ns, err := netns.GetFromName(name)
		return NsFd(ns), err
	}
}

// NPDocker returns a netns provider that when called will return the netns fd
// associated with a given docker container.
func NPDocker(containerID string) NsProvider {
	return func() (NsFd, error) {
		ns, err := netns.GetFromDocker(containerID)
		return NsFd(ns), err
	}
}

// NPPath returns a netns provider that when called will return the netns fd
// associated with the given file path.
func NPPath(path string) NsProvider {
	return func() (NsFd, error) {
		ns, err := netns.GetFromPath(path)
		return NsFd(ns), err
	}
}

// NPProcess returns a netns provider that when called will return the netns fd
// associated with the given process.
func NPProcess(pid int) NsProvider {
	return func() (NsFd, error) {
		ns, err := netns.GetFromPid(pid)
		return NsFd(ns), err
	}
}

// NPThread returns a netns provider that when called will return the netns fd
// associated with the given thread.
func NPThread(pid, tid int) NsProvider {
	return func() (NsFd, error) {
		ns, err := netns.GetFromThread(pid, tid)
		return NsFd(ns), err
	}
}

// NPNew returns a netns provider that when called will create a new unnamed
// netns (in turn switching to it) and returns the fd of the newly created ns.
func NPNew() NsProvider {
	return func() (NsFd, error) {
		ns, err := netns.New()
		return NsFd(ns), err
	}
}

// NPNewNamed returns a netns provider that when called will create a new named
// netns (in turn switching to it) and returns the fd of the newly created ns.
func NPNewNamed(name string) NsProvider {
	return func() (NsFd, error) {
		ns, err := netns.NewNamed(name)
		return NsFd(ns), err
	}
}
