package neslink

import (
	"errors"
	"fmt"
	"net"

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

// LPNewBridge creates a link provider that when called, will create a new
// bridge with the given name and then provides the newly created bridge.
func LPNewBridge(name string) LinkProvider {
	return LinkProvider{
		name: "new-bridge",
		f: func() (netlink.Link, error) {
			bridge := netlink.Bridge{
				LinkAttrs: netlink.NewLinkAttrs(),
			}
			bridge.LinkAttrs.Name = name
			if err := netlink.LinkAdd(&bridge); err != nil {
				return nil, err
			}
			return netlink.LinkByName(name)
		},
	}
}

// LPNewVeth creates a link provider that when called, will create a new veth
// pair and return the link. The names for both the new interfaces (main link
// and peer) should be provided. Only the main link is actually then provided by
// the LinkProvider. To get the peer, LAName with the peer's name should
// suffice. Any errors in creating/finding the veth pair are also returned.
func LPNewVeth(name, peerName string) LinkProvider {
	return LinkProvider{
		name: "new-veth",
		f: func() (netlink.Link, error) {
			veth := netlink.Veth{
				LinkAttrs: netlink.NewLinkAttrs(),
				PeerName:  peerName,
			}
			veth.LinkAttrs.Name = name
			if err := netlink.LinkAdd(&veth); err != nil {
				return nil, err
			}
			return netlink.LinkByName(name)
		},
	}
}

// LPNewDummy creates a link provider that when called, will create a new dummy
// link with the given name and returns it, provided no errors occur.
func LPNewDummy(name string) LinkProvider {
	return LinkProvider{
		name: "new-dummy",
		f: func() (netlink.Link, error) {
			dummy := netlink.Dummy{
				LinkAttrs: netlink.NewLinkAttrs(),
			}
			dummy.LinkAttrs.Name = name
			if err := netlink.LinkAdd(&dummy); err != nil {
				return nil, err
			}
			return netlink.LinkByName(name)
		},
	}
}

// LPNewGRETap creates a link provider that when called, creates a new gretap
// device with the given name, local IP, and remoteIP.
func LPNewGRETap(name, localIP, remoteIP string) LinkProvider {
	return LinkProvider{
		name: "new-gretap",
		f: func() (netlink.Link, error) {
			local := net.ParseIP(localIP)
			if local == nil {
				return nil, fmt.Errorf("failed to parse the local ip address of the gretap")
			}
			remote := net.ParseIP(remoteIP)
			if remote == nil {
				return nil, fmt.Errorf("failed to parse the remote ip address of the gretap")
			}
			gre := netlink.Gretap{
				LinkAttrs: netlink.NewLinkAttrs(),
				Local:     local,
				Remote:    remote,
			}
			gre.LinkAttrs.Name = name
			if err := netlink.LinkAdd(&gre); err != nil {
				return nil, err
			}
			return netlink.LinkByName(name)
		},
	}
}

// LPNewWireguard creates a link provider that when called, will create a new
// wireguard link with the given name and returns it, provided no errors occur.
// Further setup of this link should be done in custom LinkDos withwireguard
// specifc code.
func LPNewWireguard(name string) LinkProvider {
	return LinkProvider{
		name: "new-wireguard",
		f: func() (netlink.Link, error) {
			wg := netlink.Wireguard{
				LinkAttrs: netlink.NewLinkAttrs(),
			}
			wg.LinkAttrs.Name = name
			if err := netlink.LinkAdd(&wg); err != nil {
				return nil, err
			}
			return netlink.LinkByName(name)
		},
	}
}

// LPNewVxlan creates a link provider that when called, will create a new
// vxlan link with the given configuration and returns it.
func LPNewVxlan(name, localIP, groupIP string, id, port int) LinkProvider {
	return LinkProvider{
		name: "new-vxlan",
		f: func() (netlink.Link, error) {
			local := net.ParseIP(localIP)
			if local == nil {
				return nil, fmt.Errorf("failed to parse the local ip address of the vxlan")
			}
			group := net.ParseIP(groupIP)
			if group == nil {
				return nil, fmt.Errorf("failed to parse the group ip address of the vxlan")
			}
			vx := netlink.Vxlan{
				LinkAttrs: netlink.NewLinkAttrs(),
				VxlanId:   id,
				SrcAddr:   local,
				Group:     group,
				Port:      port,
			}
			vx.LinkAttrs.Name = name
			if err := netlink.LinkAdd(&vx); err != nil {
				return nil, err
			}
			return netlink.LinkByName(name)
		},
	}
}
