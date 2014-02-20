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
	"runtime"
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

const exampleAngoFile = "example/chatService.ango"

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
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error
	fmt.Println("Starting ango dev tool.")

	wd, err = os.Getwd()
	if err != nil {
		fmt.Printf("Error getting wd: %s\n", err)
		stop(1)
		select {}
	}

	pkg, err = build.ImportDir(wd, 0)
	if err != nil {
		fmt.Printf("Error loading package: %s\n", err)
		stop(1)
		select {}
	}

	if pkg.Name != "main" || !pkg.IsCommand() || filepath.Base(wd) != "ango" {
		fmt.Println("Is tool executed from the right directory? (github.com/GeertJohan/ango or a fork)?")
		fmt.Printf("Current package (%s) is invalid.\n", pkg.Name)
		stop(1)
		select {}
	}

	go rerunExample()

	go rerunAngo()

	go watchExampleAngo()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Kill, os.Interrupt)
	sig := <-sigChan
	signal.Stop(sigChan)

	fmt.Printf("Received %s, closing...\n", sig)
	if sig == os.Interrupt {
		stop(0)
	} else {
		stop(1)
	}
	select {}
}

func stop(exitCode int) {
	go func() {
		// synced stop
		close(stopCh)
		stopWg.Wait()

		// os.Exit(..)
		os.Exit(exitCode)
	}()
}

func rerunExample() {
	stopWg.Add(1)
	defer stopWg.Done()

	cmdRerun := exec.Command("rerun", filepath.Join(pkg.ImportPath, "example"))
	cmdRerun.Stdin = os.Stdin
	cmdRerun.Stdout = sgr.NewColorWriter(os.Stdout, sgr.FgYellow, false)
	cmdRerun.Stderr = sgr.NewColorWriter(os.Stderr, sgr.FgYellow, false)
	err := cmdRerun.Start()
	if err != nil {
		fmt.Printf("Error running rerun example: %s\n", err)
		stop(1)
		return
	}
	<-stopCh
	if cmdRerun.Process != nil {
		cmdRerun.Process.Signal(os.Interrupt)
	}
	err = cmdRerun.Wait()
	if err != nil && err.Error() != "exit status 2" {
		fmt.Printf("Error stopping rerun example: %s\n", err)
		return
	}
}

func rerunAngo() {
	stopWg.Add(1)
	defer stopWg.Done()

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
		fmt.Printf("Error running rerun ango build: %s\n", err)
		stop(1)
		return
	}
	<-stopCh
	if cmdRerun.Process != nil {
		cmdRerun.Process.Signal(os.Interrupt)
	}
	err = cmdRerun.Wait()
	if err != nil && err.Error() != "exit status 2" {
		fmt.Printf("Error stopping rerun ango build: %s\n", err)
		return
	}
}

func watchExampleAngo() {
	stopWg.Add(1)
	defer stopWg.Done()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Error starting watcher: %s\n", err)
		stop(1)
		return
	}
	defer watcher.Close()
	err = watcher.WatchFlags(exampleAngoFile, fsnotify.FSN_MODIFY)
	if err != nil {
		fmt.Printf("Error starting watch on example ango file: %s\n", err)
		stop(1)
		return
	}
	for {
		select {
		case <-stopCh:
			return
		case <-watcher.Event:
			angoExample()
		case err := <-watcher.Error:
			if err != nil {
				fmt.Printf("Error watching example ango file: %s\n", err)
				stop(1)
				return
			}
		}
	}
}

func angoExample() {
	stopWg.Add(1)
	defer stopWg.Done()
	fmt.Println("Running ango tool for example/chatService.ango")
	cmdAngoExample := exec.Command(filepath.Join(wd, "ango"), "--verbose", "-i", exampleAngoFile, "--js", "example/http-files", "--force-overwrite") // "--go", "example",
	cmdAngoExample.Stdin = os.Stdin
	cmdAngoExample.Stdout = sgr.NewColorWriter(os.Stdout, sgr.FgBlue, false)
	cmdAngoExample.Stderr = sgr.NewColorWriter(os.Stderr, sgr.FgBlue, false)
	err := cmdAngoExample.Run()
	if err != nil {
		fmt.Printf("Error running ango tool: %s\n", err)
	}
}
