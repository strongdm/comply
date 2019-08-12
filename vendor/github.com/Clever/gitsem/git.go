package main

import (
	"bytes"
	"os/exec"
	"strings"
)

func isRepoClean() (bool, error) {
	cmd := exec.Command("git", "status", "-s")
	result := &bytes.Buffer{}
	cmd.Stdout = result
	if err := cmd.Run(); err != nil {
		return false, err
	}
	return result.String() == "", nil
}

func repoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	result := &bytes.Buffer{}
	cmd.Stdout = result
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(result.String()), nil
}

func addFile(path string) error {
	return exec.Command("git", "add", path).Run()
}

func commit(message string) error {
	return exec.Command("git", "commit", "-m", message).Run()
}

func tag(version string) error {
	return exec.Command("git", "tag", version).Run()
}
