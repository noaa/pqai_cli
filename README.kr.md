# pqai-cli

*[English](README.md)*

[PQAI API](https://projectpq.ai/patent-search-api-by-pqai/)를 커맨드라인에서 사용하기 위한 Go CLI. 특허 선행기술 검색, 유사 문서 검색, 특허 데이터/도면 조회, CPC 분류 제안 등을 지원합니다.

각 기능이 실제로 어떻게 동작하는지 실제 API 응답 예시와 함께 보고 싶다면 [`FEATURES.kr.md`](FEATURES.kr.md)를 참고하세요 (특히 청구항을 요소별로 쪼개 대조해주는 `mapping`은 써보기 전엔 알기 어려운 기능입니다).

## 설치

`pqai`를 설치하는 방법은 세 가지입니다: 원라인 설치 스크립트(가장 쉬움, Go 불필요), 바이너리 수동 다운로드, 또는 소스에서 직접 빌드하는 방법입니다.

### 방법 A: 원라인 설치 스크립트 (가장 쉬움, Go 불필요)

**macOS / Linux** — 터미널을 열고 다음을 실행합니다:

```sh
curl -fsSL https://raw.githubusercontent.com/noaa/pqai_cli/main/install.sh | sh
```

이 스크립트는 자동으로 OS/아키텍처를 감지해서 [Releases 페이지](../../releases)에서 맞는 바이너리를 받아 `~/.local/bin/pqai`에 설치합니다. 설치 후에는 **새 터미널 창**을 열어야 셸이 갱신된 `PATH`를 인식합니다.

**macOS 참고**: 이 바이너리들은 Apple Developer 계정으로 노터라이제이션(공증)되어 있지 않으므로, 처음 실행할 때 Gatekeeper가 "개발자를 확인할 수 없기 때문에 열 수 없습니다" 경고를 띄울 수 있습니다. 이 경우 아래 명령으로 격리(quarantine) 속성을 한 번만 제거하면 됩니다:

```bash
xattr -d com.apple.quarantine ~/.local/bin/pqai
```

**Windows** — PowerShell을 열고(시작 메뉴 → "PowerShell" 검색) 다음을 실행합니다:

```powershell
irm https://raw.githubusercontent.com/noaa/pqai_cli/main/install.ps1 | iex
```

또는 명령 프롬프트(시작 메뉴 → "cmd" 검색)에서:

```cmd
curl -fsSL https://raw.githubusercontent.com/noaa/pqai_cli/main/install.bat -o "%TEMP%\pqai-install.bat" && "%TEMP%\pqai-install.bat"
```

(Windows 10 build 1803 이상 필요 — `curl`과 `tar`가 기본 내장되어 있습니다.) 설치 후에는 **새 터미널 창**을 열어야 `pqai` 명령을 사용할 수 있습니다.

미리 빌드된 바이너리 대신 소스에서 설치하고 싶다면(Go 1.21+가 이미 설치되어 있어야 함) `--source` 옵션을 추가하세요:

```sh
curl -fsSL https://raw.githubusercontent.com/noaa/pqai_cli/main/install.sh | sh -s -- --source
```

### 방법 B: 바이너리 수동 다운로드

이 리포지토리의 [Releases 페이지](../../releases)에서 자신의 OS/아키텍처에 맞는 압축 파일을 다운로드하세요:

- macOS (Apple Silicon): `pqai-darwin-arm64.tar.gz`
- macOS (Intel): `pqai-darwin-amd64.tar.gz`
- Linux (x86_64): `pqai-linux-amd64.tar.gz`
- Linux (arm64): `pqai-linux-arm64.tar.gz`
- Windows (x86_64): `pqai-windows-amd64.zip`

압축을 풀고 `PATH`에 있는 위치로 옮기거나, 그냥 압축 푼 폴더에서 바로 실행하면 됩니다:

```bash
tar -xzf pqai-darwin-arm64.tar.gz     # 윈도우는 압축 프로그램으로 압축 해제
./pqai help
```

방법 A와 동일한 macOS Gatekeeper 참고 사항이 여기에도 적용됩니다.

### 방법 C: 소스에서 직접 빌드

이 도구는 [Go](https://go.dev) 언어로 작성된 커맨드라인 프로그램입니다. Go를 한 번도 써본 적이 없어도 괜찮습니다 — Go로 코드를 직접 작성할 필요는 없고, 프로그램을 빌드(컴파일)하기 위해 딱 한 번 Go 컴파일러만 설치하면 됩니다.

#### 1단계: Go 설치 (이미 설치되어 있다면 건너뛰세요)

먼저 Go가 이미 설치되어 있는지 확인합니다:

```bash
go version
```

`go version go1.22.0 darwin/arm64`처럼 버전이 출력되면 이미 설치된 것이니 2단계로 넘어가세요.

`command not found` 같은 오류가 나오면 아래처럼 설치합니다:

- **macOS**: [Homebrew](https://brew.sh)가 없다면 먼저 설치한 뒤 `brew install go` 실행. 또는 [go.dev/dl](https://go.dev/dl/)에서 설치 파일을 직접 다운로드해도 됩니다.
- **Windows**: [go.dev/dl](https://go.dev/dl/)에서 `.msi` 설치 파일을 받아 실행하면 됩니다. 설치 마법사가 알아서 환경변수까지 설정해주므로 계속 다음(Next)만 눌러도 됩니다.
- **Linux**: 배포판의 패키지 매니저를 사용하거나(예: Ubuntu/Debian에서 `sudo apt install golang-go`), [go.dev/dl](https://go.dev/dl/)에서 압축 파일을 받아 [공식 설치 안내](https://go.dev/doc/install)를 따르세요.

설치가 끝나면 **새 터미널 창**을 열어(PATH가 갱신되도록) `go version`을 다시 실행해 정상적으로 버전이 나오는지 확인하세요.

#### 2단계: 이 리포지토리 다운로드

`git`이 설치되어 있다면, GitHub 페이지의 초록색 "Code" 버튼에서 리포 주소를 복사해 다음을 실행합니다:

```bash
git clone <repo-url>
cd pqai-cli
```

git이 없다면, GitHub 리포 페이지에서 초록색 "Code" 버튼 → "Download ZIP"을 눌러 받은 뒤 압축을 풀고, 그 폴더에서 터미널을 엽니다.

#### 3단계: CLI 빌드

프로젝트 폴더 안에서 다음을 실행합니다:

```bash
go build -o pqai .
```

이 명령은 Go 소스 코드를 컴파일해서 현재 폴더에 `pqai`(윈도우는 `pqai.exe`)라는 실행 파일 하나를 만들어줍니다. 이 과정은 최초 1회만(또는 코드가 업데이트됐을 때만) 하면 되고, 그 이후에는 매번 `pqai` 프로그램만 실행하면 되므로 평소 사용에는 Go 지식이 전혀 필요하지 않습니다.

빌드가 잘 됐는지 확인:

```bash
./pqai help
```

(윈도우 PowerShell/cmd에서는 `.\pqai.exe help` 또는 `pqai.exe help`를 사용하세요.)

사용법 텍스트가 출력되면 빌드 성공입니다 — 이 README의 나머지 부분에서 실제 명령어들을 확인하세요.

## 인증

PQAI+ 구독 계정 페이지에서 발급받은 API 토큰이 필요합니다 (도면 다운로드 라우트 제외).

**권장: 전역 설정에 한 번만 저장해두면 어느 폴더에서든 동작합니다:**

```bash
pqai config set-api-key your_token_here

# 또는 기존 .env 파일에서 가져오기:
pqai config set-api-key --from-dotenv .env
```

이 명령은 토큰을 사용자 설정 파일(macOS/Linux: `~/.config/pqai/config.env`, Windows: `%AppData%\pqai\config.env`)에 `0600` 권한으로 저장하므로, 현재 `pqai`를 실행 중인 폴더가 어디인지와 무관하게 동작합니다. 현재 설정 상태는 다음으로 확인할 수 있습니다:

```bash
pqai config show
```

그 외에 토큰을 설정하는 방법(일회성 오버라이드나 CI 환경에 유용):

- 셸에서 `PQAI_API_KEY` 환경변수를 export
- 현재 프로젝트 폴더에 `.env` 파일을 두기 — 이 방법은 해당 폴더에서 `pqai`를 실행할 때만 적용됩니다.

```
PQAI_API_KEY=your_token_here
# PQAI_ENDPOINT=https://api.projectpq.ai   # 선택, API 주소 재정의
```

여러 방법이 동시에 설정되어 있다면 다음 순서로 우선 적용됩니다 (위가 우선):

1. 셸에서 export한 `PQAI_API_KEY`
2. 현재 폴더의 `.env` 파일
3. `pqai config set-api-key`로 저장한 전역 설정 파일

**주의**: PQAI+ 요금제는 월 $20에 약 20회 호출 한도입니다. 호출은 아껴서 테스트하세요.

## 명령어

### 1. 자연어/텍스트 검색

#### `search <query>` — 선행기술 문서 검색 (`/search/102/`)

자연어 문장이나 문단을 쿼리로 넣으면 관련 특허/논문을 유사도 순으로 반환합니다. 예: `"a drone that can extinguish fires autonomously"`.

```bash
pqai search "a fire fighting drone" -n 5
pqai search "wireless charging for electric vehicles" -after 2018-01-01 -type patent
pqai search "battery thermal management" -index H01M -snip -json
```

플래그:
| 플래그 | 의미 | 예시 |
|---|---|---|
| `-n` | 결과 개수 (기본 10) | `-n 20` |
| `-offset` | 페이지네이션 오프셋 (0부터) | `-offset 10` |
| `-index` | CPC 서브클래스로 검색 범위 제한 (`auto`=자동선택) | `-index H04W` |
| `-cc` | 국가 코드 필터 (콤마 구분) | `-cc US,EP,WO` |
| `-dtype` | 컷오프 날짜 기준 (`priority`/`publication`/`filing`) | `-dtype priority` |
| `-after` | 이 날짜 이후 문서만 | `-after 2016-01-01` |
| `-before` | 이 날짜 이전 문서만 | `-before 2019-12-31` |
| `-type` | 문서 유형 (`patent`/`npl`) | `-type patent` |
| `-snip` | 검색어와 매칭되는 스니펫 포함 | `-snip` |
| `-maps` | 쿼리-문서 요소별 매핑 포함 | `-maps` |
| `-lq` | 잠재 쿼리(관련/비관련 특허로 검색 결과 보정) JSON | `-lq '{"relevant":["US123"],"irrelevant":[]}'` |
| `-json` | 사람이 읽기 좋은 요약 대신 원본 JSON 출력 | `-json` |

기본 출력은 순위, 특허번호, 유사도 점수, 공개일, 제목, 소유자, 스니펫/초록 요약을 사람이 읽기 좋은 형태로 보여줍니다. `-json`을 붙이면 원본 응답을 그대로 확인할 수 있습니다.

#### `combos <query>` — 선행기술 "조합" 검색 (`/search/103/`)

단일 문서가 아니라 여러 문서의 조합으로 청구항을 커버하는 케이스를 찾습니다 (예: 103조 자명성 검토용). 플래그는 `search`와 동일합니다.

```bash
pqai combos "battery management system with thermal runaway detection" -n 10
```

### 2. 특정 특허 기준 검색

#### `prior-art <pn>` — 해당 특허의 출원일 이전 선행기술 검색 (`/prior-art/patent/`)

```bash
pqai prior-art US7654321B2 -n 10
```

#### `similar <pn>` — 해당 특허와 유사한 문서 검색 (`/similar/`)

```bash
pqai similar US10112730B2 -n 10 -type patent
```

두 명령 공통 플래그: `-n`, `-offset`, `-index`, `-type`, `-json`.

**참고**: `search`/`combos`와 달리, `/prior-art/patent/`와 `/similar/` API 라우트는 날짜(`-dtype`/`-after`/`-before`)나 국가(`-cc`) 필터를 아예 지원하지 않습니다 — [PQAI 공식 API 문서](https://api.projectpq.ai/docs)에 따르면 이 두 라우트는 `pn`, `n`, `offset`, `index`, `type`만 받습니다. `prior-art`는 이미 "해당 특허의 출원일 이전"으로 날짜가 암묵적으로 고정되어 있고, 두 라우트 모두 국가로 좁히는 기능 자체가 없습니다. 날짜/국가로 필터링하려면 `search`/`combos`를 텍스트 쿼리로 사용해야 합니다.

### 3. 쿼리-문서 쌍 분석

#### `snippet <pn> -q <text>` — 쿼리와 매칭되는 스니펫 조회 (`/snippets/`)

```bash
pqai snippet US10112730B2 -q "autonomous drone fire suppression"
```

#### `mapping <pn> -q <text>` — 쿼리-문서 요소별 매핑 조회 (`/mappings/`)

청구항 요소별로 문서의 어느 부분이 대응되는지 매핑을 반환합니다 (특허 무효성/침해 분석에 유용). 쿼리를 세미콜론/개행으로 구분된 청구항 형태로 넣으면 요소별로 자동 분리되어 각각 매핑됩니다 (실제 예시: [`FEATURES.kr.md`](FEATURES.kr.md#4-청구항을-요소별로-쪼개서-문서와-대조한다--mapping--가장-저평가된-기능)).

```bash
pqai mapping US10112730B2 -q "a rotor assembly configured to..."
```

### 4. 데이터 조회

| 명령 | 라우트 | 설명 |
|---|---|---|
| `patent <pn>` | `/patents/:pn` | 특허 서지/텍스트 데이터 조회 |
| `document <id>` | `/documents/` | PQAI DB에서 문서(특허/논문) 조회 |
| `vector <pn> <field>` | `/patents/:pn/vectors/:field` | 특허의 임베딩 벡터 조회 (`field`: `cpcs` 또는 `abstract`) |
| `dataset -name <n> -n <i>` | `/datasets/` | 데이터셋 샘플 조회 |

```bash
pqai patent US7654321B2
pqai document US7654321B2
pqai vector US7654321B2 abstract
pqai dataset -name PoC -n 23
```

### 5. 도면

| 명령 | 토큰 필요? | 설명 |
|---|---|---|
| `drawings <pn> [-thumb]` | 필요 | 특허 도면(또는 썸네일) 목록 조회 |
| `drawing <pn> <n> [-thumb] [-w px] [-h px] [-o path]` | 불필요 | 특정 도면(PNG/JPEG)을 파일로 다운로드 |

```bash
pqai drawings US7654321B2
pqai drawing US7654321B2 4 -o drawing4.png
pqai drawing US7654321B2 4 -w 300 -o thumb4.png   # 썸네일, 너비 300px
```

`-w`와 `-h`를 동시에 지정하면 비율이 안 맞을 경우 이미지가 늘어날 수 있으니 하나만 지정하는 걸 권장합니다 (API 문서 권장사항).

### 6. 분류 제안

| 명령 | 라우트 | 설명 |
|---|---|---|
| `cpcs <text>` | `/suggest/cpcs` | 텍스트에 대한 CPC 분류 제안 |
| `gaus <text>` | `/predict/gaus` | 텍스트에 대한 USPTO Group Art Unit 제안 |
| `cpc-def <cpc>` | `/definitions/cpcs` | CPC 클래스 정의/설명 조회 |

```bash
pqai cpcs "fire fighting drones"
pqai gaus "fire fighting drones"
pqai cpc-def H04W52/02
```

## 참고 사항

- 모든 위치 인자(쿼리, 특허번호 등)와 플래그는 순서 상관없이 섞어 쓸 수 있습니다. 예: `pqai drawing US123 2 -o out.png`, `pqai drawing -o out.png US123 2` 모두 동일하게 동작합니다.
- 대부분의 라우트는 JSON을 반환하며 기본적으로 pretty-print 됩니다. `search`/`combos`/`prior-art`/`similar`는 기본적으로 사람이 읽기 좋은 요약을 보여주고, `-json`으로 원본 응답을 볼 수 있습니다.
- 도면/썸네일 개별 이미지 라우트(`drawing`)만 토큰 없이 동작하며, 나머지는 모두 `PQAI_API_KEY`가 필요합니다.
