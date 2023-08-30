<!-- meta -->
<title>
    commitlog | manual
</title>
<meta name="description" content="commits to changelog generator">
<!-- meta end -->

### [commitlog](/)

[Home](/) [Manual](/manual) [Download](/download) [API](/api)

# API

### General Guide

commitlog also comes as a pkg that you could reuse to modify the behaviour of
the commands and this is very limited at the moment since I'm still working on
the best way to get plugins to work with the original CLI instead of having to
write your own version of commitlog.

The pkg contains the 2 base command's creators and behaviour modifiers, or more
commonly known as the golang options pattern.

Briefly put, You have one function that takes in unlimited amount of functions
as parameter with each of these parameter functions being able to modify the
behaviour of the returned instance.

Easy way to explain this is with an example of the `releaser` API

```go
package main

import (
	"log"

	"github.com/barelyhuman/commitlog/v2/pkg"
)

func main() {
	versionString := "v0.0.1"
	releaser, _ := pkg.CreateNewReleaser(
		versionString,
		pkg.WithMajorIncrement(),
	)

	log.Println(releaser.String())

}
```

here the `pkg.CreateNewReleaser` takes in one mandatory value which is the
`versionString` and the 2nd parameter is optional here.

Though, since we wish for the releaser to have a custom behaviour everytime the
flags change, instead of writing entire functionalities inside various releaser
functions, which would look like so

```go
releaser.IncrementMajor()
releaser.IncrementMinor()
```

I'd be unable to expose the builders / option functions out to the public for
them to write custom behaviours that work directly with the `struct` being used
by commitlog itself and instead you'd be writing wrappers around existing
functions. Thus, adding another layer of abstraction which isn't needed for
something that wants to be extended.

This approach gives me the ability to expose a selected few properties for you
to modify while writing your own builder/option function.

The only pointer functions that releaser has is the one's that'll help with
printing or identifying the final version's state.

Since, you now know how the API is written, the go doc for this module should be
able to help you with the remaining.

[godoc&nearr;](https://pkg.go.dev/github.com/barelyhuman/commitlog)

> **Note**: if the go doc still hasn't been generated for v2.0.0, please go
> through the source code to help you with the implementation details
