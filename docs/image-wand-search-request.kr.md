# PQAI 검색 요청: Apple Image Wand 관련 선행기술

*[English](image-wand-search-request.en.md)*

## CLI 호출 명령

```bash
pqai search "A method for generating a digital image from a user's rough hand-drawn sketch made with a finger or stylus input in a note-taking application, wherein the user encircles the sketch or a blank area with a selection gesture, and the system generates one or more candidate images based on the sketch and an associated text description or surrounding note content, allowing the user to select among multiple rendering styles such as sketch, illustration, and animation." \
  -dtype filing \
  -after 2023-01-01 \
  -cc US \
  -type patent \
  -n 10 \
  -json
```

(`-json` 유무만 다르고 나머지 파라미터는 사람이 읽기 쉬운 버전과 JSON 버전 모두 동일)

## 실제 API 라우트 및 파라미터

- 라우트: `GET /search/102/` (PQAI: 텍스트 쿼리로 선행기술 문서 검색)
- 엔드포인트: `https://api.projectpq.ai/search/102/`

| 파라미터 | 값 | 의미 |
|---|---|---|
| `q` | (아래 검색식 전문) | 텍스트 쿼리 |
| `dtype` | `filing` | 컷오프 날짜 기준 = 출원일 |
| `after` | `2023-01-01` | 이 날짜(출원일) 이후 문서만 |
| `cc` | `US` | 국가 코드 필터 |
| `type` | `patent` | 문서 유형 = 특허 (논문 제외) |
| `n` | `10` | 반환 결과 개수 |
| `token` | (환경변수 `PQAI_API_KEY`에서 자동 첨부) | 인증 토큰 |

## 검색식 (q 파라미터, 영어 기능 설명)

Apple **Image Wand**(Notes 앱, Apple Intelligence, Markup 메뉴에서 손으로 그린 스케치나 메모 텍스트를 바탕으로 AI 이미지를 생성하는 기능)의 동작을 서술한 기능 설명:

> A method for generating a digital image from a user's rough hand-drawn sketch made with a finger or stylus input in a note-taking application, wherein the user encircles the sketch or a blank area with a selection gesture, and the system generates one or more candidate images based on the sketch and an associated text description or surrounding note content, allowing the user to select among multiple rendering styles such as sketch, illustration, and animation.

## 관련 파일

- [`image-wand-search.kr.md`](./image-wand-search.kr.md) — 사람이 읽기 쉬운 검색 결과 및 분석
- [`image-wand-search.json`](./image-wand-search.json) — 원본 JSON 응답
