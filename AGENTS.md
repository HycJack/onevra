# oneo-go

oneo-go is an independent Go component library built on MyGo. Read [ARCHITECTURE.md](ARCHITECTURE.md) in full before changing code, then read the relevant Subjects.

## Boundaries

- Keep Go source, tests, examples, module files and MyGo configuration in `source/`. Run Go commands from that module directory; keep documentation at the repository root and in `docs/`.

- Compose MyGo's built-in widgets; do not create a parallel renderer, input system or theme manager.
- Keep public components product-neutral. Applications own business state, localized text, persistence and asynchronous work.
- Read upstream public signatures before using MyGo APIs. Keep examples runnable independently of consumer applications.
- Use English for documentation and comments, except maintained translations and intentional multilingual fixtures.
- Keep README.md and README.zh-CN.md equivalent with reciprocal language links.

## Maintenance

Update documentation in the same change only when its reader-needed facts change. Implementation edits alone do not require prose changes.

- README.md and README.zh-CN.md cover the project introduction, basic usage and contributor onboarding. Link to deeper material instead of expanding the quick start.
- ARCHITECTURE.md and docs/ explain system boundaries, core flows, design rationale and important constraints. Keep property definitions, defaults and local implementation details in code; do not duplicate them in prose or create a document per task.
- CHANGELOG.md records only completed outcomes visible to library consumers under Unreleased. Describe supported usage or public-contract changes, not internal refactors, source layout, tests or tooling. Consolidate evolving unreleased work; do not present an unreleased experiment as a removed released feature.

Keep documentation and code on the same branch and reconcile definitions during merges. Each explanation has one owning location; link rather than repeat it.

The 16 KiB overview and 24 KiB subject budgets are ceilings, not targets. Remove redundancy before splitting; add subjects only for distinct, lasting design questions. Keep Subjects and references current.

## Verification

Follow [development checks](docs/build.md#development-checks): format, test, race-check, vet and build. Interface changes need target-platform native review, not just headless tests. Do not launch external terminals or Agents as a substitute.

For documentation, verify size budgets, Subjects coverage, relative links and heading anchors, and changelog structure. Run `git diff --check` from the repository root. Do not commit, tag, publish or push without explicit authorization.
