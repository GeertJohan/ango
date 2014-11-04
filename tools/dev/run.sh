#!/bin/sh

# install rerun
go get -u -a github.com/skelterjohn/rerun

# get dev-tool dependencies
go get -u -a github.com/GeertJohan/ango/tools/dev

# run dev-tool
go run tools/dev/*.go