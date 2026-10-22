---
title: yadi and natra
description: Run code only when a condition holds.
ne: यदि / नत्र
section: Language
order: 4
---

`yadi` means "if" and `natra` means "otherwise".

```lipi
rakha umer = 20;

yadi (umer > 17) {
  bhana("Voting garna milchha.");
} natra {
  bhana("Ajhai sano.");
}
```

```text
Voting garna milchha.
```

The condition goes in parentheses, and each block goes in braces. `natra` is optional.

## More than two choices

Put another `yadi` inside the `natra` block:

```lipi
rakha ank = 72;

yadi (ank > 79) {
  bhana("Distinction");
} natra {
  yadi (ank > 59) {
    bhana("First division");
  } natra {
    bhana("Pass");
  }
}
```

```text
First division
```

## yadi gives back a value

A `yadi` is an expression: it produces the value of the last line in the block that ran. You can store that value.

```lipi
rakha a = 8;
rakha b = 13;
rakha thulo = yadi (a > b) { a; } natra { b; };

bhana("Thulo:", thulo);
```

```text
Thulo: 13
```

If the condition is false and there is no `natra`, the result is `khali`.

## No semicolon after the brace

A `yadi` statement ends with `}`, so it does not need a `;` after it. Adding one is allowed.
