#!/usr/bin/env sh

for d in $(go tool dist list | egrep 'android|linux')
do
    go env -w GOOS=$(echo $d | cut -d "/" -f 1) GOARCH=$(echo $d | cut -d "/" -f 2)
    go build -o "./build/$d/word" main.go
done;