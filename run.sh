#!/bin/bash
docker build -t lxm-ency .
docker rm --force lxm-ency
docker run --name lxm-ency -e ENV="DEBUG" --restart=always -p8000:8000 lxm-ency