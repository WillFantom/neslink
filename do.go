package neslink

import (
	"fmt"
	"runtime"
)

// NsDo executes a given function in a specified network namespace. It does so
// in a separate OS thread in order to allow the rest of the program to continue
// on the current network namespace. Returned is the error (if any occurred) and
// a bool, that if false, shows that the function did not successfully return
// the netns back to the origin, a potentially critical error.
func NsDo(nsP NsProvider, actions ...NsAction) (error, bool) {
	// 1. get origin network namespace fd to revert back to
	originNs, err := NPNow()()
	if err != nil {
		return fmt.Errorf("failed to reference the calling network namespace: %w", err), true
	}
	defer originNs.Close()
	// 2. get new network namespace fd to switch to
	targetNs, err := nsP()
	if err != nil {
		return fmt.Errorf("failed to reference the target network namespace: %w", err), true
	}
	defer targetNs.Close()
	// 3. create error channel for new routine
	errChan := make(chan error, 1)
	defer close(errChan)
	// 4. create new go routine
	go func(oNs, tNs NsFd, nsActions ...NsAction) {
		// 1. lock os thread for goroutine
		runtime.LockOSThread()
		// 2. create defer to unlock os thread
		defer runtime.UnlockOSThread()
		// 3. switch to new netns
		if err := tNs.Set(); err != nil {
			errChan <- fmt.Errorf("failed to set netns to the target: %w", err)
			return
		}
		defer oNs.Set()
		// 4. exec ns actions
		for idx, action := range nsActions {
			if err := action.f(); err != nil {
				errChan <- fmt.Errorf("failed to perform ns action %d (%s): %w", idx+1, action.actionName, err)
				return
			}
		}
		errChan <- nil
	}(originNs, targetNs, actions...)
	// 5. get error from goroutine
	err = <-errChan
	// 6. check to make sure the final ns of the function matches the origin
	finalNs, nsErr := NPNow()()
	if nsErr != nil {
		return fmt.Errorf("failed to reference the final network namespace for post-run checks: %w", nsErr), false
	}
	defer finalNs.Close()
	// 7. return
	return err, (finalNs.ID() == originNs.ID())
	// TODO: check potential performance issue do to defer stack? maybe switch to anonymous function?
}

// LinkDo runs a set of link actions on a link that is obtained from the given
// LinkProvider. The link provider is called and the actions are performed in
// the namespace obtained by the NsProvider, thus this can be used to manage
// links in any namespace. The actions are perfromed in the given namespace via
// NsDo, so the returned outputs are much the same as with NsDo.
func LinkDo(nsP NsProvider, lP LinkProvider, actions ...LinkAction) (error, bool) {
	function := func() error {
		link, err := lP()
		if err != nil {
			return fmt.Errorf("failed to get link from provider: %w", err)
		}
		for idx, action := range actions {
			if err := action.f(link); err != nil {
				return fmt.Errorf("failed to perform action %d (%s): %w", idx+1, action.actionName, err)
			}
		}
		return nil
	}
	return NsDo(nsP, NAGeneric("link-action-set", function))
}

// MustNsDo performs a normal NsDo, however if it fails to revert back to the
// origin netns, then this will trigger a panic. Any other errors are returned
// as normal.
func MustNsDo(nsP NsProvider, actions ...NsAction) error {
	err, ok := NsDo(nsP, actions...)
	if !ok {
		panic(fmt.Errorf("nsdo failed to revert to origin ns: %w", err))
	}
	return err
}

// MustLinkDo performs a normal LinkDo, however if it fails to revert back to
// the origin netns, then this will trigger a panic. Any other errors are
// returned as normal.
func MustLinkDo(nsP NsProvider, lP LinkProvider, actions ...LinkAction) error {
	err, ok := LinkDo(nsP, lP, actions...)
	if !ok {
		panic(fmt.Errorf("linkdo failed to revert to origin ns: %w", err))
	}
	return err
}
