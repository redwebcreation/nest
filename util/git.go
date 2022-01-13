package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Repository string

func NewRepository(remote string, path string) (*Repository, error) {
	if out, err := exec.Command("git", "clone", remote, path).CombinedOutput(); err != nil {
		return nil, fmt.Errorf("%s: %s", err, out)
	}

	repo := Repository(path)
	return &repo, nil
}

func OpenRepository(path string) (*Repository, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	repo := Repository(path)

	return &repo, nil
}

func (r Repository) Exec(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = string(r)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err, out)
	}
	return bytes.TrimSpace(out), nil
}

func (r Repository) LatestCommit() ([]byte, error) {
	return r.Exec("rev-parse", "HEAD")
}

func (r Repository) Checkout(commit string) error {
	_, err := r.Exec("checkout", commit)
	return err
}

func (r Repository) Commits() ([]string, error) {
	out, err := r.Exec("log", "--pretty=%H")
	if err != nil {
		return nil, err
	}

	return strings.Split(string(out), "\n"), nil
}

func (r Repository) Pull(branch string) error {
	_, err := r.Exec("checkout", branch)
	if err != nil {
		return err
	}

	_, err = r.Exec("pull")
	return err
}

func (r Repository) Read(path string) ([]byte, error) {
	return r.Exec("show", "HEAD:"+path)
}

func (r Repository) Tree() ([]string, error) {
	out, err := r.Exec("ls-tree", "-r", "--name-only", "HEAD")
	if err != nil {
		return nil, err
	}

	return strings.Split(string(out), "\n"), nil
}
