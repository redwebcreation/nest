package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Repository interface {
	Exec(...string) ([]byte, error)
	LatestCommit() (string, error)
	Checkout(string) error
	Commits() ([]string, error)
	Pull(branch string) ([]byte, error)
	Read(string) ([]byte, error)
	Tree() ([]string, error)
}

type GitRepository string

func NewRepository(remote string, path string) (*GitRepository, error) {
	if out, err := exec.Command("git", "clone", remote, path).CombinedOutput(); err != nil {
		// split out by newline
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			fmt.Fprintf(os.Stderr, "\n  "+Gray.Fg()+"|  "+line+Reset())
		}

		fmt.Fprintln(os.Stderr)

		return nil, fmt.Errorf("could not clone the configuration")
	}

	repo := GitRepository(path)
	return &repo, nil
}

func OpenRepository(path string) (*GitRepository, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	repo := GitRepository(path)

	return &repo, nil
}

func (r GitRepository) Exec(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = string(r)
	out, err := cmd.CombinedOutput()

	return bytes.TrimSpace(out), err
}

func (r GitRepository) LatestCommit() (string, error) {
	out, err := r.Exec("rev-parse", "HEAD")
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(out)), nil
}

func (r GitRepository) Checkout(commit string) error {
	_, err := r.Exec("checkout", commit)
	return err
}

func (r GitRepository) Commits() ([]string, error) {
	out, err := r.Exec("log", "--pretty=%H")
	if err != nil {
		return nil, err
	}

	return strings.Split(string(out), "\n"), nil
}

func (r GitRepository) Pull(branch string) ([]byte, error) {
	_, err := r.Exec("checkout", branch)
	if err != nil {
		return nil, err
	}

	out, err := r.Exec("pull")
	return out, err
}

func (r GitRepository) Read(path string) ([]byte, error) {
	return r.Exec("show", "HEAD:"+path)
}

func (r GitRepository) Tree() ([]string, error) {
	out, err := r.Exec("ls-tree", "-r", "--name-only", "HEAD")
	if err != nil {
		return nil, err
	}

	return strings.Split(string(out), "\n"), nil
}
