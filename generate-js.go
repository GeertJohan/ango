package main

import (
	"fmt"
	"github.com/GeertJohan/ango/parser"
	"github.com/GeertJohan/go.ask"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type dataJs struct {
	ProtocolVersion string
	Service         *parser.Service
}

func generateJs(service *parser.Service) error {
	var err error
	var outputDir string
	if filepath.IsAbs(flags.JsDir) {
		outputDir = flags.JsDir
	} else {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting PWD: %s\n", err)
			os.Exit(1)
		}
		outputDir = filepath.Join(wd, flags.JsDir)
	}

	// create outputFile
	outputFileName := fmt.Sprintf("ango-%s.gen.js", service.Name)
	outputFileAbs := filepath.Join(outputDir, outputFileName)
	var outputFile *os.File
	outputFile, err = os.OpenFile(outputFileAbs, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		if os.IsExist(err) {
			// output file exists, ask user if we should overwrite.
			if flags.ForceOverwrite || ask.MustAskf("File '%s' exists, overwrite?", filepath.Join(flags.JsDir, outputFileName)) {
				outputFile, err = os.OpenFile(outputFileAbs, os.O_TRUNC|os.O_WRONLY, 0666)
				if err != nil {
					return err
				}
			} else {
				fmt.Println("Won't continue.")
				os.Exit(1)
			}
		} else {
			return err
		}
	}
	defer outputFile.Close()

	//prepare data
	data := &dataJs{
		ProtocolVersion: calculateVersion(service),
		Service:         service,
	}

	// intermediate io.WriteClose to abstract js-beautify or os.File away from template
	var outputWriteCloser io.WriteCloser

	// raise cmdJsb in scope to be able to wait for it to be done (when non-nil)
	var cmdJsb *exec.Cmd

	// check if js-beautify is installed
	if _, err = exec.LookPath("js-beautify"); err == nil {
		// js-breautify seems to be installed. Let's use it.
		cmdJsb = exec.Command("js-beautify", "--stdin", "--indent-with-tabs")
		outputWriteCloser, err = cmdJsb.StdinPipe()
		if err != nil {
			fmt.Printf("Error opening StdinPipe on js-beautify: %s\n", err)
			os.Exit(1)
		}
		cmdJsb.Stdout = outputFile
		err = cmdJsb.Start()
		if err != nil {
			fmt.Printf("Error starting js-beautify: %s\n", err)
			os.Exit(1)
		}
	} else {
		// write directly to file
		outputWriteCloser = outputFile
	}

	// execute template
	err = tmplJs.Execute(outputWriteCloser, data)
	if err != nil {
		fmt.Printf("Error executing javascript template: %s\n", err)
		os.Exit(1)
	}

	// close outputWriter (either closing stdin on cmdJsb or )
	err = outputWriteCloser.Close()
	if err != nil {
		fmt.Printf("Error closing outputWriteCloser: %s\n", err)
		os.Exit(1)
	}

	// wait for js-beautify to be done (if it looks like it's running)
	if cmdJsb != nil {
		err = cmdJsb.Wait()
		if err != nil {
			fmt.Printf("Error waiting for js-beautify to be done: %s\n", err)
			os.Exit(1)
		}
	}

	// all done
	return nil
}
