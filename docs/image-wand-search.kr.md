# PQAI 검색 결과: Apple Image Wand 관련 선행기술

*[English](image-wand-search.en.md)*

## 검색 조건

- 명령: `pqai search "<query>" -dtype filing -after 2023-01-01 -cc US -type patent -n 10`
- 쿼리 (영어, 기능 설명):
  > A method for generating a digital image from a user's rough hand-drawn sketch made with a finger or stylus input in a note-taking application, wherein the user encircles the sketch or a blank area with a selection gesture, and the system generates one or more candidate images based on the sketch and an associated text description or surrounding note content, allowing the user to select among multiple rendering styles such as sketch, illustration, and animation.
- 필터: 출원일(filing) 2023-01-01 이후, 국가 코드 US, 문서 유형 patent
- 전체 요청 상세(라우트, 파라미터 표): [`image-wand-search-request.kr.md`](./image-wand-search-request.kr.md)

## 결과 (사람이 읽기 쉬운 버전)

 1. **US2025200831A1** — score 0.6821 — 공개일 2025-06-19
    METHOD FOR AUTOMATICALLY GENERATING SKETCH IMAGE, APPARATUS FOR AUTOMATICALLY GENERATING SKETCH IMAGE USING THE METHOD, AND COMPUTER READABLE MEDIUM HAVING PROGRAM FOR PROCESSING THE METHOD
    소유자: Korea Advanced Institute of Science and Technology
    색상 이미지에서 형태(shape) 데이터를 추출하고 참조 이미지에서 스타일 데이터를 추출하여, 이 둘을 조합해 스케치 이미지를 출력하는 방법.

 2. **US2024221328A1** — score 0.6729 — 공개일 2024-07-04
    Method and Device for Sketch-Based Placement of Virtual Objects
    소유자: **Apple Inc.**
    사용자가 콘텐츠 생성 영역에 스케치를 입력하면, 이를 바탕으로 3D 모델을 획득하고, 카메라로 얻은 실사 이미지 위에 컴퓨터 생성 그래픽 객체(가상 객체)를 배치/표시하는 방법. (AR/3D 객체 배치에 가까우며 Image Wand의 2D 이미지 생성과는 다소 결이 다름)

 3. **US2025308125A1** — score 0.6680 — 공개일 2025-10-02
    CUSTOMIZED ANIMATION FROM VIDEO
    소유자: Snap Inc.
    비디오의 특정 영역을 사용자가 드로잉으로 선택하면, 해당 영역으로 그래픽 요소를 만들고 시각 효과를 적용해 커스텀 스티커/그래픽을 생성.

 4. **US2024378780A1** — score 0.6670 — 공개일 2024-11-14
    Method and system for generating image
    소유자: Naver Webtoon Ltd.
    사용자가 선택한 콘텐츠를 기반으로 소스 이미지를 특정하고, 해당 콘텐츠의 그림체(drawing style)를 적용해 결과 이미지를 생성.

 5. **US2024320867A1** — score 0.6631 — 공개일 2024-09-26
    Iterative Image Generation From Text
    소유자: Sony Interactive Entertainment Inc.
    텍스트-이미지 생성기가 만든 이미지에서 추가 디스크립터를 자동으로 식별해 프롬프트에 반영, 반복적으로 이미지를 개선/생성.

 6. **US2025173962A1** — score 0.6582 — 공개일 2025-05-29
    METHOD AND SYSTEM FOR CREATING 3D OBJECTS FROM ROUGHLY DRAWN SKETCH AND TEXT
    소유자: Korea Electronics Technology Institute
    비정형 스케치와 텍스트를 입력받아 2D 이미지를 생성하고, 이를 기반으로 3D 객체를 생성.

 7. **US2023343000A1** — score 0.6499 — 공개일 2023-10-26
    Method and apparatus for picture generation and storage medium
    소유자: Boe Technology Group Co., Ltd.
    타겟 픽처의 디자인에 포함된 콘텐츠 요소들을 획득하고, 드로잉 명령에 따라 그래픽 요소로 파싱한 뒤 결합하여 최종 픽처를 생성.

 8. **US2025117990A1** — score 0.6364 — 공개일 2025-04-10
    SCRIBBLE-TO-VECTOR IMAGE GENERATION
    소유자: ADOBE INC.
    객체를 묘사하는 스케치 입력을 받아 스케치 가이던스로 처리한 뒤, 이미지 생성 모델을 이용해 해당 객체를 묘사하는 합성 이미지를 생성.

 9. **US2024273308A1** — score 0.6318 — 공개일 2024-08-15
    System and method for visual content generation and iteration
    소유자: Toyota Research Institute, Inc.
    생성형 언어 모델로 다양한 텍스트를 생성하고, 이를 바탕으로 생성형 비주얼 모델로 여러 이미지를 만드는 창작 과정 지원 시스템.

10. **US12347013B2** — score 0.6143 — 공개일 2025-07-01
    Animated custom sticker creation
    소유자: Snap Inc.
    사용자가 이미지의 특정 영역을 선택하면 그래픽 요소를 생성하고 시각 효과 및 애니메이션 패턴을 적용해 커스텀 애니메이션 스티커를 생성.

## 요약 / 분석

- 이번 쿼리로는 **Apple 소유 특허는 상위 10건 중 1건(US2024221328A1)** 만 노출되었으며, 이는 Image Wand의 "스케치 → 2D 이미지 생성"보다는 "스케치 → 3D 가상 객체 배치(AR)"에 가까운 특허임.
- 나머지는 Adobe, Snap, Sony, Naver Webtoon, KAIST, 한국전자기술연구원 등 **스케치/텍스트 기반 이미지·3D 생성**이라는 인접 기술 분야의 유사 선행기술.
- Apple Image Wand(Notes 앱, Markup 메뉴, 3가지 렌더링 스타일: Sketch/Illustration/Animation)를 더 정확히 겨냥하려면:
  - 쿼리에 "note-taking application Markup menu", "multiple style presets", "on-device generative model" 등 구체적 표현 추가
  - `-lq`로 `US2024221328A1`을 `relevant`에 지정해 Apple 관련 방향으로 재검색
  - `owner`/`assignee` 필터가 없으므로, Apple 특허만 보려면 결과를 후처리로 필터링 필요

## 원본 JSON

전체 원본 응답은 [`image-wand-search.json`](./image-wand-search.json) 참고.
