package util

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"io"
)

// RoFS stands for Read-Only File System.
type RoFS interface {
	Get(string) ([]byte, error)
	Exists(string) bool
}

type Remote struct {
	Url  string
	fs   billy.Filesystem
	repo *git.Repository
}

func (r *Remote) Load() error {
	if r.fs != nil {
		return nil
	}

	fs := memfs.New()

	repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:   r.Url,
		Depth: 1,
	})

	if err != nil {
		return err
	}

	r.repo = repo
	r.fs = fs

	return nil
}

func (r *Remote) Get(path string) ([]byte, error) {
	err := r.Load()
	if err != nil {
		return nil, err
	}

	f, err := r.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}

func (r *Remote) Exists(path string) bool {
	err := r.Load()
	if err != nil {
		return false
	}

	_, err = r.fs.Stat(path)
	return err == nil
}
