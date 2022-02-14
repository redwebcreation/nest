package cloud

import (
	"github.com/mitchellh/go-homedir"
	"os"
)

// GetToken retrieves the token from the environment
func GetToken() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	token, err := os.ReadFile(home + "/.nest/.cloudtoken")
	if err != nil {
		return "", err
	}

	return string(token), nil
}

// SetToken sets the token for the cloud provider
func SetToken(token string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	return os.WriteFile(home+"/.nest/.cloudtoken", []byte(token), 0600)
}
