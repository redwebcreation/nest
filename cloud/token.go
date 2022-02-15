package cloud

import (
	"github.com/redwebcreation/nest/global"
	"os"
)

// GetToken retrieves the token from the environment
func GetToken() (string, error) {

	token, err := os.ReadFile(global.GetCloudTokenFile())
	if err != nil {
		return "", err
	}

	return string(token), nil
}

// SetToken sets the token for the cloud provider
func SetToken(token string) error {
	return os.WriteFile(global.GetCloudTokenFile(), []byte(token), 0600)
}
