### [commitlog](/)

# Install

## From goblin

Goblin will build the project for your specific system so it might take a bit of time based on the traffic

```sh
curl -sf https://goblin.reaper.im/github.com/barelyhuman/commitlog
```

## From Releases

To avoid the above work load you can download the binaries from the github releases pages
[Commitlog Releases &rarr;](https://github.com/barelyhuman/commitlog/releases)

## From Source

If you have golang installed on your system you can build the binary on your system

```sh
# clone the repository
git clone https://github.com/barelyhuman/commitlog

# cd into the folder
cd commitlog

# run go mod to download needed deps
go mod tidy

# run go build and wait for it to complete
go build

# if on unix, use the `install` command to install it onto a system used path
install goblin /usr/local/bin
# or if you need perms
sudo install goblin /usr/local/bin
```
