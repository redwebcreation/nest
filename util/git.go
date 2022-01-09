package util

import (
	"fmt"
	"os/exec"
	"strings"
)

type Repository struct {
	Path string
}

func (r Repository) callGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.Path
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err, out)
	}
	return string(out), nil
}

func (r Repository) Clone(url string) error {
	_, err := r.callGit("clone", url, r.Path)

	return err
}

func (r Repository) LatestCommit() (string, error) {
	out, err := r.callGit("rev-parse", "HEAD")

	return strings.TrimSpace(out), err
}

func (r Repository) Checkout(commit string) error {
	_, err := r.callGit("checkout", commit)
	return err
}

func (r Repository) ReadFile(path string) (string, error) {
	return r.callGit("show", "HEAD:"+path)
}

func (r Repository) Files() ([]string, error) {
	out, err := r.callGit("ls-tree", "HEAD", "--name-only", "-r")
	if err != nil {
		return nil, err
	}

	return strings.Split(out, "\n"), nil
}
