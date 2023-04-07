package main

import (
	"fmt"

	"github.com/willfantom/nescript"
	"github.com/willfantom/neslink"
)

func main() {
	// 1. define the actual script to run, apply mutations, then "compile"
	script := nescript.NewScript("sleep 3 && {{.Command}} address list").WithField("Command", "ip").MustCompile()

	// 2. create an empty process to store the newly created process in
	process := *new(nescript.Process)

	// 3. execute the script via the appropriate ns action in an nsdo call
	err := neslink.Do(neslink.NPNow(), neslink.NANewNs("example"), neslink.NAExecNescript(script, nil, &process), neslink.NADeleteNamed("example"))
	if err != nil {
		panic(err)
	}

	// 4. wait for the result from the process
	result, err := process.Result()
	if err != nil {
		panic(err)
	}

	// 5. dump the result
	fmt.Printf("Command Output:\n%s", result.StdOut)
}
