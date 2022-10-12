package neslink

// TODO: Document all link actions
// TODO: Link actions for adding/deleting routes

import (
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
	f          func(netlink.Link) error
}

// ActionName returns the name associated with the given link action.
func (la LinkAction) ActionName() string {
	return la.actionName
}

// Do will perform the link operation immediately, supporting use outide of any
// LinkDo calls.
func (la LinkAction) Do(link netlink.Link) error {
	return la.f(link)
}

// LAGeneric allows for a custom LinkAction to be created and then used in a
// LinkDo call.
func LAGeneric(name string, function func(netlink.Link) error) LinkAction {
	if name == "" {
		name = "unnamed-link-action"
	}
	return LinkAction{
		actionName: name,
		f:          function,
	}
}

// LADelete will simply delete the link when the action is executed. For obvious
// reasons this should be at the end of any LinkDo call (since the link will be
// deleted, further actions will error).
func LADelete() LinkAction {
	return LinkAction{
		actionName: "delete-link",
		f: func(l netlink.Link) error {
			return netlink.LinkDel(l)
		},
	}
}

func LASetName(name string) LinkAction {
	return LinkAction{
		actionName: "set-name",
		f: func(l netlink.Link) error {
			return netlink.LinkSetName(l, name)
		},
	}
}

func LASetAlias(alias string) LinkAction {
	return LinkAction{
		actionName: "set-alias",
		f: func(l netlink.Link) error {
			return netlink.LinkSetAlias(l, alias)
		},
	}
}

func LASetHw(addr string) LinkAction {
	return LinkAction{
		actionName: "set-hw",
		f: func(l netlink.Link) error {
			hwAddr, err := net.ParseMAC(addr)
			if err != nil {
				return err
			}
			return netlink.LinkSetHardwareAddr(l, hwAddr)
		},
	}
}

func LASetUp() LinkAction {
	return LinkAction{
		actionName: "set-state-up",
		f: func(l netlink.Link) error {
			return netlink.LinkSetUp(l)
		},
	}
}

func LASetDown() LinkAction {
	return LinkAction{
		actionName: "set-state-down",
		f: func(l netlink.Link) error {
			return netlink.LinkSetDown(l)
		},
	}
}

func LASetPromiscOn() LinkAction {
	return LinkAction{
		actionName: "set-promisc-on",
		f: func(l netlink.Link) error {
			return netlink.SetPromiscOn(l)
		},
	}
}

func LASetPromiscOff() LinkAction {
	return LinkAction{
		actionName: "set-promisc-off",
		f: func(l netlink.Link) error {
			return netlink.SetPromiscOff(l)
		},
	}
}

func LAAddAddr(cidr string) LinkAction {
	return LinkAction{
		actionName: "add-address",
		f: func(l netlink.Link) error {
			addr, err := netlink.ParseAddr(cidr)
			if err != nil {
				return fmt.Errorf("failed to parse cidr to network address: %w", err)
			}
			return netlink.AddrAdd(l, addr)
		},
	}
}

func LADelAddr(cidr string) LinkAction {
	return LinkAction{
		actionName: "del-address",
		f: func(l netlink.Link) error {
			addr, err := netlink.ParseAddr(cidr)
			if err != nil {
				return fmt.Errorf("failed to parse cidr to network address: %w", err)
			}
			return netlink.AddrDel(l, addr)
		},
	}
}
