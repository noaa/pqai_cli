package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// parseEnvFile reads KEY=VALUE pairs from a dotenv-style file. Missing files
// yield an empty map with no error.
func parseEnvFile(path string) (map[string]string, error) {
	vals := map[string]string{}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return vals, nil
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		val = strings.Trim(val, `"'`)
		if key == "" {
			continue
		}
		vals[key] = val
	}
	return vals, scanner.Err()
}

// applyEnvFile sets process environment variables from a dotenv-style file,
// without overriding variables already set in the process environment.
func applyEnvFile(path string) {
	vals, err := parseEnvFile(path)
	if err != nil {
		return
	}
	for key, val := range vals {
		if _, set := os.LookupEnv(key); !set {
			os.Setenv(key, val)
		}
	}
}

// loadDotEnv reads KEY=VALUE pairs from a .env file in the current directory
// (if present) and sets them as environment variables, without overriding
// variables already set in the process environment.
func loadDotEnv() {
	applyEnvFile(".env")
}

// loadGlobalConfig reads KEY=VALUE pairs from the user's global pqai config
// file (if present) as a fallback, so pqai works from any directory without
// a project-local .env. Variables already set (from the shell or a local
// .env) take precedence and are left untouched.
func loadGlobalConfig() {
	path, err := configFilePath()
	if err != nil {
		return
	}
	applyEnvFile(path)
}

// configFilePath returns the path to the user's global pqai config file:
//   - Windows: %AppData%\pqai\config.env
//   - Linux/macOS: $XDG_CONFIG_HOME/pqai/config.env, or ~/.config/pqai/config.env
func configFilePath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.env"), nil
}

func configDir() (string, error) {
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("AppData"); appData != "" {
			return filepath.Join(appData, "pqai"), nil
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "AppData", "Roaming", "pqai"), nil
	}

	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "pqai"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "pqai"), nil
}

// setGlobalConfigValue merges key=value into the global pqai config file,
// creating the config directory/file if needed, and writes it with 0600
// permissions since it holds an API token.
func setGlobalConfigValue(key, value string) (string, error) {
	path, err := configFilePath()
	if err != nil {
		return "", err
	}

	vals, err := parseEnvFile(path)
	if err != nil {
		return "", err
	}
	vals[key] = value

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", err
	}

	var b strings.Builder
	for k, v := range vals {
		fmt.Fprintf(&b, "%s=%s\n", k, v)
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o600); err != nil {
		return "", err
	}
	return path, nil
}
