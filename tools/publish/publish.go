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
	preseveArtifacts()

	getBranch()

	fmt.Printf("Current version is: %s\n", fullVersion())

	runBuild()

	runRice()

	moveFile()
}

func preseveArtifacts() {
	artifacts := []string{"ango-release", "ango-latest"}
	for _, art := range artifacts {
		err := exec.Command("wget", "https://drone.io/github.com/GeertJohan/ango/files/"+art).Run()
		if err != nil {
			fmt.Printf("Error preserving artifact '%s': %s\n", art, err)
			os.Exit(1)
		}
	}
}

func getBranch() {
	// get current branch
	branchBytes, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		fmt.Printf("Error getting branch: %s\n", err)
		os.Exit(1)
	}
	branch := strings.Trim(string(branchBytes), "\n")

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
	versionHash = strings.Trim(string(hashBytes), "\n")
}

func runBuild() {
	// compile build args
	buildArgs := []string{`build`}
	if isRelease || isLatest {
		buildArgs = append(buildArgs, `-ldflags`)
	}
	if isRelease {
		buildArgs = append(buildArgs, fmt.Sprintf(`-X main.versionNumber %s -X main.versionHash %s`, versionNumber, versionHash))
	}
	if isLatest {
		buildArgs = append(buildArgs, fmt.Sprintf(`-X main.versionHash %s`, versionHash))
	}

	// run build
	buildOut, err := exec.Command("go", buildArgs...).CombinedOutput()
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
