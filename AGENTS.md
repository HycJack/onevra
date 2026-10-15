# Onevra

Onevra is an independent Go component library built on MyGo. Read [ARCHITECTURE.md](ARCHITECTURE.md) in full before changing code, then read the relevant Subjects.

## Boundaries

- Keep Go source, tests, examples, module files and application packaging configuration in `source/`. Run Go commands from that module directory; keep documentation at the repository root and in `docs/`. The module path is `github.com/oneo-me/onevra/source`; follow the [publication versioning rules](docs/build.md#publication).

- When Onevra implements a component, use that component instead of MyGo's equivalent styled widget in application and gallery UI. MyGo primitives and unstyled bases remain implementation building blocks; dedicated MyGo component demonstrations are the only styled-widget exception. Check `source/components/` and `source/shells/` before choosing a control.
- Use all MyGo plugins directly, including Glass and future plugins; do not implement Onevra replacements.
- Compose MyGo primitives and unstyled bases; do not create a parallel renderer, input system or theme manager.
- Keep public components product-neutral. Applications own business state, localized text, persistence and asynchronous work.
- Prefer generic content slots for shared component customization over business-specific fields or variants. Keep domain-specific content and decisions in the application; preserve existing convenience APIs when adding slots.
- Read upstream public signatures before using MyGo APIs. MyGo documentation: https://mygo.egoist.dev/docs; Markdown index: https://mygo.egoist.dev/llms.txt; complete Markdown: https://mygo.egoist.dev/llms-full.txt; individual pages: `/docs/*.md`. Version-matched documentation and implementation live in `docs/` and `ui/` under the directory returned by `GOWORK=off go list -m -f '{{.Dir}}' github.com/egoist/mygo` from `source/`; prefer these over the latest website when versions differ. Keep examples runnable independently of consumer applications.
- On every MyGo upgrade, read the release notes and compare upstream changes against Onevra's own components and Shell, especially adapted collection behavior, input, accessibility, animation and native appearance. Adopt applicable fixes through upstream primitives rather than duplicating runtime internals; preserve intentional Onevra interaction contracts. Proactively run `GOWORK=off go generate ./examples/gallery/internal/gallery` from `source/` on every upgrade and review the generated inventory diff; do not rely on packaging to refresh it. Supplement the gallery catalog with runnable examples for newly added components and window APIs, including methods and plugins the generator cannot discover; provide localized labels, lifecycle cleanup and relevant interaction coverage. Record an explicit reason for any intentional exclusion in the decision map. Add regression tests where Onevra's own behavior changes, and report review conclusions and native-verification gaps.
- Use English for documentation and comments, except maintained translations and intentional multilingual fixtures.
- Keep README.md and README.zh-CN.md equivalent with reciprocal language links.

## Maintenance

Update documentation in the same change only when its reader-needed facts change. Implementation edits alone do not require prose changes.

- README.md and README.zh-CN.md cover the project introduction, basic usage and contributor onboarding. Link to deeper material instead of expanding the quick start.
- ARCHITECTURE.md and docs/ explain system boundaries, core flows, design rationale and important constraints. Keep property definitions, defaults and local implementation details in code; do not duplicate them in prose or create a document per task.
- skills/onevra/ guides Agents building applications on Onevra. Update it when basic usage, component replacement decisions, testing conventions or packaging requirements change; keep it concise and link to the module source for details.
- CHANGELOG.md records only completed outcomes visible to library consumers under Unreleased. Describe supported usage or public-contract changes, not internal refactors, source layout, tests or tooling. Consolidate evolving unreleased work; do not present an unreleased experiment as a removed released feature.
- Keep adjacent changelog bullets together without blank lines; retain blank lines around headings.
- Update summaries, changelog entries and release notes describe only what was added, changed or removed. Unless the user explicitly requests it, omit Verification sections and all testing or validation details, including commands, results and verification gaps. Required checks still run as part of the work.

Keep documentation and code on the same branch and reconcile definitions during merges. Each explanation has one owning location; link rather than repeat it.

The 16 KiB overview and 24 KiB subject budgets are ceilings, not targets. Remove redundancy before splitting; add subjects only for distinct, lasting design questions. Keep Subjects and references current.

## Verification

Scale [development checks](docs/build.md#development-checks) to the change. While iterating, run only the affected packages, narrowed with `-run`; before handing off, run formatting, vet, the build and `go test ./...` once. Leave the full race check, native tests and other runs over a minute to the user: list the commands and state that they were not run. Documentation-only changes need no Go checks.

Add tests for public component behavior and contracts and for fixed bugs; extend existing tests before creating new ones. Gallery demonstrations, copy, styling adjustments and removals need no new tests. Keep tests fast: avoid real-time waits and repeated full-page frames.

Interface changes need target-platform native review. Headless tests cannot replace it, so report what remains unreviewed rather than adding tests to compensate. Do not launch external terminals or Agents as a substitute.

For documentation, verify size budgets, Subjects coverage, relative links and heading anchors, and changelog structure. Run `git diff --check` from the repository root. Do not commit, tag, publish or push without explicit authorization.

## Release requests

An explicit request to publish a named version authorizes the preparation, release commit, annotated tag and pushes described in [Authorized publication](docs/build.md#authorized-publication); execute that workflow without asking again for each step. Requests to prepare a release, edit documentation or record the release rules do not authorize publication. Pushing the release tag also authorizes the [tag-triggered workflow](.github/workflows/release.yml) to create its GitHub Release page. Application archives require separate inclusion in the requested scope.
