# ~~net~~**nes**link
<!-- TODO: add go reference badge -->

NESlink is go package that allows for interaction with netlink. NESlink is simply a quality of life wrapper around these great [netlink](https://github.com/vishvananda/netlink) and [netns](https://github.com/vishvananda/netns) packages. Not all the functionality of these packages remains in NESlink, but the interactions that are included are better suited to the [NES]() platform implementation.

<!-- TODO: Add link to NES platform repository -->

The main objective of this package is to make link interaction easier when dealing with many links across many network namespaces. This includes safe parallel operations in multiple namespaces.

--- 

## Usage

At the core of this package are the `Do` functions, one for network namespaces, the other for links. These allow for low-code, go routine safe interaction with both links and namespaces (and links in namespaces).

### Namespace Interaction

Any interaction with a netns should be done via a call to `NsDo`. As an example, to list the links in a network namespace, you would simply need to provide `NsDo` with a `NsProvider` for the target namespace, and the `NsAction` for listing links. So if there was a network namespace called `example`, then the following snippet would perform the action safely:

```go
links := make([]netlink.Link, 0)
err, ok := neslink.NsDo(neslink.NPName("example"), neslink.NALinks(&links))
if err != nil || !ok {
  ...
```

> 💡 Any number of `NSActions` can be provided to a single `NsDo` call, and they will be executed in order.

Here `err` would contain any error that occurred either in switching namespaces or within the function. However, `ok` will only be true if the network namespace was successfully returned to the original netns prior to running the function.

A new netns can be created similarly by specifying with the `NPNew()` or `NPNewNamed(...)` provider. The actions set can simply be left blank if you wish to only create the netns and not interact with it.

Custom `NsActions` can be easily created too, see [this example](./examples/nsactions).

### Link Interaction

To manage links, any operation should be a `LinkAction` set in a call to `LinkDo`. Much like `NsDo`, `LinkDo` will execute a set of functions in a given netns, but applied to a specific link found via a `LinkProvider`. As an example, the below snippet will create a new bridge called _`br0`_ in a pre-existing named netns called _`example`_, then set its MAC address to _`12:23:34:45:56:67`_ and set its state to UP:

```go
if err, ok := neslink.LinkDo(
  neslink.NPName("example"),
  neslink.LPNewBridge("br0"),
  neslink.LASetHw("12:23:34:45:56:67"),
  neslink.LASetUp(),
  ); err != nil || !ok {
  ...
```

> 📝 Setting a link's netns is not a `LinkAction` but instead a `NsAction`, since after moving the link to another netns, the netns of the LinkDo goroutine should also be changed to the netns to complete any further actions on the link.

Via the `LinkProviders`, new links can be created, or already created links can be obtained via their name, index, or alias.

### NEScript Integration

Using this package, [NEScripts](https://github.com/willfantom/nescript) can be executed on any specific netns, making it easy to specify custom actions to execute via the `NsAction` system.

### Other Notes

 - Both `XXDo` functions have a `MustXXDo` variant that will still return errors, but if the function can not successfully revert to the original netns, will cause a panic.

--- 

## Motivation

Whilst the 2 packages referenced at the top of this doc for netlink and netns provide all this functionality and more, they are still somewhat low-level packages. This can result in programs that use them extensively needing a lot of wrapper code to make the provided functionality easier and safer to use. This package is that wrapper code.

--- 

## 🚧 WIP

 - [ ] Address inline TODOs
 - [ ] Add tests
 - [ ] Run tests via actions
