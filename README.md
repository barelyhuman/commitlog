# Commitlog

Changelog generator using Commit History

![Binary Builds](https://github.com/barelyhuman/commitlog/workflows/Binary%20Builds/badge.svg)

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

 
### Example Output (from this exact repository)

```sh
> commitlog path/to/repository
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

- No Tests added so is probably unstable right now

## Note 
 The current code base is a prototypal solution and while usable will be heavily refactored to follow best practices, if you'd like to take help of a better codebase if you are creating your own fork then you can take a refactored codebase from [https://github.com/percybolmer/commitlog](https://github.com/percybolmer/commitlog)

## Contribution

[Contribution Guidelines](CONTRIBUTING.md)
