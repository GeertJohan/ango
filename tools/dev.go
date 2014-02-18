package main

import (
	"fmt"
	"github.com/foize/go.sgr"
	"github.com/howeyc/fsnotify"
	"go/build"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
)

var (
	wd string // working directory

	pkg *build.Package

	monAngo *exec.Cmd // command builds ango on file change (uses gomon)

	// watcher on the ango binary
	watcher *fsnotify.Watcher

	// closed on stop
	stopCh = make(chan bool)
	stopWg sync.WaitGroup
)

type CheckWriter struct {
	wr     io.Writer
	filter string
	action func()
}

func (c *CheckWriter) Write(b []byte) (n int, err error) {
	n, err = c.wr.Write(b)
	if strings.Contains(string(b), c.filter) {
		c.action()
	}
	return n, err
}
func main() {
	var err error
	fmt.Println("Starting ango dev tool.")

	wd, err = os.Getwd()
	if err != nil {
		fmt.Printf("Error getting wd: %s\n", err)
		stop(nil)
	}

	pkg, err = build.ImportDir(wd, 0)
	if err != nil {
		fmt.Printf("Error loading package: %s\n", err)
		stop(nil)
	}

	if pkg.Name != "main" || !pkg.IsCommand() || filepath.Base(wd) != "ango" {
		fmt.Println("Is tool executed from the right directory? (github.com/GeertJohan/ango or a fork)?")
		fmt.Printf("Current package (%s) is invalid.\n", pkg.Name)
		stop(nil)
	}

	go rerunExample()

	go rerunAngo()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Kill, os.Interrupt)
	sig := <-sigChan
	signal.Stop(sigChan)

	fmt.Printf("Received %s, closing...\n", sig)
	stop(sig)
}

func stop(sig os.Signal) {
	var exitCode int
	if sig != os.Interrupt {
		exitCode = 1
	}
	if sig == nil {
		sig = os.Interrupt
	}

	// synced stop
	close(stopCh)
	stopWg.Wait()

	// os.Exit(..)
	os.Exit(exitCode)
}

func rerunExample() {
	stopWg.Add(1)

	cmdRerun := exec.Command("rerun", filepath.Join(pkg.ImportPath, "example"))
	cmdRerun.Stdin = os.Stdin
	cmdRerun.Stdout = sgr.NewColorWriter(os.Stdout, sgr.FgYellow, false)
	cmdRerun.Stderr = sgr.NewColorWriter(os.Stderr, sgr.FgYellow, false)
	err := cmdRerun.Start()
	if err != nil {
		fmt.Printf("Error running rerun example: %s\n", err)
	}
	<-stopCh
	if cmdRerun.Process != nil {
		cmdRerun.Process.Signal(os.Interrupt)
	}
	stopWg.Done()
}

func rerunAngo() {
	stopWg.Add(1)

	cmdRerun := exec.Command("rerun", "-build-only", pkg.ImportPath)
	cmdRerun.Stdin = os.Stdin
	cw := &CheckWriter{
		wr:     os.Stderr,
		filter: "build passed",
		action: func() {
			angoExample()
		},
	}
	cmdRerun.Stdout = sgr.NewColorWriter(os.Stdout, sgr.FgCyan, false)
	cmdRerun.Stderr = sgr.NewColorWriter(cw, sgr.FgCyan, false)
	err := cmdRerun.Start()
	if err != nil {
		fmt.Printf("Error running rerun example: %s\n", err)
	}
	<-stopCh
	if cmdRerun.Process != nil {
		cmdRerun.Process.Signal(os.Interrupt)
	}
	stopWg.Done()
}

func angoExample() {
	stopWg.Add(1)
	fmt.Println("Running ango tool for example/example.ango")
	cmdAngoExample := exec.Command(filepath.Join(wd, "ango"), "--verbose", "-i", "example/example.ango", "--js", "example/http-files", "--force-overwrite") // "--go", "example",
	cmdAngoExample.Stdin = os.Stdin
	cmdAngoExample.Stdout = sgr.NewColorWriter(os.Stdout, sgr.FgBlue, false)
	cmdAngoExample.Stderr = sgr.NewColorWriter(os.Stderr, sgr.FgBlue, false)
	err := cmdAngoExample.Run()
	if err != nil {
		fmt.Printf("Error running ango tool: %s\n", err)
	}
	stopWg.Done()
}
