package neslink

// TODO: Document all link actions
// TODO: Link actions for adding/deleting routes

import (
	"errors"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

// LinkAction is a singular operation that can be performed on a generic netlink
// link. Actions have a name as to identify individual actions when passed as a
// set to a LinkDo call, providing more contextual errors. They also have a
// function that take a link as a parameter. When called, the function will
// perform the operation on the provided link, returning an error if any
// occurred. These do support being executed outside of LinkDo calls, but
// using LinkDo is still recommended.
type LinkAction struct {
	actionName string
	f          func() error
}

// ActionName returns the name associated with the given link action.
func (la LinkAction) ActionName() string {
	return la.actionName
}

// act will perform the link operation immediately.
func (la LinkAction) act() error {
	return la.f()
}

// LAGeneric allows for a custom LinkAction to be created and then used in a
// LinkDo call.
func LAGeneric(actionName string, provider LinkProvider, function func() error) LinkAction {
	if actionName == "" {
		actionName = "unnamed-link-action"
	}
	return LinkAction{
		actionName: actionName,
		f:          function,
	}
}

// LADelete will simply delete the link when the action is executed. For obvious
// reasons this should be at the end of any LinkDo call (since the link will be
// deleted, further actions will error).
func LADelete(provider LinkProvider) LinkAction {
	return LinkAction{
		actionName: "delete-link",
		f: func() error {
			if l, err := provider.Provide(); err != nil {
				return errors.Join(errNoLink, err)
			} else {
				return netlink.LinkDel(l)
			}
		},
	}
}

func LASetName(provider LinkProvider, name string) LinkAction {
	return LinkAction{
		actionName: "set-name",
		f: func() error {
			if l, err := provider.Provide(); err != nil {
				return errors.Join(errNoLink, err)
			} else {
				return netlink.LinkSetName(l, name)
			}
		},
	}
}

func LASetAlias(provider LinkProvider, alias string) LinkAction {
	return LinkAction{
		actionName: "set-alias",
		f: func() error {
			if l, err := provider.Provide(); err != nil {
				return errors.Join(errNoLink, err)
			} else {
				return netlink.LinkSetAlias(l, alias)
			}
		},
	}
}

func LASetHw(provider LinkProvider, addr string) LinkAction {
	return LinkAction{
		actionName: "set-hw",
		f: func() error {
			if l, err := provider.Provide(); err != nil {
				return errors.Join(errNoLink, err)
			} else {
				hwAddr, err := net.ParseMAC(addr)
				if err != nil {
					return err
				}
				return netlink.LinkSetHardwareAddr(l, hwAddr)
			}
		},
	}
}

func LASetUp(provider LinkProvider) LinkAction {
	return LinkAction{
		actionName: "set-state-up",
		f: func() error {
			if l, err := provider.Provide(); err != nil {
				return errors.Join(errNoLink, err)
			} else {
				return netlink.LinkSetUp(l)
			}
		},
	}
}

func LASetDown(provider LinkProvider) LinkAction {
	return LinkAction{
		actionName: "set-state-down",
		f: func() error {
			if l, err := provider.Provide(); err != nil {
				return errors.Join(errNoLink, err)
			} else {
				return netlink.LinkSetDown(l)
			}
		},
	}
}

func LASetPromiscOn(provider LinkProvider) LinkAction {
	return LinkAction{
		actionName: "set-promisc-on",
		f: func() error {
			if l, err := provider.Provide(); err != nil {
				return errors.Join(errNoLink, err)
			} else {
				return netlink.SetPromiscOn(l)
			}
		},
	}
}

func LASetPromiscOff(provider LinkProvider) LinkAction {
	return LinkAction{
		actionName: "set-promisc-off",
		f: func() error {
			if l, err := provider.Provide(); err != nil {
				return errors.Join(errNoLink, err)
			} else {
				return netlink.SetPromiscOff(l)
			}
		},
	}
}

func LAAddAddr(provider LinkProvider, cidr string) LinkAction {
	return LinkAction{
		actionName: "add-address",
		f: func() error {
			if l, err := provider.Provide(); err != nil {
				return errors.Join(errNoLink, err)
			} else {
				addr, err := netlink.ParseAddr(cidr)
				if err != nil {
					return fmt.Errorf("failed to parse cidr to network address: %w", err)
				}
				return netlink.AddrAdd(l, addr)
			}
		},
	}
}

func LADelAddr(provider LinkProvider, cidr string) LinkAction {
	return LinkAction{
		actionName: "del-address",
		f: func() error {
			if l, err := provider.Provide(); err != nil {
				return errors.Join(errNoLink, err)
			} else {
				addr, err := netlink.ParseAddr(cidr)
				if err != nil {
					return fmt.Errorf("failed to parse cidr to network address: %w", err)
				}
				return netlink.AddrDel(l, addr)
			}
		},
	}
}
