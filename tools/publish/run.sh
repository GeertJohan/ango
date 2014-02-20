#!/bin/sh

# get publish-tool dependencies
go get github.com/GeertJohan/ango/tools/publish

# run publish-tool
go run tools/publish/*.go