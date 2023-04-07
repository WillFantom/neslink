package neslink

import (
	"errors"
	"fmt"
	"os"
	"path"

	"golang.org/x/sys/unix"
)

// NsProvider offers a approach to obtaining network namespace paths based on
// given conditions.
type NsProvider struct {
	name string
	f    func() (Namespace, error)
}

var (
	errNoNs error = errors.New("failed to obtain netns from provider")
)

// Provide determines the network namespace path based on the provider's
// conditions. Since some conditions are collected at the time of the provider's
// creation and others when this function is called, repeat calls are not always
// expected to produce the same result. Also note, the path is only returned,
// not opened.
func (nsp NsProvider) Provide() (Namespace, error) {
	return nsp.f()
}

// NPGeneric provides the means to create custom providers. See the docker
// provider for an example of this.
func NPGeneric(providerName string, function func() (Namespace, error)) NsProvider {
	if providerName == "" {
		providerName = "unnamed-ns-provider"
	}
	return NsProvider{
		name: providerName,
		f:    function,
	}
}

// NPNow returns a netns provider that provides the netns path for the
// process/thread that calls the Provide function.
func NPNow() NsProvider {
	return NsProvider{
		name: "now",
		f: func() (Namespace, error) {
			return Namespace(fmt.Sprintf("/proc/%d/task/%d/ns/net", os.Getpid(), unix.Gettid())), nil
		},
	}
}

// NPProcess returns a netns provider that provides the netns path for the
// process associated with the given process ID.
func NPProcess(pid int) NsProvider {
	return NsProvider{
		name: "process",
		f: func() (Namespace, error) {
			return Namespace(fmt.Sprintf("/proc/%d/ns/net", pid)), nil
		},
	}
}

// NPThread returns a netns provider that provides the netns path for the
// process associated with the given process and thread ID.
func NPThread(pid, tid int) NsProvider {
	return NsProvider{
		name: "thread",
		f: func() (Namespace, error) {
			return Namespace(fmt.Sprintf("/proc/%d/task/%d/ns/net", pid, tid)), nil
		},
	}
}

// NPName returns a netns provider that provides the netns path for a named
// (mounted) netns. This assumes the ns is mounted in the default location.
func NPName(name string) NsProvider {
	return NsProvider{
		name: "name",
		f: func() (Namespace, error) {
			return Namespace(path.Join(DefaultMountPath, name)), nil
		},
	}
}

// NPNameAt returns a netns provider that provides the netns path for a named
// (mounted) netns.
func NPNameAt(mountdir, name string) NsProvider {
	return NsProvider{
		name: "name-at",
		f: func() (Namespace, error) {
			return Namespace(path.Join(mountdir, name)), nil
		},
	}
}

// NPPath returns a netns provider that provides the netns path based on the
// path given.
func NPPath(path string) NsProvider {
	ns := Namespace(path)
	return NsProvider{
		name: "path",
		f: func() (Namespace, error) {
			return ns, nil
		},
	}
}
