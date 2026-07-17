package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// apiKeySource records where PQAI_API_KEY ended up coming from, determined
// once in main() before the command dispatch, since by the time any command
// runs, loadDotEnv/loadGlobalConfig have already merged everything into the
// process environment and the original source is no longer distinguishable.
var apiKeySource = "not set"

func cmdConfig(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: pqai config <set-api-key|show> ...")
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "set-api-key":
		return cmdConfigSetAPIKey(rest)
	case "show":
		return cmdConfigShow(rest)
	default:
		return fmt.Errorf("unknown config subcommand: %s (set-api-key | show)", sub)
	}
}

func cmdConfigSetAPIKey(args []string) error {
	fs := flag.NewFlagSet("config set-api-key", flag.ExitOnError)
	fromDotenv := fs.String("from-dotenv", "", "read PQAI_API_KEY from this dotenv file and save it")
	pos := parseArgs(fs, args)

	var key string
	if *fromDotenv != "" {
		vals, err := parseEnvFile(*fromDotenv)
		if err != nil {
			return fmt.Errorf("could not read dotenv file: %w", err)
		}
		key = vals["PQAI_API_KEY"]
		if key == "" {
			return fmt.Errorf("%s does not contain PQAI_API_KEY", *fromDotenv)
		}
	} else {
		if len(pos) < 1 {
			return fmt.Errorf("usage: pqai config set-api-key <token>  or  pqai config set-api-key --from-dotenv <path>")
		}
		key = pos[0]
	}

	path, err := setGlobalConfigValue("PQAI_API_KEY", key)
	if err != nil {
		return err
	}
	fmt.Printf("Saved: %s\n", path)
	fmt.Println("pqai will now use this token automatically from any folder.")
	return nil
}

func cmdConfigShow(args []string) error {
	path, err := configFilePath()
	if err != nil {
		return err
	}
	vals, err := parseEnvFile(path)
	if err != nil {
		return err
	}

	fmt.Printf("Config file: %s\n", path)
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		fmt.Println("Status: file does not exist yet ('pqai config set-api-key <token>' has not been run)")
	}
	fmt.Println()

	if key, ok := vals["PQAI_API_KEY"]; ok && key != "" {
		fmt.Printf("Global config PQAI_API_KEY: %s\n", maskKey(key))
	} else {
		fmt.Println("Global config PQAI_API_KEY: none")
	}
	fmt.Printf("Source of the value actually in use: %s\n", apiKeySource)
	return nil
}

func maskKey(k string) string {
	if len(k) <= 8 {
		return strings.Repeat("*", len(k))
	}
	return k[:4] + strings.Repeat("*", len(k)-8) + k[len(k)-4:]
}
