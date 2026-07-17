# What you can do with the PQAI API (backed by real responses)

*[한국어](FEATURES.kr.md)*

[PQAI](https://projectpq.ai) is an **AI-powered prior-art search API** that grew out of an open-source project involving the USPTO and several universities. `pqai-cli` is a wrapper that lets you use this API straight from your terminal.

This document focuses on the features you can't fully appreciate without trying them — and every response below is pasted in **as-is from real API calls**. (Some of the input values are illustrative/hypothetical claims, but every response shown was actually returned by `https://api.projectpq.ai`.)

---

## 1. Feed it a sentence, get related patents — `search`

Patent search usually means combining keywords (AND/OR), but PQAI accepts a **natural-language paragraph** as the query directly. Under the hood it turns the sentence into an embedding and runs semantic similarity search.

```bash
pqai search "A method for generating a digital image from a user's rough hand-drawn sketch made with a finger or stylus input in a note-taking application, wherein the user encircles the sketch or a blank area with a selection gesture, and the system generates one or more candidate images based on the sketch and an associated text description or surrounding note content, allowing the user to select among multiple rendering styles such as sketch, illustration, and animation." \
  -dtype filing -after 2023-01-01 -cc US -type patent -n 10
```

That paragraph is actually a description of **Apple's Image Wand** (the Notes-app feature that turns a hand-drawn sketch into an AI-generated image). The actual #1 result:

```
1. US2025200831A1 — score 0.6821 — published 2025-06-19
   METHOD FOR AUTOMATICALLY GENERATING SKETCH IMAGE, ...
   Owner: Korea Advanced Institute of Science and Technology
```

Just from a functional description, it surfaced patents in the same neighborhood (sketch → image generation), including one owned by Apple itself (related to 3D object placement). Full results and analysis in [`docs/image-wand-search.en.md`](docs/image-wand-search.en.md).

You can combine **date / country / document-type filters**, so this works directly for landscape checks like "has a feature we're about to file for already been published or granted in the last few years?"

---

## 2. Search by "combination," not a single document — `combos`

While `search` finds "the single document most similar to this sentence," `combos` (`/search/103/`) finds combinations where **multiple documents together** cover the whole claim. It's built for §103 (non-obviousness) analysis — instead of manually cross-checking whether "prior art A + B" renders a claim obvious, you get a narrowed-down candidate pool to start from.

```bash
pqai combos "battery management system with thermal runaway detection" -n 10
```

---

## 3. Search anchored on a specific patent — `prior-art`, `similar`

You don't need keywords at all — **just a patent number** is enough to search. `prior-art` filters to documents that predate the patent's filing date, which is exactly what you want for an invalidity search, while `similar` finds similar documents with no date restriction.

```bash
pqai prior-art US11868178B2 -n 10
```

Actually calling this on `US11868178B2` (a ring-form wearable-device patent), out of 500 candidates the top hits were:

```
1. CN218548276U — score 0.7606 — published 2023-02-28
   Wearable device (a wearable device with a conductive button assembly)
2. KR20170091346A — score 0.7472 — Ring type wearable device
3. WO2021017915A1 — score 0.7410 — Finger-ring monitoring device and monitoring system
```

The key point: knowing just one patent number is enough to pull up its overlapping prior-art pool — you never have to think up a search query yourself.

**Gotcha**: unlike `search`/`combos`, these two routes don't take date or country filters at all. Per the [official PQAI API docs](https://api.projectpq.ai/docs), `/prior-art/patent/` and `/similar/` only accept `pn`, `n`, `offset`, `index`, and `type` — there's no `-cc`, `-dtype`, `-after`, or `-before` here (`prior-art` already implicitly cuts off at the patent's own filing date, so a separate date filter wouldn't make sense; there's just no country filter at all). Trying one errors out immediately client-side (`flag provided but not defined`), before any request is even sent. If you need to filter by date/country, do it through `search`/`combos` with a text query instead.

---

## 4. Split a claim into elements and match each one against a document — `mapping` ⭐ (the most underrated feature)

`snippet`/`mapping` go beyond "this document is relevant to your query" — they slice the query into **individual claim elements (limitations)** and show which part of the document each one maps to. This is essentially an automated version of the "claim chart" that patent attorneys and examiners build by hand.

The key trick: if you feed the query in as a **claim broken into lines**, the API automatically splits it into elements and maps each one separately — something you can't tell just from reading the docs; you only discover it by actually sending a multi-line query.

```bash
Q='A drone comprising:
a rotor assembly configured to rotate about an axis and provide lift;
a camera coupled to the drone for capturing images of a fire;
a fire suppression module configured to release a fire suppressant material;
a controller configured to autonomously navigate the drone toward the fire.'

pqai mapping US10112730B2 -q "$Q"
```

Actual response (abridged):

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

It sliced the 4-element claim into separate items and, within the target patent (`US10112730B2`, a drone-related patent), paired each one with the most relevant matching passage. If you feed in a single element only (`-q "a rotor assembly..."`), only that one element gets mapped — so **the real trick to matching an entire claim at once is splitting it with semicolons/newlines.**

`snippet` is the simpler cousin: it returns a single passage most relevant to the whole query:

```bash
pqai snippet US10112730B2 -q "a rotor assembly configured to rotate about an axis and provide lift"
```
```json
{
  "snippet": "... (UAVs), fixed wing aircrafts, rotary wing aircrafts and helicopters. The one or more rotating parts of said vehicle may include rotor blades, propeller blades, turbine blades, jet/gas compressors, reciprocating engine, or similar parts. Receiving the instantaneous ..."
}
```

---

## 5. Predict CPC classification and examining unit (GAU) from text alone — `cpcs`, `gaus`

Knowing ahead of time which CPC subclass and which USPTO examining group (Art Unit) an invention is likely to be assigned to — even before filing — helps you narrow the scope of a prior-art search or gauge how tough examination might be.

```bash
pqai cpcs "an autonomous drone equipped with a fire suppressant tank and camera-based fire detection for aerial firefighting"
```

Actual top response (sorted by confidence, out of 61 total candidates):

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

Rather than just a single CPC code, it returns the **entire hierarchy** ("area conflagrations" → "delivery of fire extinguishing material by aircraft" → "spraying from the aircraft"), so you can understand the classification scheme even without knowing the code beforehand. Calling `gaus` on the same text:

```bash
pqai gaus "an autonomous drone equipped with a fire suppressant tank and camera-based fire detection for aerial firefighting"
# → ["3752", "2482", "3664"]
```

You immediately get 3 candidate USPTO Group Art Units, ranked by priority. If you just want the definition of a specific CPC code:

```bash
pqai cpc-def A62C3/0242
```

---

## 6. Drawings download without a token — `drawings`, `drawing`

Not the whole API requires a paid token. **Only the individual drawing-image download route is open without authentication** — it's documented in the README, but it's hard to be fully sure until you actually try calling it with no token, so we verified it directly.

```bash
pqai drawings US10112730B2          # list of drawings: ["1","2","3","4","5","6"]
pqai drawing US10112730B2 1 -o fig1.png    # PNG downloaded successfully with no token (321KB)
pqai drawing US10112730B2 1 -w 300 -o thumb1.png  # 300px-wide thumbnail
```

Even with `PQAI_API_KEY` removed from `.env`, the `drawing` command alone works fine. This means you can bulk-fetch or embed patent figures into slides/reports without burning through API credits.

---

## 7. Query the raw data itself — `patent`, `document`, `vector`, `dataset`

- `patent <pn>` — the patent's raw data: bibliographic info, claims, specification text, etc.
- `document <id>` — look up a document (patent or paper) in PQAI's own indexed database
- `vector <pn> <field>` — pull a patent's `cpcs` or `abstract` embedding vector directly (useful for your own similarity computations, clustering, etc.)
- `dataset -name <n> -n <i>` — fetch a sample from a benchmark/PoC dataset PQAI provides

These commands go one level deeper than the summarized information `search`-family commands return, and they're what you need when you want to do your own post-processing after search — custom re-ranking, vector-similarity computation, and so on.

---

## Note: pricing and call budget

The PQAI+ entry-level plan has a monthly call cap (e.g. roughly 20 calls for $20/month). All examples above are actual responses obtained by calling the live API through `pqai-cli`; if you're planning experiments that need lots of calls (e.g. a full invalidity search across an entire claim), budget your paid credits accordingly before you start.

For the complete command reference, see [`README.md`](README.md).
