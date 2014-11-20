package main

import (
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/GeertJohan/ango/definitions"
	// "os"
	"sort"
)

var (
	versionNumber  string
	versionHash    string
	instanceUnique = randomString(10)
)

func versionFull() string {
	// is this a 'release' build?
	if len(versionNumber) > 0 {
		return versionNumber + "-" + versionHash
	}

	// is this a 'latest' build?
	if len(versionHash) > 0 {
		return "latest-" + versionHash
	}

	// asume local (development?) build
	return "local-" + instanceUnique
}

func calculateVersion(service *definitions.Service) string {
	hasher := sha256.New()

	// wr := io.MultiWriter(hasher, os.Stdout) // debug writer
	wr := io.Writer(hasher)

	// write ango version
	fmt.Fprintln(wr, versionFull())

	// write service name
	fmt.Fprintln(wr, service.Name)

	// write server procedures
	fmt.Fprint(wr, "server:\n")
	calculateVersionProcedures(wr, service.ServerProcedures)

	// write client procedures
	fmt.Fprint(wr, "client:\n")
	calculateVersionProcedures(wr, service.ClientProcedures)

	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func calculateVersionProcedures(hasher io.Writer, procs map[string]*definitions.Procedure) {

	// sort keys for server procedures
	keys := make([]string, 0, len(procs))
	for key := range procs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// create signature
	for _, key := range keys {
		proc := procs[key]
		if proc.Oneway {
			fmt.Fprint(hasher, "oneway ")
		}
		fmt.Fprintf(hasher, "%s(", proc.Name)
		calculateVersionParams(hasher, proc.Args)
		fmt.Fprint(hasher, ")")
		if len(proc.Rets) > 0 {
			fmt.Fprint(hasher, "(")
			calculateVersionParams(hasher, proc.Rets)
			fmt.Fprint(hasher, ")")
		}
		fmt.Fprint(hasher, "\n")
	}
}

func calculateVersionParams(hasher io.Writer, params []*definitions.Param) {
	for _, param := range params {
		fmt.Fprintf(hasher, "%s %s,", param.Name, param.Type)
	}
}
