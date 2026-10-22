---
title: Galti reference
description: Every error message Lipi can show, what it means and how to fix it.
ne: गल्ती
section: More
order: 1
---

Every Lipi error starts with `Galti:` and the line number, then says what went wrong. Lipi stops at the first Galti it finds.

```text
Galti: line 3 ma semicolon birsiyo
```

## While reading your code

These come from the parser, before anything runs.

| Message | Meaning | Fix |
| --- | --- | --- |
| `semicolon birsiyo` | You forgot a semicolon. | Add `;` at the end of that line. |
| `'=' chahiyo` | Lipi expected `=` here (after `rakha naam`). | Write `rakha naam = …;` |
| `'IDENT' chahiyo` | Lipi expected a name here. | Put a name after `rakha` or in a parameter list. |
| `')' chahiyo` / `'{' chahiyo` | A bracket is missing. | Check that every `(` and `{` has its pair. |
| `'{' kholiyo tara '}' le band gariyena` | A block was opened but never closed. | Add the missing `}`. |
| `'…' bujhina` | Lipi did not understand this symbol here. | Check for a typo or a stray character. |

## While running your code

These come from the evaluator, while your program runs.

| Message | Meaning |
| --- | --- |
| `'x' bhanne naam bhetiyena` | No name called `x` was found. |
| `sankhya ra string jodna mildaina` | A number and a string can't be added. |
| `sahi/galat lai '-' garna mildaina` | `-` only works on numbers. |
| `string ma '*' chaldaina` | Strings only support `+`, `==` and `!=`. |
| `0 le bhag garna mildaina` | You can't divide by zero. |
| `'jod' lai 2 ota input chahinchha, 1 ota aayo` | The function needs 2 arguments, but got 1. |
| `'x' kaam hoina, bolauna mildaina` | `x` is not a function, so it can't be called. |
| `lambai le sankhya napna sakdaina` | `lambai` can only measure strings. |
| `… le aafailai dherai patak bolayo` | A function called itself too many times. |

## A worked example

```lipi
rakha mulya = 100;
rakha kar = 13;
bhana("Jamma: " + mulya + kar);
```

```text
Galti: line 3 ma string ra sankhya jodna mildaina
```

`"Jamma: " + mulya` tries to add a string and a number. Do the maths first, then give both values to `bhana`:

```lipi
rakha mulya = 100;
rakha kar = 13;
bhana("Jamma:", mulya + kar);
```

```text
Jamma: 113
```
