package docker

import (
	"encoding/base64"
	"encoding/json"
)

// Registry represents a docker registry
type Registry struct {
	// Name of the registry.
	Name string `yaml:"name"`
	// Host of the registry.
	Host string `yaml:"host"`
	// Port of the registry.
	Port string `yaml:"port"`
	// Username to use when authenticating with the registry.
	Username string `yaml:"username"`
	// Password to use when authenticating with the registry.
	Password string `yaml:"password"`
}

// IsZero returns true if the registry is empty
func (r Registry) IsZero() bool {
	return r.Name == "" && r.Host == "" && r.Port == "" && r.Username == "" && r.Password == ""
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
