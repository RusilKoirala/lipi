---
title: Operators
description: Arithmetic, comparison and how Lipi decides what to do first.
ne: चिह्न
section: Language
order: 3
---

## Arithmetic

| Operator | Meaning | Example | Result |
| --- | --- | --- | --- |
| `+` | add (or join strings) | `2 + 3` | `5` |
| `-` | subtract | `10 - 4` | `6` |
| `*` | multiply | `6 * 7` | `42` |
| `/` | divide (whole numbers) | `9 / 2` | `4` |
| `-x` | negative | `-5` | `-5` |

Dividing by zero is a Galti: `0 le bhag garna mildaina`.

## Comparison

| Operator | Meaning |
| --- | --- |
| `==` | equal |
| `!=` | not equal |
| `<` | less than |
| `>` | greater than |
| `!x` | not |

There is no `<=` or `>=` yet. Write `umer > 17` instead of `umer >= 18`.

## What happens first

Multiplication and division happen before addition and subtraction, and comparisons happen last. Use parentheses when you want something else.

```lipi
bhana(1 + 2 * 3);
bhana((1 + 2) * 3);
bhana(2 + 3 > 4);
```

```text
7
9
sahi
```

From strongest to weakest:

1. Calls: `jod(1, 2)`
2. Prefix: `-x`, `!x`
3. `*` `/`
4. `+` `-`
5. `<` `>` `==` `!=`

Comparisons all have the same strength, so `a < b == sahi` is read left to right as `(a < b) == sahi`.

The parser that decides this is a Pratt parser. [How Lipi works](/docs/how-it-works) explains it.
