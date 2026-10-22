---
title: Functions
description: Write reusable pieces of work with kaam and farkau.
ne: काम
section: Language
order: 5
---

`kaam` means "work". A function is a piece of work you name once and call as often as you like. `farkau` ("send back") returns a value from it.

```lipi
kaam jod(a, b) {
  farkau a + b;
}

bhana(jod(2, 3));
bhana(jod(40, 2));
```

```text
5
42
```

## Parameters and calls

Inside the parentheses are the parameters. When you call the function, each argument is given to the matching parameter. The number of arguments must match:

```lipi
kaam jod(a, b) {
  farkau a + b;
}

jod(1);
```

```text
Galti: line 5 ma 'jod' lai 2 ota input chahinchha, 1 ota aayo
```

## Functions without a name

`kaam` without a name makes a function value. Store it with `rakha`:

```lipi
rakha dobar = kaam(x) {
  farkau x * 2;
};

bhana(dobar(21));
```

```text
42
```

`kaam jod(a, b) { … }` is a short way of writing `rakha jod = kaam(a, b) { … };`.

## Functions that call themselves

A function can call itself. This is called recursion, and it is how Lipi repeats things while it has no loops:

```lipi
kaam factorial(n) {
  yadi (n < 2) {
    farkau 1;
  }
  farkau n * factorial(n - 1);
}

bhana(factorial(10));
```

```text
3628800
```

If a function calls itself forever, Lipi stops it instead of crashing:

```text
Galti: line 2 ma 'f' le aafailai dherai patak bolayo (recursion rokiyena)
```

## Closures

A function remembers the variables around it when it was made. Here the inner `kaam` remembers `x`:

```lipi
kaam gunak(x) {
  farkau kaam(y) {
    farkau x * y;
  };
}

rakha tebar = gunak(3);
bhana(tebar(7));
```

```text
21
```

## Functions as arguments

Because functions are values, you can pass one into another:

```lipi
kaam duiPatak(f, x) {
  farkau f(f(x));
}

bhana(duiPatak(kaam(n) { farkau n + 10; }, 5));
```

```text
25
```
