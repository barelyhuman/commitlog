<!-- meta -->
<title>
    commitlog | manual
</title>
<meta name="description" content="commits to changelog generator">
<!-- meta end -->

### [commitlog](/)

[Manual](/manual)\
[Download &darr;](/download)

# Manual

- [Basics](#basics)
  - [generate](#generate)
  - [release](#release)
- [Advanced](#advanced)
- [CLI](#cli)

<h3 id="basics">Basics</h3>

**_`commitlog`_** come with 2 basic functionalities

- Generating commits as changelogs in markdown
- Managing your version and tags

Both the above are available as sub commands and we'll go through each one by
one.

Just running `commitlog` no longer works as it did in v1, you'll need to specify
a sub command to get it working.

This has been done to keep it simpler to keep arbitrary behaviour to a minimum

<h4 id="generate">generate</h3>

The first one is `generate` and in most cases you'll just be using this to
handle generation of commits as a list of markdown notes

Here's an example of what it'll look like to run the `generate` sub-command

```sh
$ commitlog generate
# or 
$ commitlog g
```

**Output**:

```md
## All Changes
87ce3fc76b252908cdd81f9537713f8d67813652 chore: v2.0.0-beta.0

09fc85c0d1d52d770bbbf39ac45907103a53d354 refactor: swap go-git for push with exec

5bafe9d105298001707c7c816e60d3ef0257a816 fix: tag creation

909fd032facdeb4fb7633a4fad82dced28e8de83 chore: swap debugging logger with stderr printer
<!-- .... more commits till the end of tag -->
```

The behaviour of the `generate` command is like so

- go through all commits from the current `HEAD` (commit with the `HEAD`
  reference) till the last commit or the most recent tagged commit.
- the other case is where if the latest commit is also a tagged commit, it'll
  still go till the most recent tagged commit.

  eg: v0.0.2 -> HEAD, v0.0.1 -> SOMEHASH, would give you the commits between
  both the tags.

This is the default behavior of commitlog and can be changed by passing various
flags which you'll read about soon.

<h4 id="release">release</h3>

The `release` sub-command is a complimentary sub-command added to commitlog to
help you with managing release versions. I understand that most languages have
their own package version management and you might want to work with them but
that's also why the `release` sub-command isn't something you need, it's there
for cases where you might wanna use it.

My primary usecase is when working with repositories that work with multiple
languages and I'd prefer a unified version.

The `release` command doesn't modify your source code files but instead handles
the version with it's own file named `.commitlog.release`, which just holds the
**version** number right now and while this does have the X.X.X format (semantic
versioning format) the increments or how you handle them doesn't have to be.

The command is flexible enough for you to combine and create your own way of
versioning.

##### Initialize

```sh
$ commitlog release --init
# or 
$ commitlog r --init
```

Once initialized, you'll see a `.commitlog.release` file in your repo with the
version `v0.0.0` to start with.

##### Increments

Post that you can use the remaining flags to help you incrementing the version
as needed.

```sh
$ commitlog r --major # change to v1.0.0
$ commitlog r --minor # change to v0.1.0
$ commitlog r --patch # change to v0.0.1
$ commitlog r --pre # change to v0.0.0-beta.0
$ commitlog r --pre --pre-tag="dev" # change to v0.0.0-dev.0

# or go bonkers 
$ commitlog r --patch --pre --pre-tag="dev" # change to v0.0.1-dev.0
$ commitlog r --major --patch --pre --pre-tag="dev" # change to v1.0.1-dev.0
```

#### Git actions

Most times you'd also want a tool like this to create commits and tags for you
which is fair, so `release` does come with a few git actions.

- `--commit` (will create a commit and tag it for you)
- `--push` (will push the commits and tags for you)

The whole command would look a little something like this

```sh
$ commitlog r --patch --commit --push
```

<h3 id="advanced">Advanced</h3>

**TBD**

<h3 id="cli">CLI</h3>

This is the entire CLI reference, which you can also access by using the `-h`
flag on the downloaded binary.

**`commitlog -h`**

```
NAME:
   commitlog - commits to changelogs

USAGE:
   commitlog [global options] command [command options] [arguments...]

COMMANDS:
   generate, g  commits to changelogs
   release, r   manage .commitlog.release version
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

---

**`commitlog g -h`**

```
NAME:
   commitlog generate - commits to changelogs

USAGE:
   commitlog generate [command options] [arguments...]

OPTIONS:
   --categories value       categories to use, includes all commits by default. any text you add here will be used to create categories out of the commits
   --end END, -e END        END reference for the commit to stop including commits at.This is exclusive of the given commit reference
   --out FILE, -o FILE      path to the output FILE
   --path PATH, -p PATH     root with the '.git' folder PATH (default: ".")
   --promo                  add promo text to the end of output (default: false)
   --start START, -s START  START reference for the commit to include commits from,This is inclusive of the given commit reference
   --stdio                  print to the stdout (default: true)
```

---

**`commitlog r -h`**

```
NAME:
   commitlog release - manage .commitlog.release version

USAGE:
   commitlog release [command options] [arguments...]

OPTIONS:
   --commit              if true will create a commit and tag, of the changed version (default: false)
   --init                initialise commitlog release (default: false)
   --major               create a major version (default: false)
   --minor               create a minor version (default: false)
   --patch               create a patch version (default: false)
   --path PATH, -p PATH  PATH to a folder where .commitlog.release exists or is to be created.(note: do not use `--commit` or `--push` if the file isn't in the root) (default: ".")
   --pre                 create a pre-release version. will default to patch increment unlessspecified and not already a pre-release (default: false)
   --pre-tag value       create a pre-release version (default: "beta")
   --push                if true will create push the created release commit + tag on origin (default: false)
```
