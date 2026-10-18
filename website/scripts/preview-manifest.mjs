// Writes the website's component navigation from the Gallery's demonstrations,
// so the page labels and preview slugs stay in one place. Go must be installed.
//
//   bun run preview:manifest
import { execFile } from "node:child_process"
import { writeFile } from "node:fs/promises"
import { resolve, dirname } from "node:path"
import { fileURLToPath } from "node:url"
import { promisify } from "node:util"

const run = promisify(execFile)
const websiteDir = resolve(dirname(fileURLToPath(import.meta.url)), "..")
const sourceDir = resolve(websiteDir, "../source")
const output = resolve(websiteDir, "app/lib/preview-manifest.json")

try {
  const { stdout } = await run(
    "go",
    ["run", "./examples/gallery/web", "-manifest"],
    {
      cwd: sourceDir,
      env: { ...process.env, GOWORK: "off" },
      maxBuffer: 1 << 24,
    }
  )
  await writeFile(output, stdout)
  console.log(`Wrote ${output}`)
} catch (error) {
  const e = error
  console.error(
    `preview:manifest: go run ./examples/gallery/web -manifest failed: ${
      e.stderr || e.message
    }`
  )
  process.exitCode = 1
}
