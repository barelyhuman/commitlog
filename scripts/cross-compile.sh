#!/usr/bin/env bash

set -euxo pipefail

rm -rf ./bin

build_commands=('
    apk add make curl git \
    ; GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/linux-arm64/alvu \
    ; GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/linux-amd64/alvu \
    ; GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o bin/linux-arm/alvu \
    ; GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o bin/windows-386/alvu \
    ; GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/windows-amd64/alvu \
    ; GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/darwin-amd64/alvu \
    ; GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/darwin-arm64/alvu 
')

# run a docker container with osxcross and cross compile everything
docker run -it --rm -v $(pwd):/usr/local/src -w /usr/local/src \
	golang:alpine3.16 \
	sh -c "$build_commands"

# create archives
cd bin
for dir in $(ls -d *);
do
    cp ../README.md $dir
    cp ../LICENSE $dir
    mkdir -p $dir/docs
    cp -r ../docs/* $dir/docs
    
    # remove the download document and styles
    rm -rf $dir/docs/download.md
    rm -rf $dir/docs/styles.css

    tar cfzv "$dir".tgz $dir
    rm -rf $dir
done
cd ..