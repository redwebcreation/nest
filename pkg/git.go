package pkg

import (
	"bytes"
	"fmt"
	logger2 "github.com/redwebcreation/nest/pkg/logger"
	"log"
	"os/exec"
	"strings"
)

type GitWrapper struct {
	Logger *log.Logger
}

func (g GitWrapper) Clone(remote string, local string, branch string) error {
	if branch == "" {
		return fmt.Errorf("branch is empty")
	}

	_, err := g.run("", "clone", "-b", branch, remote, local)

	return err
}

func (g GitWrapper) Pull(dir string, branch string) ([]byte, error) {
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

func (g GitWrapper) ListCommits(dir string, branch string) (CommitList, error) {
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

func (g GitWrapper) ReadFile(dir string, commit string, file string) ([]byte, error) {
	return g.run(dir, "show", commit+":"+file)
}

func (g *GitWrapper) run(dir string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()
	out := buf.Bytes()

	g.Logger.Print(logger2.NewEvent(
		logger2.DebugLevel,
		"running git command",
		logger2.Fields{
			"cli": "git " + strings.Join(args, " "),
			"tag": "vcs.run",
		},
	))

	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, out)
	}

	return out, nil
}

func (g GitWrapper) Exists(dir, path, commit string) bool {
	_, err := g.ReadFile(dir, commit, path)
	return err == nil
}
