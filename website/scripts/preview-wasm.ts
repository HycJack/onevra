// Compiles the website's component preview, the Gallery's demonstrations for
// browsers (source/examples/gallery/web), and serves it in development and in
// the Pages artifact. Go must be installed; the website build fails without it.
import { execFile } from "node:child_process"
import { createHash } from "node:crypto"
import { mkdtemp, readFile, rm } from "node:fs/promises"
import { tmpdir } from "node:os"
import { join, resolve } from "node:path"
import { fileURLToPath } from "node:url"
import { promisify } from "node:util"
import type { Plugin } from "vite"

const run = promisify(execFile)

const websiteDir = fileURLToPath(new URL("../", import.meta.url))
const sourceDir = resolve(websiteDir, "../source")
const pagePath = resolve(websiteDir, "preview/index.html")
const base = (process.env.BASE_PATH ?? "/onevra/").replace(/^\/+|\/+$/g, "")
const basePath = base ? `/${base}/` : "/"

interface PreviewBuild {
  wasm: Buffer
  loader: string
}

let building: Promise<PreviewBuild> | undefined

/** Compiles the app to WebAssembly and reads Go's matching loader. */
function build(): Promise<PreviewBuild> {
  building ??= compile().catch((error) => {
    building = undefined
    throw error
  })
  return building
}

async function compile(): Promise<PreviewBuild> {
  const dir = await mkdtemp(join(tmpdir(), "onevra-preview-"))
  const out = join(dir, "app.wasm")
  try {
    await run(
      "go",
      [
        "build",
        "-tags=mygo_noinspector",
        "-trimpath",
        "-ldflags=-s -w",
        "-o",
        out,
        "./examples/gallery/web",
      ],
      {
        cwd: sourceDir,
        env: {
          ...process.env,
          GOOS: "js",
          GOARCH: "wasm",
          CGO_ENABLED: "0",
          GOWORK: "off",
        },
        maxBuffer: 1 << 24,
      }
    )
    const [wasm, wasmLoader] = await Promise.all([readFile(out), loader()])
    return { wasm, loader: wasmLoader }
  } catch (error) {
    const e = error as { stderr?: string; message: string }
    throw new Error(
      `onevra-preview: go build of ./examples/gallery/web failed: ${e.stderr || e.message}`
    )
  } finally {
    await rm(dir, { recursive: true, force: true })
  }
}

/** Go's loader of WebAssembly, of the Go that compiles the app. */
async function loader(): Promise<string> {
  const { stdout: goroot } = await run("go", ["env", "GOROOT"], {
    cwd: sourceDir,
  })
  return readFile(join(goroot.trim(), "lib", "wasm", "wasm_exec.js"), "utf8")
}

/** The page loading the app, with the file name of its WebAssembly. */
async function page(wasm: string): Promise<string> {
  return (await readFile(pagePath, "utf8")).replace("__WASM__", wasm)
}

const types: Record<string, string> = {
  html: "text/html; charset=utf-8",
  js: "text/javascript; charset=utf-8",
  wasm: "application/wasm",
}

export function previewWasm(): Plugin {
  return {
    name: "onevra-preview-wasm",
    applyToEnvironment: (environment) => environment.name === "client",
    async generateBundle() {
      const { wasm, loader } = await build()
      // Named by its content, so that browsers keep it while it lasts.
      const name = `app.${createHash("sha256").update(wasm).digest("hex").slice(0, 12)}.wasm`
      this.emitFile({
        type: "asset",
        fileName: `preview/${name}`,
        source: wasm,
      })
      this.emitFile({
        type: "asset",
        fileName: "preview/wasm_exec.js",
        source: loader,
      })
      this.emitFile({
        type: "asset",
        fileName: "preview/index.html",
        source: await page(name),
      })
    },
    configureServer(server) {
      // Recompile the app when its Go source changes; the next reload of the
      // preview picks it up.
      server.watcher.add(sourceDir)
      server.watcher.on("all", (event, file) => {
        if (
          (event === "change" || event === "add" || event === "unlink") &&
          file.endsWith(".go")
        ) {
          building = undefined
        }
      })

      server.middlewares.use(async (request, response, next) => {
        const { pathname } = new URL(request.url ?? "/", "http://localhost")
        const prefix = `${basePath}preview/`
        if (!pathname.startsWith(prefix)) return next()

        const file = pathname.slice(prefix.length) || "index.html"
        try {
          let source: string | Buffer
          let kind: string
          if (file === "index.html") {
            // A fixed name in development keeps the page readable; it is
            // served with no-store below. The build hashes the name instead.
            source = await page("app.wasm")
            kind = "html"
          } else if (file === "wasm_exec.js") {
            source = (await build()).loader
            kind = "js"
          } else if (file === "app.wasm") {
            source = (await build()).wasm
            kind = "wasm"
          } else {
            return next()
          }
          response.setHeader("Content-Type", types[kind])
          response.setHeader("Cache-Control", "no-store")
          response.end(source)
        } catch (error) {
          next(error)
        }
      })
    },
  }
}
