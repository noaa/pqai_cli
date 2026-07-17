# pqai-cli

*[한국어](README.kr.md)*

A Go CLI for the [PQAI API](https://projectpq.ai/patent-search-api-by-pqai/) — prior-art search, similar-document search, patent data/drawing lookup, CPC classification suggestions, and more.

Want to see how each feature actually behaves, with real API responses? Check out [`FEATURES.en.md`](FEATURES.en.md) (in particular, `mapping` — which splits a claim into elements and matches each one against a document — is a feature that's hard to appreciate until you've tried it).

## Install

There are three ways to get `pqai`: a one-line installer script (easiest, no Go required), a manual binary download, or building it yourself from source.

### Option A: One-line installer (easiest, no Go required)

**macOS / Linux** — open a terminal and run:

```sh
curl -fsSL https://raw.githubusercontent.com/noaa/pqai_cli/main/install.sh | sh
```

This detects your OS/architecture, downloads the matching binary from the [Releases page](../../releases), and installs it to `~/.local/bin/pqai`. Open a **new terminal window** afterward so your shell picks up the updated `PATH`.

**macOS note**: since these binaries aren't notarized by an Apple Developer account, Gatekeeper may block the first run with "cannot be opened because the developer cannot be verified." If that happens, remove the quarantine flag once:

```bash
xattr -d com.apple.quarantine ~/.local/bin/pqai
```

**Windows** — open PowerShell (Start menu → search "PowerShell") and run:

```powershell
irm https://raw.githubusercontent.com/noaa/pqai_cli/main/install.ps1 | iex
```

Or, from Command Prompt (Start menu → search "cmd"):

```cmd
curl -fsSL https://raw.githubusercontent.com/noaa/pqai_cli/main/install.bat -o "%TEMP%\pqai-install.bat" && "%TEMP%\pqai-install.bat"
```

(Requires Windows 10 build 1803 or later — `curl` and `tar` are built in.) Open a **new terminal window** afterward to use the `pqai` command.

If you'd rather install from source instead of a prebuilt binary, add `--source` (requires Go 1.21+ already installed):

```sh
curl -fsSL https://raw.githubusercontent.com/noaa/pqai_cli/main/install.sh | sh -s -- --source
```

### Option B: Manual binary download

Go to the [Releases page](../../releases) of this repository and download the archive for your OS/architecture, e.g.:

- macOS (Apple Silicon): `pqai-darwin-arm64.tar.gz`
- macOS (Intel): `pqai-darwin-amd64.tar.gz`
- Linux (x86_64): `pqai-linux-amd64.tar.gz`
- Linux (arm64): `pqai-linux-arm64.tar.gz`
- Windows (x86_64): `pqai-windows-amd64.zip`

Then extract it and move the binary somewhere on your `PATH` (or just run it from the folder you extracted it to):

```bash
tar -xzf pqai-darwin-arm64.tar.gz     # or unzip on Windows
./pqai help
```

The same macOS Gatekeeper note from Option A applies here too.

### Option C: Build from source

This is a command-line tool written in the [Go](https://go.dev) programming language. If you've never used Go before, don't worry — you don't need to know how to write Go code, you just need the Go compiler installed once to build the program.

#### Step 1: Install Go (skip if you already have it)

Check whether Go is already installed:

```bash
go version
```

If you see something like `go version go1.22.0 darwin/arm64`, Go is already installed — skip to Step 2.

If you get a "command not found" error, install Go:

- **macOS**: install [Homebrew](https://brew.sh) if you don't have it, then run `brew install go`. Alternatively, download the installer from [go.dev/dl](https://go.dev/dl/).
- **Windows**: download and run the installer (`.msi`) from [go.dev/dl](https://go.dev/dl/). It sets everything up automatically — just click through the installer.
- **Linux**: use your package manager (e.g. `sudo apt install golang-go` on Ubuntu/Debian), or download the archive from [go.dev/dl](https://go.dev/dl/) and follow the [official install instructions](https://go.dev/doc/install).

After installing, open a **new** terminal window (so it picks up the updated PATH) and run `go version` again to confirm.

#### Step 2: Download this repository

If you have `git` installed, copy the repo URL from the green "Code" button on the GitHub page and run:

```bash
git clone <repo-url>
cd pqai-cli
```

Otherwise, click the green "Code" button on the GitHub repo page → "Download ZIP", then unzip it and open a terminal in that folder.

#### Step 3: Build the CLI

From inside the project folder, run:

```bash
go build -o pqai .
```

This compiles the Go source files into a single executable file named `pqai` (or `pqai.exe` on Windows) in the current folder. You only need to do this once (or again after pulling new code changes) — after that, you just run the `pqai` program directly, no Go knowledge required for day-to-day use.

Verify it worked:

```bash
./pqai help
```

(On Windows, use `.\pqai.exe help` or `pqai.exe help` in PowerShell/cmd.)

If you see the usage text printed, the build succeeded and you're ready to go — the rest of this README shows the actual commands.

## Authentication

You need an API token issued from your PQAI+ subscription account page (except for the drawing-download route).

- Set it via the `PQAI_API_KEY` environment variable, or
- Drop a `.env` file in the project root and it's read automatically (an existing shell environment variable takes precedence).

```
PQAI_API_KEY=your_token_here
# PQAI_ENDPOINT=https://api.projectpq.ai   # optional, override the API base URL
```

**Note**: The PQAI+ plan is $20/month for roughly 20 calls. Test sparingly.

## Commands

### 1. Natural-language / text search

#### `search <query>` — prior-art document search (`/search/102/`)

Feed it a natural-language sentence or paragraph as the query, and it returns related patents/papers ranked by similarity. Example: `"a drone that can extinguish fires autonomously"`.

```bash
pqai search "a fire fighting drone" -n 5
pqai search "wireless charging for electric vehicles" -after 2018-01-01 -type patent
pqai search "battery thermal management" -index H01M -snip -json
```

Flags:
| Flag | Meaning | Example |
|---|---|---|
| `-n` | number of results (default 10) | `-n 20` |
| `-offset` | pagination offset (0-based) | `-offset 10` |
| `-index` | restrict search to a CPC subclass (`auto` = auto-select) | `-index H04W` |
| `-cc` | country code filter (comma-separated) | `-cc US,EP,WO` |
| `-dtype` | cutoff-date basis (`priority`/`publication`/`filing`) | `-dtype priority` |
| `-after` | only documents after this date | `-after 2016-01-01` |
| `-before` | only documents before this date | `-before 2019-12-31` |
| `-type` | document type (`patent`/`npl`) | `-type patent` |
| `-snip` | include a snippet matching the query | `-snip` |
| `-maps` | include per-element query-to-document mapping | `-maps` |
| `-lq` | latent-query JSON (steer results using relevant/irrelevant patents) | `-lq '{"relevant":["US123"],"irrelevant":[]}'` |
| `-json` | print raw JSON instead of a human-readable summary | `-json` |

By default the output shows rank, patent number, similarity score, publication date, title, owner, and a snippet/abstract summary in a human-readable format. Add `-json` to see the raw response as-is.

#### `combos <query>` — prior-art "combination" search (`/search/103/`)

Instead of a single document, this finds cases where a combination of multiple documents covers a claim (useful for §103 non-obviousness review). Flags are the same as `search`.

```bash
pqai combos "battery management system with thermal runaway detection" -n 10
```

### 2. Search anchored on a specific patent

#### `prior-art <pn>` — prior art filed before the given patent's filing date (`/prior-art/patent/`)

```bash
pqai prior-art US7654321B2 -n 10
```

#### `similar <pn>` — documents similar to the given patent (`/similar/`)

```bash
pqai similar US10112730B2 -n 10 -type patent
```

Flags shared by both commands: `-n`, `-offset`, `-index`, `-type`, `-json`.

### 3. Query-document pair analysis

#### `snippet <pn> -q <text>` — retrieve the snippet matching a query (`/snippets/`)

```bash
pqai snippet US10112730B2 -q "autonomous drone fire suppression"
```

#### `mapping <pn> -q <text>` — per-element query-to-document mapping (`/mappings/`)

Returns a mapping showing which part of the document corresponds to each element of the claim (useful for invalidity/infringement analysis). Feed the query as a claim split by semicolons/newlines, and each element is automatically separated and mapped individually (real example: [`FEATURES.en.md`](FEATURES.en.md#4-split-a-claim-into-elements-and-match-each-one-against-a-document--mapping--the-most-underrated-feature)).

```bash
pqai mapping US10112730B2 -q "a rotor assembly configured to..."
```

### 4. Data lookup

| Command | Route | Description |
|---|---|---|
| `patent <pn>` | `/patents/:pn` | look up bibliographic/text data for a patent |
| `document <id>` | `/documents/` | look up a document (patent/paper) in the PQAI database |
| `vector <pn> <field>` | `/patents/:pn/vectors/:field` | get a patent's embedding vector (`field`: `cpcs` or `abstract`) |
| `dataset -name <n> -n <i>` | `/datasets/` | fetch a dataset sample |

```bash
pqai patent US7654321B2
pqai document US7654321B2
pqai vector US7654321B2 abstract
pqai dataset -name PoC -n 23
```

### 5. Drawings

| Command | Token required? | Description |
|---|---|---|
| `drawings <pn> [-thumb]` | yes | list a patent's drawings (or thumbnails) |
| `drawing <pn> <n> [-thumb] [-w px] [-h px] [-o path]` | no | download a specific drawing (PNG/JPEG) to a file |

```bash
pqai drawings US7654321B2
pqai drawing US7654321B2 4 -o drawing4.png
pqai drawing US7654321B2 4 -w 300 -o thumb4.png   # thumbnail, 300px wide
```

If you specify both `-w` and `-h`, the image may be stretched if the aspect ratio doesn't match, so it's recommended to specify only one (per the API docs).

### 6. Classification suggestions

| Command | Route | Description |
|---|---|---|
| `cpcs <text>` | `/suggest/cpcs` | suggest CPC classifications for a text |
| `gaus <text>` | `/predict/gaus` | suggest a USPTO Group Art Unit for a text |
| `cpc-def <cpc>` | `/definitions/cpcs` | look up the definition/description of a CPC class |

```bash
pqai cpcs "fire fighting drones"
pqai gaus "fire fighting drones"
pqai cpc-def H04W52/02
```

## Notes

- All positional arguments (queries, patent numbers, etc.) and flags can be mixed in any order. E.g. `pqai drawing US123 2 -o out.png` and `pqai drawing -o out.png US123 2` behave identically.
- Most routes return JSON, pretty-printed by default. `search`/`combos`/`prior-art`/`similar` show a human-readable summary by default, and `-json` shows the raw response.
- Only the individual drawing/thumbnail image route (`drawing`) works without a token; everything else requires `PQAI_API_KEY`.
