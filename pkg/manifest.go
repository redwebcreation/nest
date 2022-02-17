package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/redwebcreation/nest/global"
	"os"
)

var (
	ErrManifestNotFound = fmt.Errorf("manifest not found")
)

type Manifest struct {
	ID         string            `json:"id"`
	Containers map[string]string `json:"containers"`
	Networks   map[string]string `json:"networks"`
}

func (m Manifest) Save() error {
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return os.WriteFile(global.GetContainerManifestFile(m.ID), bytes, 0600)
}

func NewManifest(id string) *Manifest {
	return &Manifest{
		ID:         id,
		Containers: make(map[string]string),
		Networks:   make(map[string]string),
	}
}

func LoadManifest(id string) (*Manifest, error) {
	bytes, err := os.ReadFile(global.GetContainerManifestFile(id))
	if err != nil {
		return nil, err
	}

	var m Manifest
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func GetLatestManifest() (*Manifest, error) {
	manifests, err := os.ReadDir(global.GetManifestsDir())
	if err != nil {
		return nil, err
	}

	if len(manifests) == 0 {
		return nil, ErrManifestNotFound
	}

	latest := manifests[len(manifests)-1].Name()

	// removes .json
	return LoadManifest(latest[:len(latest)-5])
}
