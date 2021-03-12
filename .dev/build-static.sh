#!/bin/sh

go build -tags timetzdata --ldflags "-X 'main.Version=$(git describe --tags)' -linkmode external -extldflags \"-static\" -s -w" -o prettycrontab .
