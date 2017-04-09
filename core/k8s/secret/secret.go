package secret

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rutmir/services/core/log"
)

const defaultSecretsPath = "/etc/secrets"

var (
	secretsPath string
)

func init() {
	secretsPath = os.Getenv("SECRETS_PATH")
	if len(secretsPath) == 0 {
		secretsPath = defaultSecretsPath
	}
}

// Loads a value from the specified secret file
func getValue(name string, panic bool) string {
	path := filepath.Join(secretsPath, name)
	secret, err := ioutil.ReadFile(path)
	if err != nil {
		if panic {
			log.Fatal("unable to read secret - %s", path)
		} else {
			log.Err("unable to read secret - %s", path)
		}
	}
	return string(secret)
}

// GetValue
func GetValue(name string) string {
	return getValue(name, false)
}

// GetValueOrPanic
func GetValueOrPanic(name string) string {
	return getValue(name, true)
}

// GetSecretByFilePath
func GetSecretByFilePath(name string) (string, bool) {
	path := filepath.Join(secretsPath, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", false
	}

	return path, true
}
