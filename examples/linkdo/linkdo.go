package main

import (
	"fmt"

	"github.com/vishvananda/netlink"
	"github.com/willfantom/neslink"
)

func main() {
	// 1. create a new named netns "example" and create a new bridge "test" and
	// set the bridge mac address. note that if the named netns already exists,
	// this will error.{
	if err := neslink.Do(neslink.NPNow(),
		neslink.NANewNs("example"),
		neslink.LANewBridge("testbr"),
		neslink.LASetHw(neslink.LPName("testbr"), "34:34:34:34:34:34"),
	); err != nil {
		panic(err)
	}

	// 2. using a do, fetch the links in the new "example" netns like in the
	// do example.
	fmt.Println("New NS Links:")
	links := make([]netlink.Link, 0)
	err := neslink.Do(neslink.NPName("example"), neslink.NALinks(&links), neslink.NADeleteNamed("example"))
	if err != nil {
		panic(err)
	}

	// 4. dump the links that were found in the "example" netns
	fmt.Printf("links found: %d\n", len(links))
	for _, l := range links {
		fmt.Printf("\t%d: %s: %s\n", l.Attrs().Index, l.Attrs().Name, l.Attrs().HardwareAddr.String())
	}
}
