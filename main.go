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

const usage = `pqai - PQAI API command-line client (https://api.projectpq.ai)

Usage:
  pqai <command> [flags] [args]

Search (token required):
  search <query>          search prior-art documents by text query (102)
  combos <query>          search prior-art "combinations" by text query (103)
  prior-art <pn>          search prior art for a patent (documents predating its filing)
  similar <pn>            search documents similar to a patent

Documents/patents (token required):
  patent <pn>             look up patent data
  document <id>           look up a document in the PQAI database
  vector <pn> <field>     look up a patent's embedding vector (field: cpcs | abstract)
  snippet <pn> -q <text>  look up the snippet for a query-document pair
  mapping <pn> -q <text>  look up the per-element mapping for a query-document pair
  dataset                 fetch a dataset sample (-name, -n)

Drawings:
  drawings <pn>           list a patent's drawings (token required)
  drawing <pn> <n>        download a patent drawing (PNG, no token required)

Classification (token required):
  cpcs <text>             suggest CPC classifications for text
  gaus <text>             suggest a Group Art Unit for text
  cpc-def <cpc>           look up a CPC class definition

Config:
  config set-api-key <token>                 save the API token to the global config file (used from any folder)
  config set-api-key --from-dotenv <path>    read the token from a dotenv file and save it to the global config
  config show                                show the global config file location and where the active token came from

Environment variables:
  PQAI_API_KEY            API access token (required, except for the drawing-download route)
  PQAI_ENDPOINT           override the API base URL (default: https://api.projectpq.ai)

The token is resolved in this order (highest priority first):
  1. PQAI_API_KEY exported in your shell
  2. A .env file in the current folder
  3. The global config file saved via 'pqai config set-api-key' (works from any folder)

Run 'pqai <command> -h' to see flags for each command.

  version                 print the CLI version
`

func main() {
	_, hadShellEnv := os.LookupEnv("PQAI_API_KEY")
	loadDotEnv()
	_, hadEnvAfterDotenv := os.LookupEnv("PQAI_API_KEY")
	loadGlobalConfig()

	switch {
	case hadShellEnv:
		apiKeySource = "shell environment variable (export PQAI_API_KEY=...)"
	case hadEnvAfterDotenv:
		apiKeySource = "local .env file in the current folder"
	default:
		if v, ok := os.LookupEnv("PQAI_API_KEY"); ok && v != "" {
			apiKeySource = "global config file (pqai config set-api-key)"
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
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func cmdSearch(c *Client, route string, args []string) error {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	n := fs.Int("n", 10, "number of results")
	offset := fs.Int("offset", 0, "pagination offset (0-based)")
	index := fs.String("index", "", "CPC subclass (e.g. H04W, auto)")
	cc := fs.String("cc", "", "country code list (e.g. US,EP,WO)")
	dtype := fs.String("dtype", "", "cutoff-date basis: priority|publication|filing")
	after := fs.String("after", "", "cutoff start date (e.g. 2016-01-01)")
	before := fs.String("before", "", "cutoff end date (e.g. 2019-12-31)")
	typ := fs.String("type", "", "document type: patent|npl")
	snip := fs.Bool("snip", false, "include a matching snippet")
	maps := fs.Bool("maps", false, "include per-element mapping")
	lq := fs.String("lq", "", `latent query JSON (e.g. {"relevant":[],"irrelevant":[]})`)
	asJSON := fs.Bool("json", false, "print raw JSON output")
	pos := parseArgs(fs, args)

	if len(pos) < 1 {
		return fmt.Errorf("please provide a search query")
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
	n := fs.Int("n", 10, "number of results")
	offset := fs.Int("offset", 0, "pagination offset (0-based)")
	index := fs.String("index", "", "CPC subclass (e.g. H04W, auto)")
	typ := fs.String("type", "", "document type: patent|npl")
	asJSON := fs.Bool("json", false, "print raw JSON output")
	pos := parseArgs(fs, args)

	if len(pos) < 1 {
		return fmt.Errorf("please provide a patent number (e.g. US7654321B2)")
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
	q := fs.String("q", "", "text query")
	pos := parseArgs(fs, args)

	if len(pos) < 1 || *q == "" {
		return fmt.Errorf("usage: pqai snippet|mapping <pn> -q <query>")
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
		return fmt.Errorf("please provide a document ID (e.g. US7654321B2)")
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
		return fmt.Errorf("please provide a patent number (e.g. US7654321B2)")
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
		return fmt.Errorf("usage: pqai vector <pn> <field>  (field: cpcs | abstract)")
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
	name := fs.String("name", "PoC", "dataset name")
	n := fs.Int("n", 0, "sample number")
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
	thumb := fs.Bool("thumb", false, "list thumbnails instead")
	pos := parseArgs(fs, args)

	if len(pos) < 1 {
		return fmt.Errorf("please provide a patent number (e.g. US7654321B2)")
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
	thumb := fs.Bool("thumb", false, "download as a thumbnail")
	w := fs.Int("w", 0, "thumbnail width (px)")
	h := fs.Int("h", 0, "thumbnail height (px)")
	out := fs.String("o", "", "output file path (default: <pn>_<n>.png)")
	pos := parseArgs(fs, args)

	if len(pos) < 2 {
		return fmt.Errorf("usage: pqai drawing <pn> <n> [-thumb] [-w 300] [-o out.png]")
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
	fmt.Printf("Saved: %s (%d bytes)\n", path, len(body))
	return nil
}

func cmdText(c *Client, route, key string, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide text")
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
