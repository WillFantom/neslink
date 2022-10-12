package neslink

import (
	"fmt"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"github.com/willfantom/nescript"
	"github.com/willfantom/nescript/local"
)

// NsAction represents an action that should be executed in a namespace via
// NsDo. The action should have a relevant name as to give context to errors (as
// multiple actions are executed in a single NsDo call). Also the action itself
// should be a function that takes no parameters and returns an error (or nil in
// the event of success).
type NsAction struct {
	actionName string
	f          func() error
}

// NAGeneric allows for a custom action (function) to be performed in a given
// network namespace. A name should be given to describe the custom function in
// a couple of words to give context to NsDo errors.
func NAGeneric(name string, function func() error) NsAction {
	if name == "" {
		name = "unnamed-action"
	}
	return NsAction{
		actionName: name,
		f:          function,
	}
}

// NASetLinkNs moves a link provided by the given link provider to the namespace
// provided by the ns provider. The link itself should br present in the
// namespace in which the wrapping NsDo is set to execute in.
func NASetLinkNs(lP LinkProvider, nsP NsProvider) NsAction {
	return NsAction{
		actionName: "set-link-ns",
		f: func() error {
			link, err := lP()
			if err != nil {
				return fmt.Errorf("failed to obtain link from provider: %w", err)
			}
			ns, err := nsP()
			if err != nil {
				return fmt.Errorf("failed to obtain target netns for link from provider: %w", err)
			}
			defer ns.Close()
			return netlink.LinkSetNsFd(link, int(ns))
		},
	}
}

// NAGetNsFd provides an open file descriptor for the network namespace it is
// called in. This fd is separate from that of the one in the enclosing NsDo, so
// it is up to the user to close the fd when it is no longer needed.
func NAGetNsFd(nsfd *NsFd) NsAction {
	return NsAction{
		actionName: "get-ns-fd",
		f: func() error {
			fd, err := netns.Get()
			if err != nil {
				return err
			}
			*nsfd = NsFd(fd)
			return nil
		},
	}
}

// NAExecNescript will execute a NEScript in the netns it is called in, most
// likely the netns of the wrapping NsDo. This opens up extensive custom
// options. Provided should be the already compiled NEScript, a subcommand to
// use for the script such as ["sh" "-c"] (or nil to use the nescript package's
// deafult), and a nescript.Process for the resulting process to be stored in.
func NAExecNescript(script nescript.Script, subcommand []string, process *nescript.Process) NsAction {
	return NsAction{
		actionName: "exec-nescript",
		f: func() error {
			p, err := script.Exec(local.Executor("", subcommand))
			if err != nil {
				return err
			}
			*process = p
			return nil
		},
	}
}

// ListLinks returns a list of all the links in the namespace obtained via the
// given provider. Any errors are returned and a boolean to express if the the
// network namespace has returned back to the origin successfully.
func NALinks(links *[]netlink.Link) NsAction {
	return NsAction{
		actionName: "get-ns-links",
		f: func() error {
			l, err := netlink.LinkList()
			if err != nil {
				return err
			}
			*links = l
			return nil
		},
	}
}

// NAGetLink gets a specific link from the given link provider when the action
// is called. The result is stored within the given link parameter. An error is
// returned if any occurred.
func NAGetLink(provider LinkProvider, link *netlink.Link) NsAction {
	return NsAction{
		actionName: "get-ns-link",
		f: func() error {
			l, err := provider()
			if err != nil {
				return err
			}
			*link = l
			return nil
		},
	}
}

// func NADumpFilepath() NsAction {
// 	return NsAction{
// 		actionName: "dump-file-path",
// 		f: func() error {
// 			nsfd, err := netns.Get()
// 			if err != nil {
// 				return err
// 			}
// 			name, err := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", int(nsfd)))
// 			if err != nil {
// 				return err
// 			}
// 			fmt.Println(name)
// 			return nil
// 		},
// 	}
// }
