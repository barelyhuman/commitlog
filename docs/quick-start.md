### [commitlog](/)

# Quick Start

## Changelog

1. Install using one of the methods from [Installation Instructions](/install)
2. Change terminal directory into the git repository you want the changelogs from.
3. run the base command

   ```sh
   # on mac
   $ commitlog | pbcopy

   # on linux with x11
   $ commitlog | xclip
   ```

## Version Management

1. Install using one of the methods from [Installation Instructions](/install)
2. Change terminal directory into the git repository you want the changelogs from.
3. Initialise version management

   ```sh
   $ commitlog release init
   ```

4. If working in an existing project, change the `.commitlog.release` file to match with the current version of the project.
5. Run the needed command to increase versino
   ```sh
   $ commitlog release
   ```
