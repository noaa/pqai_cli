# PQAI Search Request: Prior Art Related to Apple's Image Wand

*[한국어](image-wand-search-request.kr.md)*

## CLI invocation

```bash
pqai search "A method for generating a digital image from a user's rough hand-drawn sketch made with a finger or stylus input in a note-taking application, wherein the user encircles the sketch or a blank area with a selection gesture, and the system generates one or more candidate images based on the sketch and an associated text description or surrounding note content, allowing the user to select among multiple rendering styles such as sketch, illustration, and animation." \
  -dtype filing \
  -after 2023-01-01 \
  -cc US \
  -type patent \
  -n 10 \
  -json
```

(Only the presence of `-json` differs; all other parameters are the same for both the human-readable and JSON versions.)

## Actual API route and parameters

- Route: `GET /search/102/` (PQAI: prior-art document search by text query)
- Endpoint: `https://api.projectpq.ai/search/102/`

| Parameter | Value | Meaning |
|---|---|---|
| `q` | (full query text below) | text query |
| `dtype` | `filing` | cutoff-date basis = filing date |
| `after` | `2023-01-01` | only documents on/after this date (filing date) |
| `cc` | `US` | country code filter |
| `type` | `patent` | document type = patent (excludes papers) |
| `n` | `10` | number of results to return |
| `token` | (attached automatically from the `PQAI_API_KEY` env var) | auth token |

## Query (the `q` parameter, English functional description)

A functional description of Apple's **Image Wand** (the feature in the Notes app / Apple Intelligence / Markup menu that generates an AI image from a hand-drawn sketch or note text):

> A method for generating a digital image from a user's rough hand-drawn sketch made with a finger or stylus input in a note-taking application, wherein the user encircles the sketch or a blank area with a selection gesture, and the system generates one or more candidate images based on the sketch and an associated text description or surrounding note content, allowing the user to select among multiple rendering styles such as sketch, illustration, and animation.

## Related files

- [`image-wand-search.en.md`](./image-wand-search.en.md) — human-readable search results and analysis
- [`image-wand-search.json`](./image-wand-search.json) — raw JSON response
