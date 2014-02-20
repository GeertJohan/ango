package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var regexpVersionBranch = regexp.MustCompile(`^[v][0-9]+.[0-9]+$`)

var (
	isRelease     = false // branch complies with regexpVersionBranch
	isLatest      = false // barnch == "master"
	versionNumber = ""    // version number (v0.1)
	versionHash   = ""    // version sha hash (short)
)

func fullVersion() string {
	if len(versionNumber) > 0 {
		return versionNumber + "-" + versionHash
	}
	return "other-" + versionHash
}

func publishSuffix() string {
	switch true {
	case isRelease:
		return "-release"
	case isLatest:
		return "-latest"
	default:
		return ""
	}
}

func main() {
	getBranch()

	fmt.Printf("Current version is: %s\n", fullVersion())

	runBuild()

	runRice()

	moveFile()
}

func getBranch() {
	// get current branch
	branchBytes, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		fmt.Printf("Error getting branch: %s\n", err)
		os.Exit(1)
	}
	branch := string(branchBytes)

	// check branch
	if branch == "master" {
		versionNumber = "latest"
		isLatest = true
	} else {
		versionNumber = regexpVersionBranch.FindString(branch)
		if len(versionNumber) > 0 {
			isRelease = true
			fmt.Println("This is a release branch.")
		}
	}

	hashBytes, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		fmt.Printf("Error getting sha refspec: %s\n", err)
		os.Exit(1)
	}
	versionHash = string(hashBytes)
}

func runBuild() {
	// compile build string
	buildStr := `build`
	if isRelease {
		buildStr += fmt.Sprintf(` -ldflags "-X main.versionNumber %s -X main.versionHash %s"`, versionNumber, versionHash)
	}
	if isLatest {
		buildStr += fmt.Sprintf(` -ldflags "-X main.versionHash %s"`, versionHash)
	}

	// run build
	buildOut, err := exec.Command("go", strings.Split(buildStr, " ")...).CombinedOutput()
	if err != nil {
		fmt.Printf("Error running build: %s\n%s\n", err, string(buildOut))
		os.Exit(1)
	}
}

func runRice() {
	err := exec.Command("go", "get", "github.com/GeertJohan/go.rice/rice").Run()
	if err != nil {
		fmt.Printf("Error go-getting rice tool: %s\n", err)
		os.Exit(1)
	}

	err = exec.Command("sudo", "apt-get", "install", "zip").Run()
	if err != nil {
		fmt.Printf("Error installing zip: %s\n", err)
		os.Exit(1)
	}

	err = exec.Command("rice", strings.Split("-i github.com/GeertJohan/ango append --exec ango", " ")...).Run()
	if err != nil {
		fmt.Printf("Error running rice: %s\n", err)
		os.Exit(1)
	}
}

func moveFile() {
	suffix := publishSuffix()
	if len(suffix) == 0 {
		return // don't move
	}

	err := os.Rename("ango", "ango"+suffix)
	if err != nil {
		fmt.Printf("Error renaming ango file: %s\n", err)
		os.Exit(1)
	}
}
