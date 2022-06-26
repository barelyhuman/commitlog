<!-- meta -->
<title>
    commitlog
</title>
<meta name="description" content="commits to changelog generator">
<!-- meta end -->

### [commitlog](/)

[Manual](/manual)\
[Download &darr;](/download)

#### Index

- [Source](https://github.com/barelyhuman/commitlog)
- [Quick Start](#quick-start)
- [About](#about)
  - [Philosophy](#philosophy)
  - [Installation](#installation)
- [Manual](/manual)
- [License](#license)

<h4 id="quick-start">Quick Start</h4>

For the one's who already know the tool and wanna get started quickly.

You can get the CLI in the following ways

1. You can get the binaries from the [download section](/download)

2. Using `go get`

```sh
go get -u github.com/barelyhuman/commitlog
```

3. Using goblin

```sh
curl -sf https://goblin.barelyhuman.xyz/github.com/barelyhuman/commitlog | sh
```

Once installed you can just run `commitlog generate` or the shorter version
`commitlog g` to get all the changes between the recent tags.

<h4 id="about">About</h4>

<h5 id="philosophy">Philosophy</h5>

**_commitlog_** simply exists because I wasn't able to find a decent enough tool
to do this without being bound to a specific programming language and it's
toolset. The closes one to this is the `git`'s own `shortlog` and `log` commands

- **Language Agnostic** - No one needs the entire programming toolset or the
  language to be setup on the system to be able to just generate logs

- **Decently Sized** - Not everyone has a powerful system and CI systems are
  still very limited in terms of resources since they give you shared spaces and
  automating generating this is a common practice so it doesn't make sense for
  it to install a runtime worth 100MB to run it.

  The binary itself is around 10-13MB which is fairly big for a tool like this
  but this is due to go's nature of embedding the platform runtime with the
  binary and this might not be ideal for some and so there's another similar
  smaller scoped version [nimclog](https://github.com/barelyhuman/nimclog) which
  is in KB's

- **Flexible** - There are no set standards for how you write commit messages.
  Commitlint is a good standard and pretty widely used and also what commitlog
  v1 used as the base for categorizing.

  V2 however, doesn't have a set commit category standard, it accepts regex's
  that you will categorize your commits for you.

  This is more for experienced developers who have a their own pattern of
  writing commits and would like a tool that'd work with their pattern instead
  of creating friction

- **Extensible** - The entire thing is open source and MIT licenses,
  specifically for you to be able to fix and add things as needed. Each command
  has been kept away from the domain logic to be able to use `commitlog` even as
  a library and build your own tool if you wish to, though if there's something
  you wish the CLI supported, you are free to open issues and I really like
  people helping with the development since that helps keeping the tool alive
  longer.

<h5 id="installation">Installation</h5>

The installation is pretty straightforwad. You [download](/download) the binary
from the download section of this website or use something like
[goblin](https://goblin.barelyhuman.xyz) to build the binary for your specific
system (rarely needed) since the releases actually accommodate the most used
operating systems and architectures already.

**Linux/Mac (Unix Systems)**

Once downloaded, you can use the `install` command on *nix systems to link the
binary to a location that's already in your PATH variable.

eg:

```bash
# install commitlog from current directory to the /usr/local/bin directory
install $(pwd)/commitlog /usr/local/bin
```

This should give you the ability to use the `commitlog` command anywhere in your
system

**Windows**

Similar to the linux setup, you can download the `.exe` file and add it to a
location that your environment path already has. An easier thing to do is to
store the entire thing in your secondary partition and add that location to your
PATH. You can use this resource from Java docs to help you with modifying the
PATH variables

- [Modifying PATH in Windows](https://www.java.com/en/download/help/path.html)

<h4 id="manual">Manual</h4>

[Read Manual &rarr;](/manual)

<h4 id="license">License</h4>

commitlog is MIT Licensed and you can read the entire license in the [source code](https://github.com/barelyhuman/commitlog)
