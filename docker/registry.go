package docker

import (
	"encoding/base64"
	"encoding/json"
)

// Registry represents a docker registry
type Registry struct {
	// Name of the registry.
	Name string `yaml:"name" json:"name"`
	// Host of the registry.
	Host string `yaml:"host" json:"host"`
	// Port of the registry.
	Port string `yaml:"port" json:"port"`
	// Username to use when authenticating with the registry.
	Username string `yaml:"username" json:"username"`
	// Password to use when authenticating with the registry.
	Password string `yaml:"password" json:"password"`
}

// UrlFor returns the url for the registry
func (r Registry) UrlFor(image string) string {
	if r.Port != "" {
		return r.Host + ":" + r.Port + "/" + image
	}

	return r.Host + "/" + image
}

// ToBase64 returns the auth string for the registry
func (r Registry) ToBase64() (string, error) {
	auth, err := json.Marshal(map[string]string{
		"username": r.Username,
		"password": r.Password,
	})

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(auth), nil
}
