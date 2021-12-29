package global

import (
	"encoding/json"
	"fmt"
	"github.com/me/nest/util"
	"github.com/mitchellh/go-homedir"
	"os"
)

var ResolverConfigFile string
var IsConfigResolverConfigured bool
var ConfigResolver *Resolver

type Resolver struct {
	Strategy      string `json:"strategy,omitempty"`
	Provider      string `json:"provider,omitempty"`
	TransportMode string `json:"transportMode,omitempty"`
	Repository    string `json:"repository,omitempty"`
	fs            util.RoFS
}

func (r Resolver) Write() error {
	contents, err := json.Marshal(r)
	if err != nil {
		return err
	}

	return os.WriteFile(ResolverConfigFile, contents, 0644)
}

func (r *Resolver) loadFilesystem() error {
	if r.fs != nil {
		return nil
	}

	if r.Strategy != "remote" {
		return fmt.Errorf("invalid strategy: only remote is supported")
	}

	if r.Provider != "GitHub" {
		return fmt.Errorf("invalid provider: only GitHub is supported")
	}

	remote := &util.Remote{
		Url: "git@github.com:" + r.Repository,
	}

	err := remote.Load()
	if err != nil {
		return err
	}

	r.fs = remote
	return nil
}

func (r Resolver) Get(path string) ([]byte, error) {
	err := r.loadFilesystem()
	if err != nil {
		return nil, err
	}

	return r.fs.Get(path)
}

func (r Resolver) Exists(path string) bool {
	err := r.loadFilesystem()
	if err != nil {
		return false
	}

	return r.fs.Exists(path)
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	ResolverConfigFile = home + "/.nest.json"

	_, err = os.Stat(ResolverConfigFile)
	IsConfigResolverConfigured = err == nil

	if IsConfigResolverConfigured {
		contents, err := os.ReadFile(ResolverConfigFile)
		if err != nil {
			panic(err)
		}

		var f Resolver
		err = json.Unmarshal(contents, &f)
		if err != nil {
			panic(err)
		}

		ConfigResolver = &f
	} else {
		ConfigResolver = &Resolver{}
	}
}
