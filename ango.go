package main

import (
	"fmt"
	"github.com/GeertJohan/ango/parser"
	goflags "github.com/jessevdk/go-flags"
	"os"
	"strings"
)

var flagsParser *goflags.Parser
var flags struct {
	Verbose        bool   `long:"verbose" short:"v" description:"Enable verbose logging"`
	ForceOverwrite bool   `long:"force-overwrite" description:"Force overwrite (don't ask user)"`
	InputFile      string `long:"input" short:"i" description:"Input file" required:"true"`
	GoDir          string `long:"go" description:"Go output directory"`
	JsDir          string `long:"js" description:"Javascript output directory"`
}

// Version constant is used in calculating the protocol version
const Version = `0.1`

func main() {
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

	parser.PrintParseErrors = false
	parseTree, err := parser.Parse(inputFile)
	if err != nil {
		fmt.Printf("Error parsing ango definitions: %s\n", err)
		os.Exit(1)
	}

	protocolVersion := calculateVersion(parseTree)
	verbosef("Calculated protocol version is: %s\n", protocolVersion)

	// warn user about abscent --go and --js
	// give version string
	if len(flags.GoDir) == 0 && len(flags.JsDir) == 0 {
		fmt.Printf("Parsed input file. There were no errors.\nUse options `--go <outputDir>` and `--js <outputDir>` to generate code.\nVersion string is: %s\n", protocolVersion)
	}

	if len(flags.JsDir) > 0 {
		generateJs(parseTree)
	}

	if len(flags.GoDir) > 0 {
		generateGo(parseTree)
	}

	verbosef("ango main() completed\n")
}

func verbosef(format string, data ...interface{}) {
	if flags.Verbose {
		fmt.Printf(format, data...)
	}
}
