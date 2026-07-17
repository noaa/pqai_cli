# PQAI Search Results: Prior Art Related to Apple's Image Wand

*[한국어](image-wand-search.kr.md)*

## Search conditions

- Command: `pqai search "<query>" -dtype filing -after 2023-01-01 -cc US -type patent -n 10`
- Query (English, functional description):
  > A method for generating a digital image from a user's rough hand-drawn sketch made with a finger or stylus input in a note-taking application, wherein the user encircles the sketch or a blank area with a selection gesture, and the system generates one or more candidate images based on the sketch and an associated text description or surrounding note content, allowing the user to select among multiple rendering styles such as sketch, illustration, and animation.
- Filters: filing date on/after 2023-01-01, country code US, document type patent
- Full request details (route, parameter table): [`image-wand-search-request.en.md`](./image-wand-search-request.en.md)

## Results (human-readable version)

 1. **US2025200831A1** — score 0.6821 — published 2025-06-19
    METHOD FOR AUTOMATICALLY GENERATING SKETCH IMAGE, APPARATUS FOR AUTOMATICALLY GENERATING SKETCH IMAGE USING THE METHOD, AND COMPUTER READABLE MEDIUM HAVING PROGRAM FOR PROCESSING THE METHOD
    Owner: Korea Advanced Institute of Science and Technology
    A method that extracts shape data from a color image and style data from a reference image, then combines the two to output a sketch image.

 2. **US2024221328A1** — score 0.6729 — published 2024-07-04
    Method and Device for Sketch-Based Placement of Virtual Objects
    Owner: **Apple Inc.**
    A method where the user sketches into a content-creation area, a 3D model is derived from it, and a computer-generated graphic object (virtual object) is placed/displayed over a real-world camera image. (Closer to AR/3D object placement, somewhat different in character from Image Wand's 2D image generation.)

 3. **US2025308125A1** — score 0.6680 — published 2025-10-02
    CUSTOMIZED ANIMATION FROM VIDEO
    Owner: Snap Inc.
    The user draws to select a specific area of a video, and a graphic element is created for that area with visual effects applied to generate a custom sticker/graphic.

 4. **US2024378780A1** — score 0.6670 — published 2024-11-14
    Method and system for generating image
    Owner: Naver Webtoon Ltd.
    Identifies a source image based on content selected by the user, and applies that content's drawing style to generate a resulting image.

 5. **US2024320867A1** — score 0.6631 — published 2024-09-26
    Iterative Image Generation From Text
    Owner: Sony Interactive Entertainment Inc.
    A text-to-image generator automatically identifies additional descriptors from a generated image, feeds them back into the prompt, and iteratively refines/generates images.

 6. **US2025173962A1** — score 0.6582 — published 2025-05-29
    METHOD AND SYSTEM FOR CREATING 3D OBJECTS FROM ROUGHLY DRAWN SKETCH AND TEXT
    Owner: Korea Electronics Technology Institute
    Takes an unstructured sketch and text as input, generates a 2D image, then creates a 3D object based on it.

 7. **US2023343000A1** — score 0.6499 — published 2023-10-26
    Method and apparatus for picture generation and storage medium
    Owner: Boe Technology Group Co., Ltd.
    Acquires content elements included in a target picture's design, parses them into graphic elements according to drawing commands, then combines them to produce a final picture.

 8. **US2025117990A1** — score 0.6364 — published 2025-04-10
    SCRIBBLE-TO-VECTOR IMAGE GENERATION
    Owner: ADOBE INC.
    Takes a sketch input depicting an object, processes it as sketch guidance, and uses an image-generation model to produce a synthetic image depicting that object.

 9. **US2024273308A1** — score 0.6318 — published 2024-08-15
    System and method for visual content generation and iteration
    Owner: Toyota Research Institute, Inc.
    A creative-support system where a generative language model produces various texts, which then drive a generative visual model to create multiple images.

10. **US12347013B2** — score 0.6143 — published 2025-07-01
    Animated custom sticker creation
    Owner: Snap Inc.
    When the user selects a specific area of an image, a graphic element is generated with visual effects and animation patterns applied to create a custom animated sticker.

## Summary / analysis

- With this query, only **1 out of the top 10 results was an Apple-owned patent (US2024221328A1)**, and even that one is closer to "sketch → 3D virtual object placement (AR)" than Image Wand's "sketch → 2D image generation."
- The rest are adjacent-field prior art in **sketch/text-based image and 3D generation** from Adobe, Snap, Sony, Naver Webtoon, KAIST, the Korea Electronics Technology Institute, and others.
- To target Apple's Image Wand (Notes app, Markup menu, 3 rendering styles: Sketch/Illustration/Animation) more precisely:
  - Add more specific phrasing to the query, such as "note-taking application Markup menu," "multiple style presets," "on-device generative model"
  - Re-search with `-lq`, marking `US2024221328A1` as `relevant`, to steer results toward the Apple direction
  - There's no `owner`/`assignee` filter, so if you only want Apple patents you'd need to post-process/filter the results yourself

## Raw JSON

See [`image-wand-search.json`](./image-wand-search.json) for the full raw response.
