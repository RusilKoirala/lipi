---
title: Getting started
description: Install lipi, write a .lip file and run it.
ne: सुरु
section: Start
order: 2
---

## Install

On **macOS or Linux**, paste this into a terminal:

```sh
curl -fsSL https://raw.githubusercontent.com/RusilKoirala/lipi/master/install.sh | sh
```

On **Windows**, download and run [lipi-setup.exe](https://github.com/RusilKoirala/lipi/releases/latest/download/lipi-setup.exe).

Both add `lipi` to your PATH. Open a new terminal and check:

```sh
lipi --version
```

The [install page](/install) has step-by-step help and other options.

## Write a program

Lipi programs live in files ending in `.lip`. Create `namaste.lip`:

```lipi
// Lines starting with // are comments.
rakha naam = "Sita";
rakha umer = 21;

bhana("Namaste,", naam);
bhana("Arko barsa", umer + 1, "barsa pugnuhunchha.");
```

## Run it

```sh
lipi namaste.lip
```

```text
Namaste, Sita
Arko barsa 22 barsa pugnuhunchha.
```

`bhana` prints every value you give it, separated by spaces, and ends the line.

## The rules in short

- Every statement ends with a semicolon `;`.
- A block that ends with `}`, such as a `yadi` or a named `kaam`, does not need one.
- Names can use letters, digits and `_`, but cannot start with a digit: `naam`, `umer2`, `mero_naam`.
- Comments start with `//` and run to the end of the line.

If you break one of these rules, Lipi stops and tells you where. Try deleting a `;` and running again:

```text
Galti: line 2 ma semicolon birsiyo
```
