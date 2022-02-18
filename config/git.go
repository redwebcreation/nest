package config

import (
	"bytes"
	"fmt"
	"github.com/redwebcreation/nest/loggy"
	"log"
	"os/exec"
	"strings"
)

type Git struct {
	Logger *log.Logger
}

func (g Git) Clone(remote string, local string, branch string) error {
	if branch == "" {
		return fmt.Errorf("branch is empty")
	}

	_, err := g.run("", "clone", "-b", branch, remote, local)

	return err
}

func (g Git) Pull(dir string, branch string) ([]byte, error) {
	return g.run(dir, "pull", "origin", branch)
}

type Commit struct {
	Hash    string
	Message string
}

type CommitList []Commit

func (c CommitList) Hashes() []string {
	var hashes []string
	for _, commit := range c {
		hashes = append(hashes, commit.Hash)
	}
	return hashes
}

func (g Git) ListCommits(dir string, branch string) (CommitList, error) {
	out, err := g.run(dir, "log", "--pretty=%H=%s", "--no-merges", branch)
	if err != nil {
		return nil, err
	}

	var commits []Commit

	for _, line := range strings.Split(string(out), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid Commit line: %s", line)
		}

		commits = append(commits, Commit{
			Hash:    parts[0],
			Message: parts[1],
		})
	}

	return commits, nil
}

func (g Git) ReadFile(dir string, commit string, file string) ([]byte, error) {
	return g.run(dir, "show", commit+":"+file)
}

func (g *Git) run(dir string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()
	out := buf.Bytes()

	g.Logger.Print(loggy.NewEvent(
		loggy.DebugLevel,
		"running git command",
		loggy.Fields{
			"cli": "git " + strings.Join(args, " "),
			"tag": "vcs.run",
		},
	))

	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, out)
	}

	return out, nil
}

func (g Git) Exists(dir, path, commit string) bool {
	_, err := g.ReadFile(dir, commit, path)
	return err == nil
}
