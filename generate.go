package main

import (
	"github.com/GeertJohan/ango/parser"
	"github.com/davecgh/go-spew/spew"
)

func generate(service *parser.Service) {
	spew.Dump(service)

	verbosef("TODO: generate files\n")
}
