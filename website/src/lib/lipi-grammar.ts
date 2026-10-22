// Shiki (TextMate) grammar + themes for highlighting Lipi on static pages and in docs.
// The editor uses the same colours through the --syn-* CSS variables in global.css.

export const lipiGrammar = {
  name: "lipi",
  scopeName: "source.lipi",
  aliases: ["lip"],
  patterns: [
    { include: "#comment" },
    { include: "#string" },
    { include: "#keyword" },
    { include: "#constant" },
    { include: "#builtin" },
    { include: "#number" },
    { include: "#operator" },
  ],
  repository: {
    comment: { name: "comment.line.double-slash.lipi", match: "//.*$" },
    string: { name: "string.quoted.double.lipi", begin: '"', end: '"' },
    keyword: { name: "keyword.control.lipi", match: "\\b(rakha|kaam|farkau|yadi|natra)\\b" },
    constant: { name: "constant.language.lipi", match: "\\b(sahi|galat)\\b" },
    builtin: { name: "support.function.builtin.lipi", match: "\\b(bhana|lambai)\\b" },
    number: { name: "constant.numeric.lipi", match: "\\b\\d+\\b" },
    operator: { name: "keyword.operator.lipi", match: "==|!=|[-+*/<>=!]" },
  },
};

type Palette = {
  fg: string;
  comment: string;
  keyword: string;
  string: string;
  constant: string;
  builtin: string;
  operator: string;
};

function theme(name: string, type: "light" | "dark", c: Palette) {
  return {
    name,
    type,
    colors: { "editor.background": "#00000000", "editor.foreground": c.fg },
    tokenColors: [
      { scope: ["comment"], settings: { foreground: c.comment, fontStyle: "italic" } },
      { scope: ["keyword.control"], settings: { foreground: c.keyword } },
      { scope: ["string"], settings: { foreground: c.string } },
      { scope: ["constant.numeric", "constant.language"], settings: { foreground: c.constant } },
      { scope: ["support.function"], settings: { foreground: c.builtin } },
      { scope: ["keyword.operator"], settings: { foreground: c.operator } },
    ],
  };
}

export const lipiLight = theme("lipi-light", "light", {
  fg: "#1c2030",
  comment: "#8b90a0",
  keyword: "#c8102e",
  string: "#1f4ea3",
  constant: "#7a4a12",
  builtin: "#1c2030",
  operator: "#6b7080",
});

export const lipiDark = theme("lipi-dark", "dark", {
  fg: "#e7e9f0",
  comment: "#6f7487",
  keyword: "#ff5d72",
  string: "#93b6ff",
  constant: "#e3b37a",
  builtin: "#e7e9f0",
  operator: "#a0a5b5",
});
