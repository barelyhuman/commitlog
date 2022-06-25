#!/usr/bin/env bash

mkdir -p docs 
cp ./README.md ./docs/index.md
curl -o mudkip.tgz -L https://github.com/barelyhuman/mudkip/releases/latest/download/linux-amd64.tgz
tar -xvzf mudkip.tgz
install linux-amd64/mudkip /usr/local/bin

mudkip --stylesheet="./docs/styles.css" --baseurl='/commitlog/'


