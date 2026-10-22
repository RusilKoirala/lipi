// Loads the real Lipi interpreter (Go compiled to WebAssembly) and runs code with it.
import { href } from "./site";

type GoInstance = {
  importObject: WebAssembly.Imports;
  run(instance: WebAssembly.Instance): Promise<void>;
};

declare global {
  interface Window {
    Go?: new () => GoInstance;
    runLipi?: (code: string) => { output: string; errors: string[] };
  }
}

export type LipiResult = {
  output: string;
  errors: string[];
  ms: number;
};

let ready: Promise<void> | null = null;

function loadScript(src: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const s = document.createElement("script");
    s.src = src;
    s.onload = () => resolve();
    s.onerror = () => reject(new Error(`Could not load ${src}`));
    document.head.append(s);
  });
}

async function instantiate(go: GoInstance): Promise<WebAssembly.Instance> {
  const url = href("/lipi.wasm");
  try {
    const { instance } = await WebAssembly.instantiateStreaming(fetch(url), go.importObject);
    return instance;
  } catch {
    // Some servers send .wasm without the right content type; fall back to a plain fetch.
    const bytes = await (await fetch(url)).arrayBuffer();
    const { instance } = await WebAssembly.instantiate(bytes, go.importObject);
    return instance;
  }
}

function reset() {
  ready = null;
  window.runLipi = undefined;
}

/** Load the interpreter once. Safe to call many times. */
export function loadLipi(): Promise<void> {
  if (ready) return ready;

  ready = (async () => {
    if (!window.Go) await loadScript(href("/wasm_exec.js"));
    const go = new window.Go!();
    const instance = await instantiate(go);

    // If the Go program ever exits, load a fresh one next time.
    go.run(instance).finally(reset);

    for (let i = 0; i < 200 && !window.runLipi; i++) {
      await new Promise((r) => setTimeout(r, 10));
    }
    if (!window.runLipi) throw new Error("The Lipi interpreter did not start.");
  })();

  ready.catch(reset);
  return ready;
}

export async function runLipi(code: string): Promise<LipiResult> {
  await loadLipi();
  const t0 = performance.now();
  try {
    const r = window.runLipi!(code);
    return { output: r.output, errors: Array.from(r.errors ?? []), ms: performance.now() - t0 };
  } catch {
    // The browser ran out of stack inside Go. Start a fresh interpreter for the next run.
    reset();
    return {
      output: "",
      errors: ["Galti: program dherai gahirai pugyo, interpreter feri suru garincha"],
      ms: performance.now() - t0,
    };
  }
}

/** "Galti: line 3 ma ..." -> 3 */
export function errorLine(message: string): number | null {
  const m = /line (\d+)/.exec(message);
  return m ? Number(m[1]) : null;
}
