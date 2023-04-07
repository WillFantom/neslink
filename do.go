package neslink

import (
	"errors"
	"fmt"
	"runtime"
)

var (
	// errDirtyThread is returned when some action that moves a thread over to a
	// netns fails to return the thread back to the netns of the caller. In the
	// scenario where this happens, the os thread can be considered dirty and
	// should not be reused. This error may also be wrapped into others, so
	// errors.Is should be used to check for its presence.
	errDirtyThread error = errors.New("system thread failed to move to expected final network namespace")
)

// NsDo executes a given function in a specified network namespace. It does so
// in a separate OS thread in order to allow the rest of the program to continue
// on the current network namespace. An error is returned if any netns move
// fails or any provided action fails. Do note that if the spawned system thread
// fails to be reverted to the network namespace of the caller, the thread is
// considered dirty and is never unlocked (thus can not be reused).
func NsDo(nsP NsProvider, actions ...NsAction) error {
	// 1. get origin network namespace fd to revert back to
	originNsFd, err := NPNow().Provide().open()
	if err != nil {
		return fmt.Errorf("failed to open the origin netns file descriptor: %w", err)
	}
	defer originNsFd.close()

	// 2. get new network namespace fd to switch to
	targetNsFd, err := nsP.Provide().open()
	if err != nil {
		return fmt.Errorf("failed to open the target netns file descriptor: %w", err)
	}
	defer targetNsFd.close()

	// 3. create error channel for new routine
	errChan := make(chan error, 1)
	defer close(errChan)

	// 4. create new go routine
	go func(oNs, tNs NsFd, nsActions ...NsAction) {

		// 1. lock os thread for goroutine
		runtime.LockOSThread()

		// 2. switch to new netns
		if err := tNs.set(); err != nil {
			errChan <- fmt.Errorf("failed to set netns to the target: %w", err)
			return
		}

		// -?- thread now dirty - perpare for cleanup
		errSet := errors.Join(nil)

		// 3. exec ns actions
		for idx, action := range nsActions {
			if err := action.act(); err != nil {
				errSet = errors.Join(errSet, fmt.Errorf("failed to perform ns action %d (%s)", idx+1, action.actionName), err)
				break
			}
		}

		// 4. switch to origin netns
		if err := oNs.set(); err != nil {
			errSet = errors.Join(errSet, fmt.Errorf("failed to switch to origin ns"), err, errDirtyThread)
		}

		// 5. if thread is dirty, dont unlock thread and sleep routine forever
		if !errors.Is(errSet, errDirtyThread) {
			runtime.UnlockOSThread()
			errChan <- errSet
		} else {
			errChan <- errSet
			dirtyThreadSleeper := make(chan struct{})
			<-dirtyThreadSleeper
		}

	}(originNsFd, targetNsFd, actions...)

	// 5. get error from goroutine and return
	return <-errChan

}

// LinkDo runs a set of link actions on a link that is obtained from the given
// LinkProvider. The link provider is called and the actions are performed in
// the namespace obtained by the NsProvider, thus this can be used to manage
// links in any namespace. The actions are perfromed in the given namespace via
// NsDo, so the returned outputs are much the same as with NsDo.
func LinkDo(nsP NsProvider, lP LinkProvider, actions ...LinkAction) error {
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

func init() {
	runtime.LockOSThread()
}
