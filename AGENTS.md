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

Update the owning subject, affected overview and usage instructions alongside mechanism changes. Record completed user-visible API and gallery changes under Unreleased in CHANGELOG.md; pure refactors and tests do not need release notes. Keep documentation and code on the same branch and reconcile definitions and verification evidence during merges.

Keep ARCHITECTURE.md within 16 KiB and each subject within 24 KiB; split by responsibility and update Subjects when needed.

## Verification

Follow [development checks](docs/build.md#development-checks): format, test, race-check, vet and build. Interface changes need target-platform native review, not just headless tests. Do not launch external terminals or Agents as a substitute.

For documentation, verify size budgets, Subjects coverage, relative links and heading anchors, and changelog structure. Run `git diff --check` from the repository root. Do not commit, tag, publish or push without explicit authorization.
