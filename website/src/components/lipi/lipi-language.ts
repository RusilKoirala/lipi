// CodeMirror support for Lipi: a small tokenizer + colours from the --syn-* CSS variables.
import { HighlightStyle, StreamLanguage } from "@codemirror/language";
import { tags as t } from "@lezer/highlight";

const KEYWORDS = new Set(["rakha", "kaam", "farkau", "yadi", "natra"]);
const ATOMS = new Set(["sahi", "galat"]);
const BUILTINS = new Set(["bhana", "lambai"]);

export const lipiLanguage = StreamLanguage.define<{ inString: boolean }>({
  name: "lipi",
  startState: () => ({ inString: false }),
  token(stream, state) {
    if (state.inString) {
      while (!stream.eol()) {
        if (stream.next() === '"') {
          state.inString = false;
          break;
        }
      }
      return "string";
    }

    if (stream.eatSpace()) return null;

    if (stream.match("//")) {
      stream.skipToEnd();
      return "comment";
    }

    if (stream.peek() === '"') {
      stream.next();
      state.inString = true;
      while (!stream.eol()) {
        if (stream.next() === '"') {
          state.inString = false;
          break;
        }
      }
      return "string";
    }

    if (stream.match(/^\d+/)) return "number";

    if (stream.match(/^[A-Za-z_][A-Za-z0-9_]*/)) {
      const word = stream.current();
      if (KEYWORDS.has(word)) return "keyword";
      if (ATOMS.has(word)) return "atom";
      if (BUILTINS.has(word)) return "builtin";
      return "variableName";
    }

    if (stream.match(/^(==|!=|[-+*/<>=!])/)) return "operator";

    stream.next();
    return null;
  },
  languageData: {
    commentTokens: { line: "//" },
    closeBrackets: { brackets: ["(", "{", '"'] },
  },
});

export const lipiHighlight = HighlightStyle.define([
  { tag: t.keyword, color: "var(--syn-keyword)", fontWeight: "500" },
  { tag: t.string, color: "var(--syn-string)" },
  { tag: [t.number, t.atom, t.bool], color: "var(--syn-constant)" },
  { tag: t.standard(t.variableName), color: "var(--syn-builtin)", fontWeight: "500" },
  { tag: t.comment, color: "var(--syn-comment)", fontStyle: "italic" },
  { tag: t.operator, color: "var(--syn-operator)" },
]);
