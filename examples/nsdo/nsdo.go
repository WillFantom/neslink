package main

import (
	"fmt"

	"github.com/vishvananda/netlink"
	"github.com/willfantom/neslink"
)

func main() {
	// 1. create an empty links slice to store the result in
	links := make([]netlink.Link, 0)

	// 2. call do with a provider (in this case will be a new unnamed netns) and
	// the Links action
	err := neslink.Do(neslink.NPNow(), neslink.NANewNs("example"), neslink.NALinks(&links), neslink.NADeleteNamed("example"))
	if err != nil {
		panic(err)
	}

	// 3. dump the links that were found (since its a new netns, should just be
	// loopback)
	fmt.Printf("links found: %d\n", len(links))
	for _, l := range links {
		fmt.Printf("\t%d: %s\n", l.Attrs().Index, l.Attrs().Name)
	}
}
