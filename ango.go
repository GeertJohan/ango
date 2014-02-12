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
	Verbose   bool   `long:"verbose" short:"v" description:"Enable verbose logging"`
	OutputDir string `long:"output" short:"o" description:"Output directory" required:"true"`
	InputFile string `long:"input" short:"i" description:"Input file" required:"true"`
}

func main() {
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

	inputFile, err := os.Open(flags.InputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %s\n", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	parseTree, err := parser.Parse(inputFile)
	if err != nil {
		fmt.Printf("Error parsing ango definitions: %s\n", err)
		os.Exit(1)
	}

	generate(parseTree)

	verbosef("ango main() completed\n")
}

func verbosef(format string, data ...interface{}) {
	if flags.Verbose {
		fmt.Printf(format, data...)
	}
}
