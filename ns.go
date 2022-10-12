package neslink

import (
	"github.com/vishvananda/netns"
	"golang.org/x/sys/unix"
)

// NsFd is a file descriptor referencing a network namespace. This can be used
// to interact and mange the netns.
type NsFd int

// nsHandle casts the NsFd to a netlink package NsHandle for interaction with
// the netlink package.
func (fd NsFd) nsHandle() netns.NsHandle {
	return netns.NsHandle(int(fd))
}

// Fd simply returns the file descriptor as an int.
func (fd NsFd) Fd() int {
	return int(fd)
}

// ID provides a unique ID for the network namespace in string form. This can be
// useful for determining if 2 NsFds reference the same netns.
func (fd NsFd) ID() string {
	return fd.nsHandle().UniqueId()
}

// Set sets the current thread's network namespace to the one associated with
// the netns file descriptor. In most cases, the Do functions in this package
// should be used, so only use this function when the Do functions do not
// suffice.
func (fd NsFd) Set() error {
	return netns.Set(fd.nsHandle())
}

// Close effectively closes the reference to the netns, thus after calling the
// file descriptor should not be used again. This is handled already by the Do
// functions in this package, so is only needed to be used when manually
// managaging netns handles.
func (fd *NsFd) Close() error {
	// code taken from netns package to avoid excessive casting...
	if err := unix.Close(int(*fd)); err != nil {
		return err
	}
	(*fd) = -1
	return nil
}

// DeleteNamedNs literally just called the delete named from the netlink
// package, but is here just incase of a use case where there is no other
// reason to directly import that package.
func DeleteNamedNs(name string) error {
	return netns.DeleteNamed(name)
}
