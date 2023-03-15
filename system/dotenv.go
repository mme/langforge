package system

import (
	"os"

	"github.com/joho/godotenv"
)

func ReadEnv(path string) (map[string]string, error) {
	return godotenv.Read(path)
}

func WriteEnv(path string, env map[string]string) error {
	return godotenv.Write(env, path)
}

func EnsureEnv(path string, keys []string) error {
	env := make(map[string]string)

	// Check if the env file exists
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		// Read the env file if it exists
		env, err = ReadEnv(path)
		if err != nil {
			return err
		}
	}

	// Ensure all keys are present
	changed := false
	for _, key := range keys {
		if _, ok := env[key]; !ok {
			envValue, ok := os.LookupEnv(key)
			if !ok {
				envValue = ""
			}
			env[key] = envValue
			changed = true
		}
	}

	// If there were changes, write the updated env back to the file
	if changed {
		err := WriteEnv(path, env)
		if err != nil {
			return err
		}
	}

	return nil
}

func UnsetAPIKeys(path string, keys []string) ([]string, error) {
	// Load the environment from the .env file
	env, err := godotenv.Read(path)
	if err != nil {
		return nil, err
	}

	unsetKeys := []string{}
	for _, key := range keys {
		value, ok := env[key]
		if !ok || value == "" {
			unsetKeys = append(unsetKeys, key)
		}
	}

	return unsetKeys, nil
}
