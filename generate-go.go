package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"

	"github.com/GeertJohan/ango/definitions"
	"github.com/GeertJohan/go.ask"
)

type dataGo struct {
	PackageName     string
	ProtocolVersion string
	Service         *definitions.Service
}

func generateGo(service *definitions.Service) error {
	var err error
	var outputDir string
	if filepath.IsAbs(flags.GoDir) {
		outputDir = flags.GoDir
	} else {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting PWD: %s\n", err)
			os.Exit(1)
		}
		if flags.GoDir == "" {
			outputDir = filepath.Join(wd, filepath.Dir(flags.InputFile), service.Name)
		} else {
			outputDir = filepath.Join(wd, flags.GoDir)
		}
	}

	//prepare data
	data := &dataGo{
		PackageName:     service.Name,
		ProtocolVersion: calculateVersion(service),
		Service:         service,
	}

	// execute template into buffer
	generatedSourceBuffer := &bytes.Buffer{}
	err = tmplGo.Execute(generatedSourceBuffer, data)
	if err != nil {
		fmt.Printf("Error executing go template: %s\n", err)
		os.Exit(1)
	}

	// format generated source
	formattedSource, err := format.Source(generatedSourceBuffer.Bytes())
	if err != nil {
		fmt.Printf("error formatting source: %v\n", err)
		formattedSource = generatedSourceBuffer.Bytes()
	}

	// create outputFile
	outputFileName := "server.gen.go"
	outputFileAbs := filepath.Join(outputDir, outputFileName)
	var outputFile *os.File
	outputFile, err = os.OpenFile(outputFileAbs, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		if os.IsExist(err) {
			// output file exists, ask user if we should overwrite.
			if flags.ForceOverwrite || ask.MustAskf("File '%s' exists, overwrite?", outputFileAbs) {
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

	_, err = outputFile.Write(formattedSource)
	if err != nil {
		fmt.Printf("error writing formatted source to file: %v\n", err)
	}

	// all done
	return nil
}
