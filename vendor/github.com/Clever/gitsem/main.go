package main

import (
	"flag"
	"fmt"
	"gopkg.in/blang/semver.v1"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func commitMessage(message, version string) string {
	if strings.Contains(message, "%s") {
		return fmt.Sprintf(message, version)
	}
	return message
}

func getCurrentVersion(path string) (*semver.Version, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &semver.Version{}, nil
	}
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return semver.New(strings.TrimSpace(string(contents)))
}

const versionFileName = "VERSION"

func exitWithError(message string) {
	fmt.Fprintf(os.Stderr, message+"\n\n")
	flag.Usage()
	os.Exit(1)
}

func bump(old *semver.Version, part string) *semver.Version {
	// We don't want to mutate the input, but there's no Clone or Copy method on a semver.Version,
	// so we make a new one by parsing the string version of the old one.
	// We ignore any errors because we know it's valid semver.
	new, _ := semver.New(old.String())
	switch part {
	case "major":
		new.Major++
		new.Minor = 0
		new.Patch = 0
	case "minor":
		new.Minor++
		new.Patch = 0
	case "patch":
		new.Patch++
	}
	return new
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [options] version\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "version can be one of: newversion | patch | minor | major\n\n")
		fmt.Fprintf(os.Stderr, "options:\n")
		flag.PrintDefaults()
	}
	message := flag.String("m", "%s", "commit message for version commit")
	help := flag.Bool("h", false, "print usage and exit")
	shouldTag := flag.Bool("tag", true, "whether or not to make a tag at the version commit")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *message == "" {
		exitWithError("missing message")
	}

	if clean, err := isRepoClean(); err != nil {
		log.Fatal(err)
	} else if !clean {
		log.Fatal("repo isn't clean")
	}

	root, err := repoRoot()
	if err != nil {
		log.Fatal(err)
	}
	versionFile := filepath.Join(root, versionFileName)
	version, err := getCurrentVersion(versionFile)
	if err != nil {
		log.Fatal(err)
	}
	if len(flag.Args()) != 1 {
		exitWithError("gitsem takes exactly one non-flag argument: version")
	}

	newVersion := flag.Args()[0]
	switch newVersion {
	case "patch", "minor", "major":
		version = bump(version, newVersion)
	default:
		if version, err = semver.New(newVersion); err != nil {
			log.Fatalf("failed to parse %s as semver: %s", newVersion, err.Error())
		}
	}

	if err := ioutil.WriteFile(versionFile, []byte(version.String()), 0666); err != nil {
		log.Fatal(err)
	}
	if err := addFile(versionFile); err != nil {
		log.Fatal(err)
	}
	versionString := "v" + version.String()
	*message = commitMessage(*message, versionString)
	if err := commit(*message); err != nil {
		log.Fatal(err)
	}
	if *shouldTag {
		if err := tag(versionString); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(versionString)
}
