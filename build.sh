#!/bin/bash

echo start building

env GOOS=linux GOARCH=amd64 go build -mod=vendor -o lxm-ency main.go
tar czvf lxm-ency.tar.gz lxm-ency data
scp lxm-ency.tar.gz root@39.96.21.121:/home/works/chenzheye
rm lxm-ency lxm-ency.tar.gz