---
title: Introduction
description: What Lipi is, who it is for, and what a Lipi program looks like.
ne: परिचय
section: Start
order: 1
---

Lipi (लिपि, "script") is a small programming language whose keywords are Nepali words written in Roman letters. Instead of `let`, `if` and `return`, you write `rakha`, `yadi` and `farkau`. When something goes wrong, the error is in Nepali too, and it always says which line to look at.

```lipi
rakha naam = "Ram";
bhana("Namaste, " + naam);

kaam jod(a, b) {
  farkau a + b;
}

yadi (jod(2, 3) > 4) {
  bhana("thulo!");
} natra {
  bhana("sano");
}
```

Running this prints:

```text
Namaste, Ram
thulo!
```

## Who it is for

- **Nepali speakers learning to code.** The ideas (variables, conditions, functions) are the same as in any language, but the words are ones you already know.
- **Anyone curious about interpreters.** Lipi is small enough to read in an afternoon. It follows [Writing an Interpreter in Go](https://interpreterbook.com) closely, so you can read the book and the code side by side.

## What it has

| Feature | Lipi | Meaning |
| --- | --- | --- |
| Variables | `rakha x = 5;` | let |
| Functions | `kaam jod(a, b) { … }` | function |
| Return | `farkau x;` | return |
| Conditions | `yadi (…) { … } natra { … }` | if / else |
| Booleans | `sahi`, `galat` | true, false |
| Printing | `bhana(…)` | print |
| Length | `lambai("…")` | len |

Numbers are whole numbers (integers). Strings go in double quotes and can contain Devanagari. Functions are values: you can store them, pass them around and return them from other functions.

## Two ways to run it

1. **In your browser**, in the [playground](/playground). Nothing to install.
2. **On your computer**, with the `lipi` command. See [Getting started](/docs/getting-started).
