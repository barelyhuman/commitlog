<!-- meta -->
<title>
    commitlog | manual
</title>
<meta name="description" content="commits to changelog generator">
<!-- meta end -->

### [commitlog](/)

[Home](/) [Manual](/manual) [Download](/download) [API](/api)

# Manual

- [Basics](#basics)
  - [generate](#generate)
  - [release](#release)
- [Advanced](#advanced)
  - [Selective Log Ranges](#selective-log-ranges)
  - [From a different repo](#from-a-different-repo)
  - [Categorize](#categorize)
- [CLI](#cli)
- [API](/api)

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

The default behaviour of **_commitlog_** is to print to the stdout of your
terminal which might not be what you wish to do.

So, if you wish to ever write to a file , you can use the `--out FILE` flag to
write the output to a file instead.

```sh
$ commitlog g -o ./CHANGELOG.md
```

> **Note**: commitlog will overwrite the file with the new output, make sure you
> create a backup of the file when working with this.

> **Note**: In coming versions an append mode will be added to the CLI to avoid
> having to deal with this manually

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

These are not really hard to guess or achieve but are in the advanced section as
this isn't something a beginner user might need to mess with. If you have other
use cases that you feel like should be added here, raise a issue.

<h4 id="selective-log-ranges">Selective Log Ranges</h4>
There's enough cases where the entire history between the tags makes no sense. You might just want a specific set of commits and this is very easy to do with commitlog

```sh
$ commitlog g --start HEAD --end HEAD~3 
# would return the last 3 commits
```

And as you might have understood, commitlog does work with git's revision system
so you can reuse git's own revision pointers as needed. Though it's mostly
easier to just use the HASH since, finding the hashes in `git log` is easier
than finding the number you have to type after the `HEAD~`

<h4 id="from-a-different-repo">From a different repo</h4>
A lot more often you might need to generate logs for a sub-repo when working with a huge multirepo folder which a lot of sub-modules and switching worktree for the changelog might not be something you wish to do.

In cases like these the `-p` or `--path` flag can be used with the `generate`
flag

```sh
$ commitlog g -p ./sub/module -s HEAD~3 
# would return the commits from HEAD~3 to the next tag or end of all commits.
```

<h4 id="categorize">Categorize</h4>
You wouldn't be using commitlog otherwise so this obviously has to exist.

As for why is this turned off by default? I just don't use commit standards and
mostly just need the commits between tags and very rarely need the categories to
help me but doesn't change the fact that someone else might need it.

So, how does commitlog categorize? using regex's , or more like comma seperated
regex's.

The implementation approach might bring some issue but I'm going to wait for
them instead of over fixing something that might never really happen.

Anyway, a simple example for a non-scoped commitlint styled categorization would
look something like this .

```sh
$ commitlog g --categories="feat:,fix:"
```

**Output:**

```md
## feat:
110c9b260e6de1e2be9af28b3592c373ab6fc501 feat: handling for pre increment

9080580575d3579b3d11dd114830eb128f0c8130 feat: releaser base maj,min,patch setup

b513e53cbbc30e5b667fbba373b9589911a47ac0 feat: generator subcommand implementation

## fix:
5bafe9d105298001707c7c816e60d3ef0257a816 fix: tag creation
```

This could be made a little more neater using a better regex like so

```sh
$ commitlog g --categories="^feat,^fix"
```

**Output**:

```md
## ^feat
e0ce9c693c554158b7b4841b799524164d9c3e83 feat(releaser): add init and pre-tag handlers

110c9b260e6de1e2be9af28b3592c373ab6fc501 feat: handling for pre increment

9080580575d3579b3d11dd114830eb128f0c8130 feat: releaser base maj,min,patch setup

b513e53cbbc30e5b667fbba373b9589911a47ac0 feat: generator subcommand implementation

## ^fix
5bafe9d105298001707c7c816e60d3ef0257a816 fix: tag creation

540bb3d7419aab314fdb999dd57eeac3c53f5fad fix doc generation

87d40ddbd463ad71b3af8b7edb8ac8176ecbf2f5 fix tests for new param
```

As visible, you find more commits in the 2nd one since the regex matches more
cases and gives you more information. You are free to modify these to make them
work with other patterns that you might follow.

Also the titles are going to be the regex's you pass. A future version will add
a way for you to add label's to the regex though that's a little harder to
achieve without making the flags very verbose so I've avoided it for now.

> **Note**: This is the initial implementation of this and might change in a
> future major version since this flag feels a little more frictional than the
> other approaches I have in mind. Feel free to recommend other ways to do this.

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
