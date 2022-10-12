package main

import (
	"fmt"

	"github.com/vishvananda/netlink"
	"github.com/willfantom/neslink"
)

func main() {
	// 1. create a new named netns "example" and create a new bridge "test" and
	// set the bridge mac address. note that if the named netns already exists,
	// this will error.
	if err, ok := neslink.LinkDo(
		neslink.NPNewNamed("example"),
		neslink.LPNewBridge("test"),
		neslink.LASetHw("34:34:34:34:34:34"),
	); err != nil || !ok {
		panic(err)
	}

	// 2. using an nsdo, fetch the links in the new "example" netns like in the
	// linkdo example.
	fmt.Println("New NS Links:")
	links := make([]netlink.Link, 0)
	err, ok := neslink.NsDo(neslink.NPName("example"), neslink.NALinks(&links))
	if err != nil {
		panic(err)
	} else if !ok {
		panic(fmt.Errorf("failed to return to origin namespace"))
	}

	// 3. delete the "example" netns
	if err := neslink.DeleteNamedNs("example"); err != nil {
		panic(err)
	}

	// 4. dump the links that were found in the "example" netns
	fmt.Printf("links found: %d\n", len(links))
	for _, l := range links {
		fmt.Printf("\t%d: %s: %s\n", l.Attrs().Index, l.Attrs().Name, l.Attrs().HardwareAddr.String())
	}
}
