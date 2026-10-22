import { forwardRef, useEffect, useImperativeHandle, useRef, useState } from "react";
import { closeBrackets, closeBracketsKeymap } from "@codemirror/autocomplete";
import { defaultKeymap, history, historyKeymap, indentWithTab } from "@codemirror/commands";
import { bracketMatching, indentOnInput, syntaxHighlighting } from "@codemirror/language";
import { EditorState, Prec, StateEffect, StateField } from "@codemirror/state";
import {
  Decoration,
  type DecorationSet,
  drawSelection,
  EditorView,
  highlightActiveLine,
  highlightActiveLineGutter,
  keymap,
  lineNumbers,
} from "@codemirror/view";
import { cn } from "@/lib/utils";
import { lipiHighlight, lipiLanguage } from "./lipi-language";

export type EditorHandle = {
  setValue(value: string): void;
  getValue(): string;
  jumpTo(line: number): void;
  focus(): void;
};

type Props = {
  initialValue: string;
  onChange?: (value: string) => void;
  onRun?: () => void;
  errorLine?: number | null;
  className?: string;
  ariaLabel?: string;
};

// Red band on the line a Galti points at. Cleared as soon as the code changes.
const setErrorLine = StateEffect.define<number | null>();
const errorLineField = StateField.define<DecorationSet>({
  create: () => Decoration.none,
  update(deco, tr) {
    for (const e of tr.effects) {
      if (e.is(setErrorLine)) {
        const n = e.value;
        if (n == null || n < 1 || n > tr.state.doc.lines) return Decoration.none;
        const line = tr.state.doc.line(n);
        return Decoration.set([Decoration.line({ class: "cm-errorLine" }).range(line.from)]);
      }
    }
    return tr.docChanged ? Decoration.none : deco;
  },
  provide: (f) => EditorView.decorations.from(f),
});

const theme = EditorView.theme({
  "&": { height: "100%", fontSize: "14px", backgroundColor: "transparent", color: "var(--foreground)" },
  ".cm-scroller": { fontFamily: "var(--font-mono)", lineHeight: "1.7", overflow: "auto" },
  ".cm-content": { padding: "14px 0", caretColor: "var(--foreground)" },
  ".cm-line": { padding: "0 16px 0 8px" },
  ".cm-gutters": { backgroundColor: "transparent", border: "none", color: "var(--muted-foreground)" },
  ".cm-lineNumbers .cm-gutterElement": { padding: "0 6px 0 16px", minWidth: "3ch", opacity: "0.55" },
  ".cm-activeLineGutter": { backgroundColor: "transparent" },
  ".cm-activeLineGutter.cm-gutterElement": { opacity: "1", color: "var(--foreground)" },
  ".cm-activeLine": { backgroundColor: "color-mix(in oklch, var(--foreground) 3.5%, transparent)" },
  "&.cm-focused": { outline: "none" },
  "&.cm-focused .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection": {
    backgroundColor: "color-mix(in oklch, var(--flagblue) 24%, transparent) !important",
  },
  ".cm-cursor": { borderLeftColor: "var(--foreground)", borderLeftWidth: "2px" },
  ".cm-matchingBracket": {
    backgroundColor: "color-mix(in oklch, var(--foreground) 10%, transparent)",
    outline: "none",
  },
  ".cm-errorLine": {
    backgroundColor: "var(--crimson-soft)",
    boxShadow: "inset 2px 0 0 var(--crimson)",
  },
});

export const Editor = forwardRef<EditorHandle, Props>(function Editor(
  { initialValue, onChange, onRun, errorLine, className, ariaLabel = "Lipi code editor" },
  ref,
) {
  const host = useRef<HTMLDivElement>(null);
  const view = useRef<EditorView | null>(null);
  const [mounted, setMounted] = useState(false);

  const onChangeRef = useRef(onChange);
  const onRunRef = useRef(onRun);
  onChangeRef.current = onChange;
  onRunRef.current = onRun;

  useEffect(() => {
    const v = new EditorView({
      parent: host.current!,
      state: EditorState.create({
        doc: initialValue,
        extensions: [
          lineNumbers(),
          highlightActiveLineGutter(),
          highlightActiveLine(),
          history(),
          drawSelection(),
          indentOnInput(),
          bracketMatching(),
          closeBrackets(),
          lipiLanguage,
          syntaxHighlighting(lipiHighlight),
          errorLineField,
          theme,
          EditorState.tabSize.of(2),
          Prec.highest(
            keymap.of([
              {
                key: "Mod-Enter",
                run: () => {
                  onRunRef.current?.();
                  return true;
                },
              },
            ]),
          ),
          keymap.of([...closeBracketsKeymap, ...defaultKeymap, ...historyKeymap, indentWithTab]),
          EditorView.updateListener.of((u) => {
            if (u.docChanged) onChangeRef.current?.(u.state.doc.toString());
          }),
          EditorView.contentAttributes.of({
            "aria-label": ariaLabel,
            spellcheck: "false",
            autocorrect: "off",
            autocapitalize: "off",
          }),
        ],
      }),
    });
    view.current = v;
    setMounted(true);
    return () => {
      v.destroy();
      view.current = null;
    };
    // The editor is created once; later changes go through the handle.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    view.current?.dispatch({ effects: setErrorLine.of(errorLine ?? null) });
  }, [errorLine]);

  useImperativeHandle(ref, () => ({
    setValue(value) {
      const v = view.current;
      if (!v) return;
      v.dispatch({
        changes: { from: 0, to: v.state.doc.length, insert: value },
        selection: { anchor: 0 },
        scrollIntoView: true,
      });
    },
    getValue() {
      return view.current?.state.doc.toString() ?? initialValue;
    },
    jumpTo(line) {
      const v = view.current;
      if (!v) return;
      const n = Math.min(Math.max(1, line), v.state.doc.lines);
      const l = v.state.doc.line(n);
      v.dispatch({ selection: { anchor: l.to }, effects: EditorView.scrollIntoView(l.from, { y: "center" }) });
      v.focus();
    },
    focus() {
      view.current?.focus();
    },
  }));

  return (
    <div className={cn("relative min-h-0", className)}>
      {!mounted && (
        // Shown before JavaScript loads so the code is visible immediately.
        <pre className="absolute inset-0 overflow-hidden py-[14px] pr-4 pl-[3.25rem] font-mono text-[14px] leading-[1.7] whitespace-pre text-foreground/85">
          {initialValue}
        </pre>
      )}
      <div ref={host} className="h-full" />
    </div>
  );
});
