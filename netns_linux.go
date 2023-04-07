//go:build linux
// +build linux

package neslink

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// close closes the file descriptor. This should be used to clean up any opned
// file descriptor.
func (n NsFd) close() error {
	if err := unix.Close(n.Int()); err != nil {
		return fmt.Errorf("failed to close netns file descriptor %d: %w", n.Int(), err)
	}
	return nil
}

// open opens the file for the network namespace and returns the file
// descriptor.
func (ns Namespace) open() (NsFd, error) {
	fd, err := unix.Open(ns.String(), unix.O_RDONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		return NsFdNone, fmt.Errorf("failed to open namespace: %w", err)
	}
	return NsFd(fd), nil
}

// set sets the current namespace to the one associated with the given file
// descriptor.
func (ns NsFd) set() error {
	return unix.Setns(ns.Int(), unix.CLONE_NEWNET)
}
