import { useCallback, useEffect, useRef, useState } from "react";
import { Check, Link2, Loader2, Play, RotateCcw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { EXAMPLES, getExample } from "@/lib/examples";
import { errorLine, loadLipi, runLipi, type LipiResult } from "@/lib/lipi";
import { href } from "@/lib/site";
import { cn } from "@/lib/utils";
import { Editor, type EditorHandle } from "./Editor";

type Status = "idle" | "loading" | "ready" | "running" | "failed";

const DRAFT_KEY = "lipi-playground-draft";

// UTF-8 safe base64url, so Devanagari survives in share links.
function encode(code: string): string {
  const bytes = new TextEncoder().encode(code);
  let bin = "";
  bytes.forEach((b) => (bin += String.fromCharCode(b)));
  return btoa(bin).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

function decode(s: string): string | null {
  try {
    const bin = atob(s.replace(/-/g, "+").replace(/_/g, "/"));
    return new TextDecoder().decode(Uint8Array.from(bin, (c) => c.charCodeAt(0)));
  } catch {
    return null;
  }
}

function readHash(): { code: string; exampleId: string | null } | null {
  const params = new URLSearchParams(window.location.hash.slice(1));
  const code = params.get("code");
  if (code) {
    const decoded = decode(code);
    if (decoded != null) return { code: decoded, exampleId: null };
  }
  const ex = getExample(params.get("example"));
  if (ex) return { code: ex.code, exampleId: ex.id };
  return null;
}

const isMac = typeof navigator !== "undefined" && /Mac|iPhone|iPad/.test(navigator.platform || navigator.userAgent);

type Props = {
  /** Small version for the landing page. */
  compact?: boolean;
};

export default function Playground({ compact = false }: Props) {
  const first = EXAMPLES[0];
  const editor = useRef<EditorHandle>(null);
  const code = useRef(first.code);

  const [exampleId, setExampleId] = useState<string>(first.id);
  const [fileName, setFileName] = useState(first.file);
  const [status, setStatus] = useState<Status>("idle");
  const [result, setResult] = useState<LipiResult | null>(null);
  const [errLine, setErrLine] = useState<number | null>(null);
  const [copied, setCopied] = useState(false);
  const [modKey, setModKey] = useState("Ctrl");

  const run = useCallback(async () => {
    setStatus("running");
    setErrLine(null);
    try {
      const r = await runLipi(code.current);
      setResult(r);
      setErrLine(r.errors[0] ? errorLine(r.errors[0]) : null);
      setStatus("ready");
    } catch {
      setStatus("failed");
    }
  }, []);

  // Pick the starting code (link, saved draft or first example), then load WebAssembly and run once.
  useEffect(() => {
    setModKey(isMac ? "⌘" : "Ctrl");

    if (!compact) {
      const fromHash = readHash();
      let draft: string | null = null;
      try {
        draft = localStorage.getItem(DRAFT_KEY);
      } catch {
        /* storage blocked */
      }
      const start = fromHash ?? (draft ? { code: draft, exampleId: null } : null);
      if (start && start.code !== code.current) {
        code.current = start.code;
        editor.current?.setValue(start.code);
        const ex = getExample(start.exampleId);
        setExampleId(ex ? ex.id : "");
        setFileName(ex ? ex.file : "mero-program.lip");
      }
    }

    setStatus("loading");
    loadLipi()
      .then(run)
      .catch(() => setStatus("failed"));
  }, [compact, run]);

  const onChange = useCallback(
    (value: string) => {
      code.current = value;
      if (!compact) {
        try {
          localStorage.setItem(DRAFT_KEY, value);
        } catch {
          /* storage blocked */
        }
      }
    },
    [compact],
  );

  function loadExample(id: string) {
    const ex = getExample(id);
    if (!ex) return;
    setExampleId(ex.id);
    setFileName(ex.file);
    code.current = ex.code;
    editor.current?.setValue(ex.code);
    history.replaceState(null, "", `#example=${ex.id}`);
    run();
  }

  async function share() {
    const url = `${location.origin}${href("/playground")}#code=${encode(code.current)}`;
    history.replaceState(null, "", `#code=${encode(code.current)}`);
    try {
      await navigator.clipboard.writeText(url);
    } catch {
      /* clipboard blocked: the address bar still has the link */
    }
    setCopied(true);
    setTimeout(() => setCopied(false), 1800);
  }

  const lines = result?.output ? result.output.replace(/\n$/, "").split("\n") : [];
  const busy = status === "loading" || status === "running";

  const runButton = (
    <Button
      onClick={run}
      disabled={status === "loading" || status === "failed"}
      size={compact ? "sm" : "default"}
      className="gap-2 bg-crimson text-white hover:bg-crimson/90"
    >
      {busy ? <Loader2 className="animate-spin" /> : <Play className="fill-current" />}
      Run
      <kbd className="hidden rounded border border-white/35 px-1 font-mono text-[10px] leading-4 font-normal opacity-85 sm:inline">
        {modKey} ↵
      </kbd>
    </Button>
  );

  return (
    <div
      className={cn(
        "flex min-h-0 flex-col overflow-hidden rounded-xl border bg-card shadow-[0_1px_0_0_var(--border),0_24px_48px_-32px_oklch(0.2_0.02_268/0.35)]",
        compact ? "h-[34rem]" : "h-full",
      )}
    >
      {/* Toolbar */}
      <div className="flex flex-wrap items-center gap-2 border-b px-3 py-2">
        <span className="flex items-center gap-2 font-mono text-xs text-muted-foreground">
          <span className="size-2 rounded-[2px] bg-crimson" aria-hidden />
          {fileName}
        </span>
        <div className="ml-auto flex flex-wrap items-center justify-end gap-2">
          {!compact && (
            <>
              <Select value={exampleId} onValueChange={loadExample}>
                <SelectTrigger className="w-[9.5rem] sm:w-[11.5rem]" aria-label="Load an example">
                  <SelectValue placeholder="Load an example" />
                </SelectTrigger>
                <SelectContent>
                  {EXAMPLES.map((ex) => (
                    <SelectItem key={ex.id} value={ex.id}>
                      {ex.title}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Button variant="outline" onClick={share} className="gap-2" aria-live="polite">
                {copied ? <Check /> : <Link2 />}
                <span className="hidden sm:inline">{copied ? "Link copied" : "Share"}</span>
              </Button>
            </>
          )}
          {compact && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => {
                code.current = first.code;
                editor.current?.setValue(first.code);
                run();
              }}
              aria-label="Reset the code"
              className="text-muted-foreground"
            >
              <RotateCcw />
            </Button>
          )}
          {runButton}
        </div>
      </div>

      {/* Editor + output */}
      <div className={cn("grid min-h-0 flex-1 grid-cols-1", compact ? "grid-rows-[minmax(0,1fr)_auto]" : "grid-rows-[minmax(0,1.4fr)_minmax(0,1fr)] md:grid-cols-[minmax(0,1.45fr)_minmax(0,1fr)] md:grid-rows-1")}>
        <Editor
          ref={editor}
          initialValue={first.code}
          onChange={onChange}
          onRun={run}
          errorLine={errLine}
          className="min-h-0"
        />

        <section
          aria-label="Output"
          aria-live="polite"
          className={cn(
            "flex min-h-0 flex-col border-t bg-paper",
            compact ? "max-h-[9.5rem] min-h-[7.5rem]" : "md:border-t-0 md:border-l",
          )}
        >
          <div className="flex items-center justify-between px-4 pt-3 pb-1 text-[11px] font-medium tracking-[0.08em] text-muted-foreground uppercase">
            <span>Output</span>
            <StatusText status={status} result={result} />
          </div>
          <div className="min-h-0 flex-1 overflow-auto px-4 pb-4 font-mono text-[13.5px] leading-6">
            {status === "loading" && !result && (
              <p className="text-muted-foreground">Loading the interpreter…</p>
            )}
            {status === "failed" && (
              <p className="text-crimson">
                The interpreter could not load. Check your connection and refresh the page.
              </p>
            )}
            {lines.map((line, i) => (
              <div key={i} className="break-words whitespace-pre-wrap">
                <span className="mr-2 text-muted-foreground/60 select-none">›</span>
                {line}
              </div>
            ))}
            {result && !result.errors.length && !lines.length && (
              <p className="font-sans text-sm text-muted-foreground">
                The program ran but printed nothing. Use <code className="font-mono">bhana(...)</code> to
                print.
              </p>
            )}
            {result?.errors.map((e, i) => {
              const line = errorLine(e);
              return (
                <div
                  key={i}
                  className="mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 rounded-md border border-crimson/25 bg-crimson-soft px-3 py-2 text-crimson"
                >
                  <span className="break-words whitespace-pre-wrap">{e}</span>
                  {line && (
                    <button
                      type="button"
                      onClick={() => editor.current?.jumpTo(line)}
                      className="font-sans text-xs underline underline-offset-2 hover:no-underline"
                    >
                      Go to line {line}
                    </button>
                  )}
                </div>
              );
            })}
          </div>
        </section>
      </div>
    </div>
  );
}

function StatusText({ status, result }: { status: Status; result: LipiResult | null }) {
  if (status === "loading") return <span>Starting…</span>;
  if (status === "running") return <span>Running…</span>;
  if (status === "failed") return <span className="text-crimson">Offline</span>;
  if (!result) return null;
  if (result.errors.length) return <span className="text-crimson normal-case tracking-normal">● Galti</span>;
  const ms = result.ms < 1 ? "<1" : Math.round(result.ms);
  return (
    <span className="normal-case tracking-normal tabular-nums">
      <span className="text-emerald-600 dark:text-emerald-400">●</span> Chalyo · {ms} ms
    </span>
  );
}
