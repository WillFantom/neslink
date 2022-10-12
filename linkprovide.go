package neslink

// TODO: Support more new link types (Dummy, GRETap, Wireguard, VxLan)

import "github.com/vishvananda/netlink"

type LinkProvider func() (netlink.Link, error)

// LPName creates a link provider that when called, will provide the link with
// the given name (in the namespace this is called in). If no matches are
// found, an error is returned.
func LPName(name string) LinkProvider {
	return func() (netlink.Link, error) {
		return netlink.LinkByName(name)
	}
}

// LPAlias creates a link provider that when called, will provide the link with
// the given alias (in the namespace this is called in). If no matches are
// found, an error is returned.
func LPAlias(alias string) LinkProvider {
	return func() (netlink.Link, error) {
		return netlink.LinkByAlias(alias)
	}
}

// LPIndex creates a link provider that when called, will provide the link with
// the given index (in the namespace this is called in). If no matches are
// found, an error is returned.
func LPIndex(index int) LinkProvider {
	return func() (netlink.Link, error) {
		return netlink.LinkByIndex(index)
	}
}

// LPNewBridge creates a link provider that when called, will create a new
// bridge with the given name and then provides the newly created bridge.
func LPNewBridge(name string) LinkProvider {
	return func() (netlink.Link, error) {
		bridge := netlink.Bridge{
			LinkAttrs: netlink.NewLinkAttrs(),
		}
		bridge.LinkAttrs.Name = name
		if err := netlink.LinkAdd(&bridge); err != nil {
			return nil, err
		}
		return netlink.LinkByName(name)
	}
}

// LPNewVeth creates a link provider that when called, will create a new veth
// pair and return the link. The names for both the new interfaces (main link
// and peer) should be provided. Only the main link is actually then provided by
// the LinkProvider. To get the peer, LAName with the peer's name should
// suffice. Any errors in creating/finding the veth pair are also returned.
func LPNewVeth(name, peerName string) LinkProvider {
	return func() (netlink.Link, error) {
		veth := netlink.Veth{
			LinkAttrs: netlink.NewLinkAttrs(),
			PeerName:  peerName,
		}
		veth.LinkAttrs.Name = name
		if err := netlink.LinkAdd(&veth); err != nil {
			return nil, err
		}
		return netlink.LinkByName(name)
	}
}
