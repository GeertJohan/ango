#!/bin/sh

# install rerun
go get github.com/GeertJohan/rerun

# get dev-tool dependencies
go get github.com/GeertJohan/ango/tools/dev

# run dev-tool
go run tools/dev/*.go