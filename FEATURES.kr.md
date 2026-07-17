# PQAI API로 할 수 있는 것들 (실제 응답 기반)

*[English](FEATURES.en.md)*

[PQAI](https://projectpq.ai)는 USPTO와 여러 대학이 참여한 오픈소스 프로젝트에서 시작된 **AI 기반 특허 선행기술 검색 API**입니다. `pqai-cli`는 이 API를 터미널에서 바로 쓸 수 있게 감싼 도구입니다.

이 문서는 "일단 써봐야 아는" 기능들 위주로, **실제 API 호출 결과**를 그대로 붙여 넣어 설명합니다. (호출에 사용한 값 일부는 예시용 가상의 청구항이며, 아래 응답은 모두 실제 `https://api.projectpq.ai`가 반환한 원본입니다.)

---

## 1. 문장만 넣으면 관련 특허를 찾아준다 — `search`

특허 검색이라고 하면 보통 키워드(AND/OR)를 조합해야 하지만, PQAI는 **자연어 문단**을 그대로 쿼리로 받습니다. 내부적으로 문장을 임베딩으로 변환해 의미 기반 유사도 검색을 합니다.

```bash
pqai search "A method for generating a digital image from a user's rough hand-drawn sketch made with a finger or stylus input in a note-taking application, wherein the user encircles the sketch or a blank area with a selection gesture, and the system generates one or more candidate images based on the sketch and an associated text description or surrounding note content, allowing the user to select among multiple rendering styles such as sketch, illustration, and animation." \
  -dtype filing -after 2023-01-01 -cc US -type patent -n 10
```

이 문장은 사실 **Apple의 Image Wand**(Notes 앱에서 손그림을 AI 이미지로 바꿔주는 기능)를 설명한 것입니다. 실제 결과 1위:

```
1. US2025200831A1 — score 0.6821 — 공개일 2025-06-19
   METHOD FOR AUTOMATICALLY GENERATING SKETCH IMAGE, ...
   소유자: Korea Advanced Institute of Science and Technology
```

기능 설명만으로 관련 분야(스케치 → 이미지 생성) 특허들을 찾아냈고, 그중에는 Apple 소유 특허(3D 객체 배치 관련)도 포함되어 있었습니다. 전체 결과와 분석은 [`docs/image-wand-search.kr.md`](docs/image-wand-search.kr.md) 참고.

**날짜/국가/문서유형 필터**를 조합할 수 있어서, "우리 회사가 출원 예정인 기능이 최근 몇 년 내 이미 등록/공개됐는지" 같은 사전 조사에 바로 쓸 수 있습니다.

---

## 2. 단일 문서가 아니라 "조합"으로 찾는다 — `combos`

`search`가 "이 문장과 가장 비슷한 문서 하나"를 찾는다면, `combos`(`/search/103/`)는 **여러 문서를 조합했을 때** 청구항 전체를 커버하는 조합을 찾습니다. 특허법 103조(자명성, non-obviousness) 검토에 대응하는 기능으로, "선행문헌 A + B를 합치면 이 청구항이 자명해지는가"를 사람이 직접 대조하지 않고 후보군부터 좁힐 수 있습니다.

```bash
pqai combos "battery management system with thermal runaway detection" -n 10
```

---

## 3. 특정 특허를 기준으로 검색한다 — `prior-art`, `similar`

키워드 없이 **특허번호만 넣어도** 검색이 됩니다. `prior-art`는 해당 특허의 출원일 "이전" 문서만 걸러서 보여줘 무효성 조사(invalidity search)에 바로 쓸 수 있고, `similar`는 날짜 제한 없이 유사 문서를 찾습니다.

```bash
pqai prior-art US11868178B2 -n 10
```

`US11868178B2`(반지형 웨어러블 디바이스 특허)로 실제 호출한 결과, 500건의 후보 중 1위는:

```
1. CN218548276U — score 0.7606 — 공개일 2023-02-28
   Wearable device (도전성 버튼 어셈블리를 가진 웨어러블 기기)
2. KR20170091346A — score 0.7472 — Ring type wearable device
3. WO2021017915A1 — score 0.7410 — Finger-ring monitoring device and monitoring system
```

특허 하나만 알면 그와 겹치는 선행기술 풀을 바로 얻을 수 있다는 점이 핵심입니다 — 검색어를 직접 고민할 필요가 없습니다.

**주의할 점**: `search`/`combos`와 달리 이 두 라우트는 날짜/국가 필터를 아예 지원하지 않습니다. [PQAI 공식 API 문서](https://api.projectpq.ai/docs)에 따르면 `/prior-art/patent/`와 `/similar/`는 `pn`, `n`, `offset`, `index`, `type`만 받고, `-cc`/`-dtype`/`-after`/`-before`는 존재하지 않습니다 (`prior-art`는 이미 "해당 특허의 출원일"로 날짜가 암묵적으로 고정되어 있어 별도 날짜 필터가 필요 없고, 국가 필터는 애초에 라우트 자체에 없습니다). 이 플래그들을 시도하면 요청이 서버로 가기도 전에 클라이언트 단에서 바로 에러(`flag provided but not defined`)가 납니다. 날짜/국가로 필터링해야 한다면 `search`/`combos`를 텍스트 쿼리로 사용하세요.

---

## 4. 청구항을 요소별로 쪼개서 문서와 대조한다 — `mapping` ⭐ (가장 저평가된 기능)

`snippet`/`mapping`은 "이 문서가 내 쿼리와 관련 있다"는 결과 이상으로, **청구항의 각 구성요소(limitation)가 문서의 어느 서술에 대응하는지**를 요소별로 잘라 보여줍니다. 이게 바로 변리사/심사관이 손으로 만드는 "claim chart"를 자동화한 것과 같은 발상입니다.

핵심은 쿼리를 **줄바꿈(개행)으로 구분된 청구항 형식**으로 넣으면, API가 자동으로 요소를 분리해 각각을 매핑해준다는 점 — 이건 문서만 봐서는 알기 어렵고 실제로 멀티라인 쿼리를 넣어봐야 드러납니다.

```bash
Q='A drone comprising:
a rotor assembly configured to rotate about an axis and provide lift;
a camera coupled to the drone for capturing images of a fire;
a fire suppression module configured to release a fire suppressant material;
a controller configured to autonomously navigate the drone toward the fire.'

pqai mapping US10112730B2 -q "$Q"
```

실제 응답(요약):

```json
{
  "mapping": [
    {
      "element": "a rotor assembly configured to rotate about an axis and provide lift;",
      "mapping": "...fixed wings, rotor wings, helicopters... dynamic equation of motion mode..."
    },
    {
      "element": "a camera coupled to the drone for capturing images of a fire;",
      "mapping": "...camera based speed-recognition system... database of vibration signatures..."
    },
    {
      "element": "a fire suppression module configured to release a fire suppressant material;",
      "mapping": "...numerically estimate EOM parameters... propeller characteristics using neural networks..."
    },
    {
      "element": "a controller configured to autonomously navigate the drone toward the fire.",
      "mapping": "...design payload of standard commercial drones..."
    }
  ]
}
```

청구항 4개 요소를 각각 별도 항목으로 잘라, 대상 특허(`US10112730B2`, 드론 관련 특허) 안에서 가장 관련성 높은 서술을 하나씩 짝지어 돌려줍니다. 단일 요소만 넣으면(`-q "a rotor assembly..."`) 그 요소 하나만 매핑되므로, **청구항 전체를 한 번에 대조하고 싶다면 세미콜론/개행으로 나눠 넣는 것이 핵심 활용법**입니다.

`snippet`은 이보다 단순한 버전으로, 쿼리 전체와 가장 관련 있는 스니펫 하나만 반환합니다:

```bash
pqai snippet US10112730B2 -q "a rotor assembly configured to rotate about an axis and provide lift"
```
```json
{
  "snippet": "... (UAVs), fixed wing aircrafts, rotary wing aircrafts and helicopters. The one or more rotating parts of said vehicle may include rotor blades, propeller blades, turbine blades, jet/gas compressors, reciprocating engine, or similar parts. Receiving the instantaneous ..."
}
```

---

## 5. 텍스트만으로 특허 분류(CPC)와 심사부서(GAU)를 예측한다 — `cpcs`, `gaus`

출원 전 초안 단계에서 "이 발명이 어느 CPC 서브클래스, 어느 USPTO 심사부(Art Unit)로 배정될지"를 미리 알 수 있으면 유사 선행기술 검색 범위를 좁히거나 심사 난이도를 가늠하는 데 도움이 됩니다.

```bash
pqai cpcs "an autonomous drone equipped with a fire suppressant tank and camera-based fire detection for aerial firefighting"
```

실제 응답 1위 (신뢰도 순 정렬, 총 61개 후보 중 상위 항목):

```json
{
  "cpc": "A62C3/0242",
  "definition": [
    ["A", "HUMAN NECESSITIES"],
    ["A62", "LIFE-SAVING; FIRE-FIGHTING"],
    ["A62C", "FIRE-FIGHTING"],
    ["A62C3/00", "Fire prevention, containment or extinguishing specially adapted for particular objects or places"],
    ["A62C3/02", "for area conflagrations, e.g. forest fires, subterranean fires"],
    ["A62C3/0228", "with delivery of fire extinguishing material by air or aircraft"],
    ["A62C3/0242", "by spraying extinguishants from the aircraft"]
  ],
  "confidence": 0.07
}
```

CPC 코드 하나만 던지는 게 아니라 **트리 전체 계층("에어리어 전소 화재" → "항공기로 소화제 전달" → "항공기에서 살포")**을 함께 반환해서, 코드를 몰라도 분류 체계를 바로 이해할 수 있습니다. 같은 텍스트로 `gaus`를 호출하면:

```bash
pqai gaus "an autonomous drone equipped with a fire suppressant tank and camera-based fire detection for aerial firefighting"
# → ["3752", "2482", "3664"]
```

USPTO Group Art Unit 후보 3개(우선순위 순)를 즉시 받습니다. 개별 CPC 코드의 정의만 다시 찾고 싶다면:

```bash
pqai cpc-def A62C3/0242
```

---

## 6. 도면은 토큰 없이도 바로 다운로드된다 — `drawings`, `drawing`

API 전체가 유료 토큰이 필요한 건 아닙니다. **개별 도면 이미지 다운로드 라우트만 인증 없이 열려 있습니다** — README에 적혀 있지만 실제로 토큰을 빼고 호출해보기 전에는 확신하기 어려운 부분이라 직접 검증했습니다.

```bash
pqai drawings US10112730B2          # 도면 목록: ["1","2","3","4","5","6"]
pqai drawing US10112730B2 1 -o fig1.png    # 토큰 없이 PNG 다운로드 성공 (321KB)
pqai drawing US10112730B2 1 -w 300 -o thumb1.png  # 너비 300px 썸네일
```

`.env`의 `PQAI_API_KEY`를 지운 상태에서도 `drawing` 명령만은 정상 동작합니다. 도면을 대량으로 스크래핑하거나 PPT/보고서에 특허 도면을 넣어야 할 때 API 크레딧을 소모하지 않고 쓸 수 있습니다.

---

## 7. 데이터 그 자체를 조회한다 — `patent`, `document`, `vector`, `dataset`

- `patent <pn>` — 서지사항, 청구항, 명세서 텍스트 등 특허 원문 데이터
- `document <id>` — PQAI 자체 색인 DB에서 문서(특허/논문 포함) 조회
- `vector <pn> <field>` — 특허의 `cpcs` 또는 `abstract` 임베딩 벡터를 직접 꺼낼 수 있음 (자체 유사도 계산, 클러스터링 등에 재사용 가능)
- `dataset -name <n> -n <i>` — PQAI가 제공하는 벤치마크/PoC 데이터셋 샘플 조회

이 명령들은 `search` 계열이 반환하는 요약 정보보다 한 단계 더 깊이 들어가, 검색 결과 이후의 후처리(자체 재랭킹, 벡터 유사도 계산 등)를 하고 싶을 때 필요합니다.

---

## 참고: 요금제와 호출 예산

PQAI+ 무료/저가 요금제는 월 호출 횟수 제한이 있습니다(예: 월 $20에 약 20회). 위 예시들은 실제로 `pqai-cli`를 통해 라이브 API에 호출해 얻은 응답이며, 대량 호출이 필요한 실험(예: 청구항 전체 무효성 조사)은 유료 크레딧을 고려해 계획하고 진행하는 것을 권장합니다.

전체 명령어 레퍼런스는 [`README.kr.md`](README.kr.md)를 참고하세요.
