# Commitlog

Changelog generator using Commit History

![Binary Builds](https://github.com/barelyhuman/commitlog/workflows/Binary%20Builds/badge.svg)

# UNDER DEV

it is usable by manualling installing the go package for now and isn't available as a binary right now, still under heavy development

To see an example of this in action, you can check the actions file for this repo. Yes it uses itself to generate the release logs

```sh
commitlog path/to/repository
```

### Example Output (from this exact repository)

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

## Contribution

You are free to raise PR's for any kind of improvements and or feature additions, they'll be added to the repository accordingly, you can also fork up your own version
