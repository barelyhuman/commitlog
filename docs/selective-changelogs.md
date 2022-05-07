### [commitlog](/)

# Selective Changelogs

## Classification

Commitlog was built on a minor subset of classification flags from the commitlint spec
which includes the following type of commits.

`ci|refactor|docs|fix|feat|test|chore|other`

These are all turned on by default on the cli and will use the commit message's subject to classify commits for you.

Eg:
The following git message would generate a commitlog like so

```
# git message
feat: added handling for feature flags
```

```markdown
<!-- COMMITLOG -->

# Features

53cb64436a3df5b22702cf6d012a4d09df93d930 - added handling for feature flags
```

Everything that cannot be classified in the above types will be put into `other`

### Selected types

If you only want to include feature and fix flags for changelogs you could do so in the following manner

```sh
$ commitlog -i "fix|feat"
# or
$ commitlog -i "fix,feat"
```

### Skip classification

If you wish to avoid classification all together and just want commits between a range you can do so as well (though you are better off without the tool for this usecase since git handles this natively)

```sh
$ commitlog --skip
```

## Commits in range

A very common use case is to get commits between a range if you aren't using tags in your current workflow.
This can be done by using **git revision handles** or the **commit hashes**.

Eg:
Let's say I want the last 3 commits

```sh
$ commitlog -s HEAD -e HEAD~4
```

<small>`HEAD~4` since the end commit is excluded from the range.</small>

Both `-s` and `-e` are individual flags and don't have to be used together.

So something like will also give you the last 4 commits.

```sh
$ commitlog -e HEAD~4
```

and the following would give you all commits till the next tagged instance starting form the given commit

```sh
$ commitlog -s HEAD~4
```
