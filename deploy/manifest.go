package deploy

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	ErrNotFound = fmt.Errorf("manifest not found")
)

type Manifest struct {
	ID         string            `json:"id"`
	Containers map[string]string `json:"containers"`
	Networks   map[string]string `json:"networks"`
}

// Manager contains the path to the manifest file and methods to manage manifests.
type Manager struct {
	Path string
}

func (m Manager) NewManifest(id string) *Manifest {
	return &Manifest{
		ID:         id,
		Containers: make(map[string]string),
		Networks:   make(map[string]string),
	}
}

func (m Manager) LoadWithID(path string) (*Manifest, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	err = json.Unmarshal(bytes, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}

func (m Manager) Latest() (*Manifest, error) {
	manifests, err := os.ReadDir(m.Path)
	if err != nil {
		return nil, err
	}

	if len(manifests) == 0 {
		return nil, ErrNotFound
	}

	latest := manifests[len(manifests)-1].Name()

	// removes .json
	return m.LoadWithID(latest[:len(latest)-5])
}

func (m Manager) Save(manifest *Manifest) error {
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return os.WriteFile(m.Path+"/"+manifest.ID+".json", bytes, 0600)
}
