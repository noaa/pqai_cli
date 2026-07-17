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
var apiKeySource = "설정되지 않음"

func cmdConfig(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("사용법: pqai config <set-api-key|show> ...")
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "set-api-key":
		return cmdConfigSetAPIKey(rest)
	case "show":
		return cmdConfigShow(rest)
	default:
		return fmt.Errorf("알 수 없는 config 하위 명령: %s (set-api-key | show)", sub)
	}
}

func cmdConfigSetAPIKey(args []string) error {
	fs := flag.NewFlagSet("config set-api-key", flag.ExitOnError)
	fromDotenv := fs.String("from-dotenv", "", "이 dotenv 파일에서 PQAI_API_KEY를 읽어와 저장")
	pos := parseArgs(fs, args)

	var key string
	if *fromDotenv != "" {
		vals, err := parseEnvFile(*fromDotenv)
		if err != nil {
			return fmt.Errorf("dotenv 파일을 읽을 수 없습니다: %w", err)
		}
		key = vals["PQAI_API_KEY"]
		if key == "" {
			return fmt.Errorf("%s 파일에 PQAI_API_KEY가 없습니다", *fromDotenv)
		}
	} else {
		if len(pos) < 1 {
			return fmt.Errorf("사용법: pqai config set-api-key <token>  또는  pqai config set-api-key --from-dotenv <path>")
		}
		key = pos[0]
	}

	path, err := setGlobalConfigValue("PQAI_API_KEY", key)
	if err != nil {
		return err
	}
	fmt.Printf("저장됨: %s\n", path)
	fmt.Println("이제 어느 폴더에서 pqai를 실행하든 이 토큰이 자동으로 사용됩니다.")
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

	fmt.Printf("설정 파일: %s\n", path)
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		fmt.Println("상태: 파일 없음 (아직 'pqai config set-api-key <token>'을 실행한 적 없음)")
	}
	fmt.Println()

	if key, ok := vals["PQAI_API_KEY"]; ok && key != "" {
		fmt.Printf("전역 설정의 PQAI_API_KEY: %s\n", maskKey(key))
	} else {
		fmt.Println("전역 설정의 PQAI_API_KEY: 없음")
	}
	fmt.Printf("현재 실제로 사용될 값의 출처: %s\n", apiKeySource)
	return nil
}

func maskKey(k string) string {
	if len(k) <= 8 {
		return strings.Repeat("*", len(k))
	}
	return k[:4] + strings.Repeat("*", len(k)-8) + k[len(k)-4:]
}
