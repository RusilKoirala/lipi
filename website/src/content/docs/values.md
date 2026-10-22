---
title: Values and types
description: Numbers, strings, sahi and galat, khali, and functions.
ne: मान
section: Language
order: 2
---

Every value in Lipi has a type. Errors use these Nepali type names, so it helps to know them.

| Type | Written as | Example |
| --- | --- | --- |
| `sankhya` (number) | digits | `42` |
| `string` (text) | double quotes | `"Kathmandu"` |
| `sahi/galat` (boolean) | `sahi` or `galat` | `sahi` |
| `khali` (nothing) | — | what a `yadi` with no `natra` gives back when it doesn't run |
| `kaam` (function) | `kaam(…) { … }` | `kaam(x) { farkau x * 2; }` |

## Numbers

Numbers are whole numbers. Division drops the remainder.

```lipi
bhana(7 / 2);
bhana(-7 + 10);
```

```text
3
3
```

## Strings

Strings are text in double quotes. Join them with `+`. They can hold Devanagari:

```lipi
rakha pahilo = "नमस्ते";
bhana(pahilo + ", Lipi!");
bhana(lambai(pahilo));
```

```text
नमस्ते, Lipi!
6
```

A number and a string can't be added. Lipi will not guess what you meant:

```lipi
bhana("Umer: " + 21);
```

```text
Galti: line 1 ma string ra sankhya jodna mildaina
```

Pass them to `bhana` separately instead: `bhana("Umer:", 21);`.

## sahi and galat

Comparisons give back `sahi` (true) or `galat` (false):

```lipi
bhana(5 > 3);
bhana("ram" == "shyam");
bhana(!sahi);
```

```text
sahi
galat
galat
```

When Lipi checks a condition, only `galat` and `khali` count as false. Every other value, including `0` and `""`, counts as true.
