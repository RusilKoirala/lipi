---
title: Variables
description: Store values with rakha and use them by name.
ne: राख
section: Language
order: 1
---

`rakha` means "keep". It stores a value under a name so you can use it later.

```lipi
rakha naam = "Hari";
rakha umer = 30;
rakha bihe_bhayo = galat;

bhana(naam, umer, bihe_bhayo);
```

```text
Hari 30 galat
```

## Using a variable

Once a name has a value, write the name anywhere you would write the value:

```lipi
rakha mulya = 250;
rakha sankhya = 4;
rakha jamma = mulya * sankhya;

bhana("Jamma:", jamma);
```

```text
Jamma: 1000
```

## Giving a name a new value

Write `rakha` again with the same name. The new value replaces the old one from that point on.

```lipi
rakha score = 10;
rakha score = score + 5;
bhana(score);
```

```text
15
```

## Names Lipi doesn't know

Using a name before it has a value is a Galti:

```lipi
bhana(shahar);
```

```text
Galti: line 1 ma 'shahar' bhanne naam bhetiyena
```

"No name called 'shahar' was found." Check the spelling, or add a `rakha` line above it.
