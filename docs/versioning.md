### [commitlog](/)

## Version Management

There's not much to version management as it's a very simple handler of semver increments.

### Initialise Versioning

To work with `commitlog release` you'll need a `.commitlog.release` file in your repo and that can be added using the following command. This will set the initial version to be `v0.0.0` when starting and you can sync it with you current package version manually by changing the file.
(We wish to add automatic syncing with major language version management systems but I'll need help due to low bandwidth)

```sh
$ commitlog release init
```

### Increments

You can then run the `commitlog release` command to initialise the TUI flow for selecting the next version or you can do the same using the flags if you are in a **non-interactive** environment.

**Note**: Increments are of the standard semver spec with the minor change of pre-releases being called **beta** releases
