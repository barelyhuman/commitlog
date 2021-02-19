# Commitlog

Changelog generator using Commit History

![Binary Builds](https://github.com/barelyhuman/commitlog/workflows/Binary%20Builds/badge.svg) 
![Binary Builds Pre-Releases](https://github.com/barelyhuman/commitlog/workflows/Binary%20Builds%20Pre-Releases/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/barelyhuman/commitlog)](https://goreportcard.com/report/github.com/barelyhuman/commitlog)

To see an example of this in action, you can check the actions file for this repo. Yes it uses itself to generate the release logs

## Binaries 
Binaries are available on the Github [releases](https://github.com/barelyhuman/commitlog/releases)

## Build 
This step is for people who want to use the latest version from the repositories which hasn't been added to releases as a binary yet and for people viewing this on sourcehut , as the binaries aren't uploaded to sourcehut.

- Make sure you have go installed on your system , minimum version `1.15`
- You can either clone the whole repo from sourcehut / github or downlad a tar.gz from sourcehut / github 

### With Clone 

```sh 
git clone https://git.sr.ht/~reaper/commitlog
# or 
git clone https://github.com/barelyhuman/commitlog

cd commitlog 
go build 
```

### With Tarballs 
```sh
tar -xvzf commitlog-<hash>.tar.gz
cd commitlog 
go build 
```

```sh
# to install it to go's bin folder and use the commitlog command during dev or as a perm install
go install 
```

## Usage 

The usage is pretty simple, this cli tool assumes that you use [commitlint standards](https://github.com/conventional-changelog/commitlint#what-is-commitlint) while writing your commits, if not it's okay everything will be classified under `Other Changes` instead of being grouped according to type of commit. 

**Simple Overview**  
`ci: <message>` - for ci/cd changes   
`feat|feature: <message>` - for feature changes  
`docs: <message>` - for documents or comment updations in code  
`refactor: <message>` - performance / code clean up changes or total BLOC changes  
`fix: <message>` - for fixes (self-explanatory)  


### Supported Flags (as of v0.0.4-dev.3)

The below mentioned are the flags supported by the current branch and older tags might not support the flags 
or certain inputs in flags, use the `-h` flag to see what's supported on the version you are using.

```sh
Usage of commitlog:
  -e string
        commit hash string / revision (ex. HEAD, HEAD^, HEAD~2) 
         to stop collecting commit message at
  -i string
        commit types to be includes (default "ci,refactor,docs,fix,feat,test,chore,other")
  -p string
        path to the repository, points to the current working directory by default (default ".")
  -s string
        commit hash string / revision (ex. HEAD, HEAD^, HEAD~2) 
         to start collecting commit messages from
  -skip
        if true will skip trying to classify and just give a list of changes
```

 
### Example Output (from this exact repository)

```sh
> commitlog 
```

```markdown
# Changelog

## Fixes

97c582b3eb5a6796ef9c250d9653ad90dce63cbe - fix: example fix

## Other Changes

da6d837eb3134f836bfbe401de7882f2e0818ba8 - Create LICENSE
b0f1b1d2bc4265cb72b70b3ae5b60f8e65f47b12 - initial commit
```

## Current Limitations

- No Tests added so is probably unstable right now ( If you'd like to help writing tests, feel free to raise a PR)
- Doesn't work inside a nested folder in the repository, needs to be on the root of the repository to work. (you can use the `-p` flag to set a path to generate logs from)


## Contribution

[Contribution Guidelines](CONTRIBUTING.md)
