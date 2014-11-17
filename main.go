package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/GeertJohan/ango/parser"
	goflags "github.com/jessevdk/go-flags"
)

var flagsParser *goflags.Parser
var flags struct {
	Verbose        bool   `long:"verbose" short:"v" description:"Enable verbose logging"`
	ForceOverwrite bool   `long:"force-overwrite" description:"Force overwrite (don't ask user)"`
	InputFile      string `long:"input" short:"i" description:"Input file" required:"true"`
	GoDir          string `long:"go-path" description:"Go output directory"`
	JsDir          string `long:"js-path" description:"Javascript output directory"`
	SkipJs         bool   `long:"skip-js" description:"Skip generation of Javascript code"`
	SkipGo         bool   `long:"skip-go" description:"Skip generation of Go code"`
}

var (
	// ErrNotImplementedYet is returned when something is not yet implemented
	ErrNotImplementedYet = errors.New("not implemented yet")
)

func main() {
	verbosef("ango version %s\n", versionFull())

	var err error
	flagsParser := goflags.NewParser(&flags, goflags.Default)

	args, err := flagsParser.Parse()
	if err != nil {
		_, ok := err.(*goflags.Error)
		if !ok {
			fmt.Printf("Error parsing flags: %s\n", err)
			os.Exit(1)
		}
		os.Exit(1)
	}
	if len(args) > 0 {
		fmt.Printf("Unexpected argument(s): '%s'\n", strings.Join(args, " "))
		os.Exit(1)
	}

	setupTemplates()

	inputFile, err := os.Open(flags.InputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %s\n", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	verbosef("Parsing %s.\n", flags.InputFile)
	angoParser := parser.NewParser(&parser.Config{
		PrintParseErrors: false,
	})
	service, err := angoParser.Parse(inputFile)
	if err != nil {
		fmt.Printf("Error parsing ango definitions: %s\n", err)
		os.Exit(1)
	}
	if service.Name != strings.TrimSuffix(filepath.Base(flags.InputFile), ".ango") {
		fmt.Println("Warning: .ango filename doesn't match name clause in file.")
	}
	verbosef("File %s parsed.", flags.InputFile)

	protocolVersion := calculateVersion(service)
	verbosef("Calculated protocol version is: %s\n", protocolVersion)

	if flags.SkipGo && flags.SkipJs {
		fmt.Printf("Parsed input file successfully.\nSkipping both Go and Javscript generation.\nGenerated version string was: %s\n", protocolVersion)
	}

	if !flags.SkipJs {
		err = generateJs(service)
		if err != nil {
			fmt.Printf("Error generating Javascript: %s\n", err)
			os.Exit(1)
		}
	}

	if !flags.SkipGo {
		err = generateGo(service)
		if err != nil {
			fmt.Printf("Error generating Go: %s\n", err)
			os.Exit(1)
		}
	}

	verbosef("ango main() completed\n")
}

func verbosef(format string, data ...interface{}) {
	if flags.Verbose {
		fmt.Printf(format, data...)
	}
}
