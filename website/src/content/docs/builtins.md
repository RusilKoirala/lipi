---
title: Built-in functions
description: bhana and lambai, the two functions Lipi gives you for free.
ne: भन
section: Language
order: 6
---

Built-in functions are ready to use without writing them yourself.

## bhana

`bhana` ("say") prints its arguments on one line, separated by spaces.

```lipi
bhana("Namaste");
bhana("Jod:", 2 + 3);
bhana(sahi, galat, 10);
```

```text
Namaste
Jod: 5
sahi galat 10
```

It gives back `khali`.

## lambai

`lambai` ("length") counts the letters in a string. It counts Devanagari letters and marks, not bytes.

```lipi
bhana(lambai("Pokhara"));
bhana(lambai("पोखरा"));
```

```text
7
5
```

It takes exactly one string. Anything else is a Galti:

```lipi
lambai(42);
```

```text
Galti: line 1 ma lambai le sankhya napna sakdaina
```
