package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var version = "dev"

const usage = `pqai - PQAI API 커맨드라인 클라이언트 (https://api.projectpq.ai)

사용법:
  pqai <command> [flags] [args]

검색 (토큰 필요):
  search <query>          텍스트 쿼리로 선행기술 문서 검색 (102)
  combos <query>          텍스트 쿼리로 선행기술 조합 검색 (103)
  prior-art <pn>          특허의 선행기술 검색 (출원일 이전 문서)
  similar <pn>            특허와 유사한 문서 검색

문서/특허 (토큰 필요):
  patent <pn>             특허 데이터 조회
  document <id>           PQAI 데이터베이스에서 문서 조회
  vector <pn> <field>     특허 벡터 조회 (field: cpcs | abstract)
  snippet <pn> -q <text>  쿼리-문서 쌍의 스니펫 조회
  mapping <pn> -q <text>  쿼리-문서 쌍의 요소별 매핑 조회
  dataset                 데이터셋 샘플 조회 (-name, -n)

도면:
  drawings <pn>           특허 도면 목록 조회 (토큰 필요)
  drawing <pn> <n>        특허 도면 다운로드 (PNG, 토큰 불필요)

분류 (토큰 필요):
  cpcs <text>             텍스트에 대한 CPC 분류 제안
  gaus <text>             텍스트에 대한 Group Art Unit 제안
  cpc-def <cpc>           CPC 클래스 정의 조회

설정:
  config set-api-key <token>          API 토큰을 전역 설정 파일에 저장 (모든 폴더에서 사용됨)
  config set-api-key --from-dotenv <path>   dotenv 파일에서 토큰을 읽어와 전역 설정에 저장
  config show                         전역 설정 파일 위치와 현재 사용 중인 토큰의 출처 확인

환경변수:
  PQAI_API_KEY            API 액세스 토큰 (필수, 도면 다운로드 라우트 제외)
  PQAI_ENDPOINT           API 주소 재정의 (기본: https://api.projectpq.ai)

토큰은 다음 순서로 우선 적용됩니다 (위가 우선):
  1. 셸에서 export한 PQAI_API_KEY
  2. 현재 폴더의 .env 파일
  3. 'pqai config set-api-key'로 저장한 전역 설정 파일 (모든 폴더에서 동작)

각 명령의 플래그는 'pqai <command> -h'로 확인하세요.

  version                 CLI 버전 출력
`

func main() {
	_, hadShellEnv := os.LookupEnv("PQAI_API_KEY")
	loadDotEnv()
	_, hadEnvAfterDotenv := os.LookupEnv("PQAI_API_KEY")
	loadGlobalConfig()

	switch {
	case hadShellEnv:
		apiKeySource = "셸 환경변수 (export PQAI_API_KEY=...)"
	case hadEnvAfterDotenv:
		apiKeySource = "현재 폴더의 .env 파일"
	default:
		if v, ok := os.LookupEnv("PQAI_API_KEY"); ok && v != "" {
			apiKeySource = "전역 설정 파일 (pqai config set-api-key)"
		}
	}

	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	c := NewClient()
	cmd, args := os.Args[1], os.Args[2:]

	var err error
	switch cmd {
	case "search":
		err = cmdSearch(c, "/search/102/", args)
	case "combos":
		err = cmdSearch(c, "/search/103/", args)
	case "prior-art":
		err = cmdPatentSearch(c, "/prior-art/patent/", args)
	case "similar":
		err = cmdPatentSearch(c, "/similar/", args)
	case "snippet":
		err = cmdSnippet(c, "/snippets/", args)
	case "mapping":
		err = cmdSnippet(c, "/mappings/", args)
	case "document":
		err = cmdDocument(c, args)
	case "patent":
		err = cmdPatent(c, args)
	case "vector":
		err = cmdVector(c, args)
	case "dataset":
		err = cmdDataset(c, args)
	case "drawings":
		err = cmdDrawings(c, args)
	case "drawing":
		err = cmdDrawing(c, args)
	case "cpcs":
		err = cmdText(c, "/suggest/cpcs", "text", args)
	case "gaus":
		err = cmdText(c, "/predict/gaus", "text", args)
	case "cpc-def":
		err = cmdText(c, "/definitions/cpcs", "cpc", args)
	case "config":
		err = cmdConfig(args)
	case "help", "-h", "--help":
		fmt.Print(usage)
	case "version", "-v", "--version":
		fmt.Println("pqai " + version)
	default:
		fmt.Fprintf(os.Stderr, "알 수 없는 명령: %s\n\n", cmd)
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "오류:", err)
		os.Exit(1)
	}
}

func cmdSearch(c *Client, route string, args []string) error {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	n := fs.Int("n", 10, "결과 개수")
	offset := fs.Int("offset", 0, "페이지네이션 오프셋 (0부터)")
	index := fs.String("index", "", "CPC 서브클래스 (예: H04W, auto)")
	cc := fs.String("cc", "", "국가 코드 목록 (예: US,EP,WO)")
	dtype := fs.String("dtype", "", "컷오프 날짜 기준: priority|publication|filing")
	after := fs.String("after", "", "컷오프 시작 날짜 (예: 2016-01-01)")
	before := fs.String("before", "", "컷오프 종료 날짜 (예: 2019-12-31)")
	typ := fs.String("type", "", "문서 유형: patent|npl")
	snip := fs.Bool("snip", false, "스니펫 포함")
	maps := fs.Bool("maps", false, "요소별 매핑 포함")
	lq := fs.String("lq", "", `잠재 쿼리 JSON (예: {"relevant":[],"irrelevant":[]})`)
	asJSON := fs.Bool("json", false, "원본 JSON 출력")
	pos := parseArgs(fs, args)

	if len(pos) < 1 {
		return fmt.Errorf("검색 쿼리를 입력하세요")
	}
	params := url.Values{"q": {pos[0]}, "n": {strconv.Itoa(*n)}}
	if *offset > 0 {
		params.Set("offset", strconv.Itoa(*offset))
	}
	setIf(params, "index", *index)
	setIf(params, "cc", *cc)
	setIf(params, "dtype", *dtype)
	setIf(params, "after", *after)
	setIf(params, "before", *before)
	setIf(params, "type", *typ)
	setIf(params, "lq", *lq)
	if *snip {
		params.Set("snip", "1")
	}
	if *maps {
		params.Set("maps", "1")
	}

	body, err := c.Get(route, params, true)
	if err != nil {
		return err
	}
	if *asJSON {
		printJSON(body)
	} else {
		printSearchResults(body)
	}
	return nil
}

func cmdPatentSearch(c *Client, route string, args []string) error {
	fs := flag.NewFlagSet("patent-search", flag.ExitOnError)
	n := fs.Int("n", 10, "결과 개수")
	offset := fs.Int("offset", 0, "페이지네이션 오프셋 (0부터)")
	index := fs.String("index", "", "CPC 서브클래스 (예: H04W, auto)")
	typ := fs.String("type", "", "문서 유형: patent|npl")
	asJSON := fs.Bool("json", false, "원본 JSON 출력")
	pos := parseArgs(fs, args)

	if len(pos) < 1 {
		return fmt.Errorf("특허 번호를 입력하세요 (예: US7654321B2)")
	}
	params := url.Values{"pn": {pos[0]}, "n": {strconv.Itoa(*n)}}
	if *offset > 0 {
		params.Set("offset", strconv.Itoa(*offset))
	}
	setIf(params, "index", *index)
	setIf(params, "type", *typ)

	body, err := c.Get(route, params, true)
	if err != nil {
		return err
	}
	if *asJSON {
		printJSON(body)
	} else {
		printSearchResults(body)
	}
	return nil
}

func cmdSnippet(c *Client, route string, args []string) error {
	fs := flag.NewFlagSet("snippet", flag.ExitOnError)
	q := fs.String("q", "", "텍스트 쿼리")
	pos := parseArgs(fs, args)

	if len(pos) < 1 || *q == "" {
		return fmt.Errorf("사용법: pqai snippet|mapping <pn> -q <query>")
	}
	params := url.Values{"q": {*q}, "pn": {pos[0]}}
	body, err := c.Get(route, params, true)
	if err != nil {
		return err
	}
	printJSON(body)
	return nil
}

func cmdDocument(c *Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("문서 ID를 입력하세요 (예: US7654321B2)")
	}
	body, err := c.Get("/documents/", url.Values{"id": {args[0]}}, true)
	if err != nil {
		return err
	}
	printJSON(body)
	return nil
}

func cmdPatent(c *Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("특허 번호를 입력하세요 (예: US7654321B2)")
	}
	body, err := c.Get("/patents/"+url.PathEscape(args[0]), nil, true)
	if err != nil {
		return err
	}
	printJSON(body)
	return nil
}

func cmdVector(c *Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("사용법: pqai vector <pn> <field>  (field: cpcs | abstract)")
	}
	route := "/patents/" + url.PathEscape(args[0]) + "/vectors/" + url.PathEscape(args[1])
	body, err := c.Get(route, nil, true)
	if err != nil {
		return err
	}
	printJSON(body)
	return nil
}

func cmdDataset(c *Client, args []string) error {
	fs := flag.NewFlagSet("dataset", flag.ExitOnError)
	name := fs.String("name", "PoC", "데이터셋 이름")
	n := fs.Int("n", 0, "샘플 번호")
	parseArgs(fs, args)

	params := url.Values{"dataset": {*name}, "n": {strconv.Itoa(*n)}}
	body, err := c.Get("/datasets/", params, true)
	if err != nil {
		return err
	}
	printJSON(body)
	return nil
}

func cmdDrawings(c *Client, args []string) error {
	fs := flag.NewFlagSet("drawings", flag.ExitOnError)
	thumb := fs.Bool("thumb", false, "썸네일 목록 조회")
	pos := parseArgs(fs, args)

	if len(pos) < 1 {
		return fmt.Errorf("특허 번호를 입력하세요 (예: US7654321B2)")
	}
	kind := "drawings"
	if *thumb {
		kind = "thumbnails"
	}
	body, err := c.Get("/patents/"+url.PathEscape(pos[0])+"/"+kind, nil, true)
	if err != nil {
		return err
	}
	printJSON(body)
	return nil
}

func cmdDrawing(c *Client, args []string) error {
	fs := flag.NewFlagSet("drawing", flag.ExitOnError)
	thumb := fs.Bool("thumb", false, "썸네일로 다운로드")
	w := fs.Int("w", 0, "썸네일 너비 (px)")
	h := fs.Int("h", 0, "썸네일 높이 (px)")
	out := fs.String("o", "", "저장할 파일 경로 (기본: <pn>_<n>.png)")
	pos := parseArgs(fs, args)

	if len(pos) < 2 {
		return fmt.Errorf("사용법: pqai drawing <pn> <n> [-thumb] [-w 300] [-o out.png]")
	}
	pn, n := pos[0], pos[1]
	kind := "drawings"
	if *thumb || *w > 0 || *h > 0 {
		kind = "thumbnails"
	}
	params := url.Values{}
	if *w > 0 {
		params.Set("w", strconv.Itoa(*w))
	}
	if *h > 0 {
		params.Set("h", strconv.Itoa(*h))
	}
	route := "/patents/" + url.PathEscape(pn) + "/" + kind + "/" + url.PathEscape(n)
	body, err := c.Get(route, params, false)
	if err != nil {
		return err
	}
	path := *out
	if path == "" {
		path = fmt.Sprintf("%s_%s.png", pn, n)
	}
	if err := os.WriteFile(path, body, 0o644); err != nil {
		return err
	}
	fmt.Printf("저장됨: %s (%d bytes)\n", path, len(body))
	return nil
}

func cmdText(c *Client, route, key string, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("텍스트를 입력하세요")
	}
	body, err := c.Get(route, url.Values{key: {args[0]}}, true)
	if err != nil {
		return err
	}
	printJSON(body)
	return nil
}

// parseArgs allows flags and positional arguments to appear in any order,
// which Go's flag package does not support natively (it stops parsing flags
// at the first positional argument). It returns the positional arguments.
func parseArgs(fs *flag.FlagSet, args []string) []string {
	var flagArgs, positional []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			positional = append(positional, args[i+1:]...)
			break
		}
		if len(a) > 1 && a[0] == '-' {
			flagArgs = append(flagArgs, a)
			name := strings.TrimLeft(a, "-")
			if strings.Contains(name, "=") {
				continue // value embedded, no extra token to consume
			}
			if f := fs.Lookup(name); f != nil {
				if bv, ok := f.Value.(interface{ IsBoolFlag() bool }); ok && bv.IsBoolFlag() {
					continue // boolean flags don't consume the next token
				}
			}
			if i+1 < len(args) {
				i++
				flagArgs = append(flagArgs, args[i])
			}
			continue
		}
		positional = append(positional, a)
	}
	fs.Parse(flagArgs)
	return positional
}

func setIf(params url.Values, key, val string) {
	if val != "" {
		params.Set(key, val)
	}
}
