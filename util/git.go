package util

import (
	"bytes"
	"fmt"
	"github.com/redwebcreation/nest/global"
	"os/exec"
	"strings"
)

type VCS struct {
	Cmd string // name of the binary to invoke

	CloneCmd       string
	PullCmd        string
	ListCommitsCmd string
	ReadFileCmd    string
}

var VcsGit = &VCS{
	Cmd: "git",

	CloneCmd:       "clone {remote} {local}",
	PullCmd:        "pull origin {branch}", // todo: remote name is hardcoded
	ListCommitsCmd: "log --pretty='%H@%s' --no-merges {branch}",
	ReadFileCmd:    "show {commit}:{file}",
}

func (v VCS) Clone(remote string, local string) error {
	_, err := v.run("", v.CloneCmd, "remote", remote, "local", local)
	return err
}

func (v VCS) Pull(dir string, branch string) ([]byte, error) {
	return v.run(dir, v.PullCmd, "branch", branch)
}

type Commit struct {
	Hash    string
	Message string
}

func (v VCS) ListCommits(dir string, branch string) ([]Commit, error) {
	out, err := v.run(dir, v.ListCommitsCmd, "branch", branch)
	if err != nil {
		return nil, err
	}

	var commits []Commit

	for _, line := range strings.Split(string(out), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(strings.Trim(line, "'"), "@")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid commit line: %s", line)
		}

		commits = append(commits, Commit{
			Hash:    parts[0],
			Message: parts[1],
		})
	}

	return commits, nil
}

func (v VCS) ReadFile(dir string, commit string, file string) ([]byte, error) {
	return v.run(dir, v.ReadFileCmd, "commit", commit, "file", file)
}

func (v *VCS) run(dir string, cmdline string, keyval ...string) ([]byte, error) {
	m := make(map[string]string)
	for i := 0; i < len(keyval); i += 2 {
		m[keyval[i]] = keyval[i+1]
	}

	args := strings.Fields(cmdline)

	for i, arg := range args {
		args[i] = expand(arg, m)
	}

	_, err := exec.LookPath(v.Cmd)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(v.Cmd, args...)
	cmd.Dir = dir

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err = cmd.Run()
	out := buf.Bytes()

	global.InternalLogger.Log(
		global.LevelDebug,
		"running git command",
		global.Fields{
			"cmd": v.Cmd + " " + strings.Join(args, " "),
			"tag": "vcs.run",
		},
	)

	if err != nil {
		return nil, fmt.Errorf("%v: %s", err, out)
	}

	return out, nil
}

func (v VCS) Exists(dir, path, commit string) bool {
	_, err := v.ReadFile(dir, commit, path)
	return err == nil
}

func expand(s string, bindings map[string]string) string {
	for k, v := range bindings {
		s = strings.Replace(s, "{"+k+"}", v, -1)
	}

	return s
}
