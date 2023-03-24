package system

import (
	"os"
	"path/filepath"

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

func GetEnv(dir string) (map[string]string, error) {

	// get the .env file path
	dotEnvPath := filepath.Join(dir, ".env")

	// Check if the env file exists
	if _, err := os.Stat(dotEnvPath); err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		} else {
			return nil, err
		}
	} else {
		return ReadEnv(dotEnvPath)
	}

}

func SetDefaultEnv(apiKeys []string, env map[string]string) map[string]string {
	// Ensure all apiKeys are present in env
	for _, apiKey := range apiKeys {
		if _, ok := env[apiKey]; !ok {
			// If the environment variable is set in the system, use its value; otherwise, use the empty string
			envValue, ok := os.LookupEnv(apiKey)
			if !ok {
				envValue = ""
			}
			env[apiKey] = envValue
		}
	}

	return env
}
