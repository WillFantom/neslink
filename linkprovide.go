package neslink

import (
	"errors"

	"github.com/vishvananda/netlink"
)

type LinkProvider struct {
	name string
	f    func() (netlink.Link, error)
}

var (
	errNoLink error = errors.New("failed to obtain link from provider")
)

// Provide determines the network namespace path based on the provider's
// conditions. Since some conditions are collected at the time of the provider's
// creation and others when this function is called, repeat calls are not always
// expected to produce the same result. Also note, the path is only returned,
// not opened.
func (lp LinkProvider) Provide() (netlink.Link, error) {
	return lp.f()
}

// LPName creates a link provider that when called, will provide the
// pre-existing link with the given name (in the namespace this is called in).
// If no matches are found, an error is returned.
func LPName(name string) LinkProvider {
	return LinkProvider{
		name: "name",
		f: func() (netlink.Link, error) {
			return netlink.LinkByName(name)
		},
	}
}

// LPAlias creates a link provider that when called, will provide the
// pre-existing link with the given alias (in the namespace this is called in).
// If no matches are found, an error is returned.
func LPAlias(alias string) LinkProvider {
	return LinkProvider{
		name: "alias",
		f: func() (netlink.Link, error) {
			return netlink.LinkByAlias(alias)
		},
	}
}

// LPIndex creates a link provider that when called, will provide the
// pre-existing link with the given index (in the namespace this is called in).
// If no matches are found, an error is returned.
func LPIndex(index int) LinkProvider {
	return LinkProvider{
		name: "index",
		f: func() (netlink.Link, error) {
			return netlink.LinkByIndex(index)
		},
	}
}
