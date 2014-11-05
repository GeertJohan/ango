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
	"time"
)

var (
	wd string // working directory

	pkg *build.Package

	monAngo *exec.Cmd // command builds ango on file change (uses gomon)

	// last ango generation
	timeLastGenCodeOverwrite time.Time

	// genLock is write-locked when generated source was detected modified unexpectedly
	genLock sync.RWMutex

	// closed on stop
	stopCh = make(chan bool)
	stopWg sync.WaitGroup
)

const exampleAngoFile = "example/chatservice.ango"

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

	// watch generated source and lock on modification
	go watchGeneratedSource()

	// watch and rebuild/restart example cmd
	go rerunExample()

	// watch and rebuild ango tool
	go rebuildAngo()

	// watch ango tool and templates and re-generate on change
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

func watchGeneratedSource() {
	stopWg.Add(1)
	defer stopWg.Done()

	haveLine := make(chan bool)
	go func() {
		// clear stdin and read line when required
		for {
			fmt.Scanln()
			select {
			case haveLine <- true:
			default:
			}
		}
	}()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Error starting watcher: %s\n", err)
		stop(1)
		return
	}
	defer watcher.Close()
	err = watcher.WatchFlags(`example/chatservice/server.gen.go`, fsnotify.FSN_MODIFY)
	if err != nil {
		fmt.Printf("Error starting watch on go generated file: %s\n", err)
		stop(1)
		return
	}
	err = watcher.WatchFlags(`example/chatservice.gen.js`, fsnotify.FSN_MODIFY)
	if err != nil {
		fmt.Printf("Error starting watch on js generated file: %s\n", err)
		stop(1)
		return
	}
	for {
		select {
		case <-stopCh:
			return
		case <-watcher.Event:
			// short timeout consuming another event becuase sublime sometimes saves (modifies) the file twice
			select {
			case <-time.After(100 * time.Millisecond):
			case <-watcher.Event:
			}
			if time.Now().Sub(timeLastGenCodeOverwrite) > time.Duration(3*time.Second) {
				genLock.Lock()
				sgr.Println("[fg-red]Detected modification on generated source. Hit enter to continue devtool (will overwrite changes!).")
				<-haveLine
				timeLastGenCodeOverwrite = time.Now()
				genLock.Unlock()
			}
			select {
			case <-time.After(100 * time.Millisecond):
			case <-watcher.Event:
			}
		case err := <-watcher.Error:
			if err != nil {
				fmt.Printf("Error watching example ango file: %s\n", err)
				stop(1)
				return
			}
		}
	}
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

func rebuildAngo() {
	stopWg.Add(1)
	defer stopWg.Done()

	cmdRerun := exec.Command("rerun", "--no-run", "--build", pkg.ImportPath)
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
	err = watcher.WatchFlags(`templates/ango-service.tmpl.go`, fsnotify.FSN_MODIFY)
	if err != nil {
		fmt.Printf("Error starting watch on go template file: %s\n", err)
		stop(1)
		return
	}
	err = watcher.WatchFlags(`templates/ango-service.tmpl.js`, fsnotify.FSN_MODIFY)
	if err != nil {
		fmt.Printf("Error starting watch on js template file: %s\n", err)
		stop(1)
		return
	}
	for {
		select {
		case <-stopCh:
			return
		case <-watcher.Event:
			// short timeout consuming other events becuase sublime sometimes saves (modifies) the file twice
			for {
				select {
				case <-time.After(100 * time.Millisecond):
					goto ango
				case <-watcher.Event:
				}
			}
		ango:
			angoExample()
			// short timeout consuming another event because somehow ango modifies the template files!?!?!?!?
			select {
			case <-time.After(100 * time.Millisecond):
			case <-watcher.Event:
			}
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
	genLock.RLock()
	defer genLock.RUnlock()
	timeLastGenCodeOverwrite = time.Now()
	fmt.Println("Running ango tool for example/chatService.ango")
	cmdAngoExample := exec.Command(filepath.Join(wd, "ango"), "--verbose", "-i", exampleAngoFile, "--js-path=example/http-files", "--force-overwrite")
	cmdAngoExample.Stdin = os.Stdin
	cmdAngoExample.Stdout = sgr.NewColorWriter(os.Stdout, sgr.FgBlue, false)
	cmdAngoExample.Stderr = sgr.NewColorWriter(os.Stderr, sgr.FgBlue, false)
	err := cmdAngoExample.Run()
	if err != nil {
		fmt.Printf("Error running ango tool: %s\n", err)
	}
}
