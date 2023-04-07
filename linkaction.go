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

// name simply returns the name of the link action.
func (la LinkAction) name() string {
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

// LANewBridge creates a new bridge with the given name.
func LANewBridge(name string) LinkAction {
	return LinkAction{
		actionName: "new-bridge",
		f: func() error {
			bridge := netlink.Bridge{
				LinkAttrs: netlink.NewLinkAttrs(),
			}
			bridge.LinkAttrs.Name = name
			return netlink.LinkAdd(&bridge)
		},
	}
}

// LANewVeth will create a new veth pair. The names for both the new interfaces
// (main link and peer) should be provided.
func LANewVeth(name, peerName string) LinkAction {
	return LinkAction{
		actionName: "new-veth",
		f: func() error {
			veth := netlink.Veth{
				LinkAttrs: netlink.NewLinkAttrs(),
				PeerName:  peerName,
			}
			veth.LinkAttrs.Name = name
			return netlink.LinkAdd(&veth)
		},
	}
}

// LANewDummy creates a new dummy link with the given name.
func LANewDummy(name string) LinkAction {
	return LinkAction{
		actionName: "new-dummy",
		f: func() error {
			dummy := netlink.Dummy{
				LinkAttrs: netlink.NewLinkAttrs(),
			}
			dummy.LinkAttrs.Name = name
			return netlink.LinkAdd(&dummy)
		},
	}
}

// LANewGRETap creates a new gretap device with the given name, local IP, and
// remoteIP.
func LANewGRETap(name, localIP, remoteIP string) LinkAction {
	return LinkAction{
		actionName: "new-gretap",
		f: func() error {
			local := net.ParseIP(localIP)
			if local == nil {
				return fmt.Errorf("failed to parse the local ip address of the gretap")
			}
			remote := net.ParseIP(remoteIP)
			if remote == nil {
				return fmt.Errorf("failed to parse the remote ip address of the gretap")
			}
			gre := netlink.Gretap{
				LinkAttrs: netlink.NewLinkAttrs(),
				Local:     local,
				Remote:    remote,
			}
			gre.LinkAttrs.Name = name
			return netlink.LinkAdd(&gre)
		},
	}
}

// LANewWireguard creates a new wireguard link with the given name. Further
// setup of this link should be done in custom LinkActions with wireguard
// specifc code.
func LANewWireguard(name string) LinkAction {
	return LinkAction{
		actionName: "new-wireguard",
		f: func() error {
			wg := netlink.Wireguard{
				LinkAttrs: netlink.NewLinkAttrs(),
			}
			wg.LinkAttrs.Name = name
			return netlink.LinkAdd(&wg)
		},
	}
}

// LANewVxlan creates a new vxlan link with the given configuration.
func LANewVxlan(name, localIP, groupIP string, id, port int) LinkAction {
	return LinkAction{
		actionName: "new-vxlan",
		f: func() error {
			local := net.ParseIP(localIP)
			if local == nil {
				return fmt.Errorf("failed to parse the local ip address of the vxlan")
			}
			group := net.ParseIP(groupIP)
			if group == nil {
				return fmt.Errorf("failed to parse the group ip address of the vxlan")
			}
			vx := netlink.Vxlan{
				LinkAttrs: netlink.NewLinkAttrs(),
				VxlanId:   id,
				SrcAddr:   local,
				Group:     group,
				Port:      port,
			}
			vx.LinkAttrs.Name = name
			return netlink.LinkAdd(&vx)
		},
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
