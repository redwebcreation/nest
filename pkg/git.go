package pkg

import (
	"bytes"
	"fmt"
	"github.com/redwebcreation/nest/global"
	"os/exec"
	"strings"
)

var Git = &git{}

type git struct{}

func (g git) Clone(remote string, local string, branch string) error {
	if branch == "" {
		return fmt.Errorf("branch is empty")
	}

	_, err := g.run("", "clone", "-b", branch, remote, local)

	return err
}

func (g git) Pull(dir string, branch string) ([]byte, error) {
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

func (g git) ListCommits(dir string, branch string) (CommitList, error) {
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

func (g git) ReadFile(dir string, commit string, file string) ([]byte, error) {
	return g.run(dir, "show", commit+":"+file)
}

func (g *git) run(dir string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()
	out := buf.Bytes()

	global.LogI(
		global.LevelDebug,
		"running git command",
		global.Fields{
			"cli": "git " + strings.Join(args, " "),
			"tag": "vcs.run",
		},
	)

	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, out)
	}

	return out, nil
}

func (g git) Exists(dir, path, commit string) bool {
	_, err := g.ReadFile(dir, commit, path)
	return err == nil
}
