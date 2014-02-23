package main

import (
	"fmt"
	"github.com/GeertJohan/ango/parser"
	"github.com/GeertJohan/go.ask"
	"os"
	"path/filepath"
)

type dataGo struct {
	PackageName     string
	ProtocolVersion string
	Service         *parser.Service
}

func generateGo(service *parser.Service) error {
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
		outputDir = filepath.Join(wd, flags.GoDir)
	}

	// create outputFile
	outputFileName := fmt.Sprintf("ango-%s.gen.go", service.Name)
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
	data := &dataGo{
		PackageName:     flags.GoPackage,
		ProtocolVersion: calculateVersion(service),
		Service:         service,
	}

	// execute template
	err = tmplGo.Execute(outputFile, data)
	if err != nil {
		fmt.Printf("Error executing go template: %s\n", err)
		os.Exit(1)
	}

	// all done
	return nil
}
