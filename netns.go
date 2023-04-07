package neslink

import "os"

// Namespace is a path to a file associated with a network namespace.
type Namespace string

// NsFd is a file descriptor for an open Namespace file.
type NsFd int

const (
	NsFdNone         NsFd   = NsFd(-1)
	DefaultMountPath string = "/run/netns"
)

// String returns the Namespace file path as a string.
func (n Namespace) String() string {
	return string(n)
}

// Int returns the Namespace file descriptor as an int.
func (n NsFd) Int() int {
	return int(n)
}

// Valid determines if the file descriptor is valid. This can be used to
// determine if a returned NsFd is ok regardless of the error.
func (n NsFd) Valid() bool {
	return n.Int() > 0
}

// Exists determines if the path used for the namespace exists. Whilst not an
// exhaustive check, this can help debug namespace providers.
func (ns Namespace) Exists() bool {
	if info, err := os.Stat(ns.String()); err != nil {
		if !info.IsDir() {
			return true
		}
	}
	return false
}
