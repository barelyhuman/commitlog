# Commitlog

Changelog generator using Commit History

![Binary Builds](https://github.com/barelyhuman/commitlog/workflows/Binary%20Builds/badge.svg)

# UNDER DEV

You can use the exisiting binaries from the releases page or install using `go get` though it is still under heavy development and will not be getting any test coverage until I'm done with most of the mvp functionality.

To see an example of this in action, you can check the actions file for this repo. Yes it uses itself to generate the release logs

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

- Cannot Specify a commit/tag to consider as the source of log
- No Tests added so is probably unstable right now

## Note 
 The current code base is a prototypal solution and while usable will be heavily refactored to follow best practices, if you'd like to take help of a better codebase if you are creating your own fork then you can take a refactored codebase from [https://github.com/percybolmer/commitlog](https://github.com/percybolmer/commitlog)

## Contribution

[Contribution Guidelines](CONTRIBUTING.md)
