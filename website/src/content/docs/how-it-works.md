---
title: How Lipi works
description: The lexer, the Pratt parser and the tree-walking evaluator, in plain words.
ne: भित्र
section: More
order: 2
---

Lipi is a tree-walking interpreter written in Go. It follows [Writing an Interpreter in Go](https://interpreterbook.com) by Thorsten Ball, which builds a language called Monkey. Lipi swaps in Nepali keywords and error messages, and tracks line numbers on every token so each Galti can point at a line.

The code is split into small packages, one for each step:

```text
lipi/
├── token/      the kinds of tokens and the Nepali keyword table
├── lexer/      text → tokens
├── ast/        the shapes of the syntax tree
├── parser/     tokens → tree (Pratt parser)
├── object/     runtime values and environments
├── evaluator/  tree → result, plus bhana and lambai
├── wasm/       the bridge used by the online playground
└── main.go     the lipi command
```

## 1. Lexer

The lexer reads the file one character at a time and groups characters into tokens. `rakha naam = "Ram";` becomes:

```text
RAKHA "rakha"   IDENT "naam"   ASSIGN "="   STRING "Ram"   SEMICOLON ";"
```

Each token also records its line and column. Words are checked against a small keyword table: `rakha`, `kaam`, `farkau`, `yadi`, `natra`, `sahi` and `galat` are keywords, and everything else is a name.

## 2. Parser

The parser turns the flat list of tokens into a tree. For `1 + 2 * 3`:

```text
      +
     / \
    1   *
       / \
      2   3
```

Lipi uses a **Pratt parser**. Each operator has a strength (`*` is stronger than `+`), and one loop keeps attaching operators as long as the next one is stronger than the current one. A function call is just another operator: `(` after a value means "call it", with the highest strength of all.

`yadi` and `kaam` are expressions in Lipi, so they are parsed the same way as `5` or `-x`, from the token that starts them.

## 3. Evaluator

The evaluator walks the tree and computes a value for each node. A `+` node evaluates its left side, then its right side, then adds them.

Variables are kept in an **environment**: a table of names and values with a link to the environment around it. Calling a function creates a new environment that links back to the one where the function was made. That link is what lets a closure remember its variables after the outer function has returned.

`farkau` wraps its value so it can pass up through nested blocks until it reaches the function call, where it is unwrapped. A Galti travels up the same way and stops the program.

## 4. Command line and browser

`main.go` reads a `.lip` file, runs it, and prints any Galti. The `wasm/` package compiles the same interpreter to WebAssembly and exposes one function, `runLipi(code)`, to JavaScript. The playground on this site calls it, so the browser runs exactly the same code as your terminal.
