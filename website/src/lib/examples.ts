export type Example = {
  id: string;
  title: string;
  /** Short Nepali gloss shown under the title. */
  ne: string;
  file: string;
  description: string;
  code: string;
  /** What running it prints (shown on the examples page). */
  output: string;
};

export const EXAMPLES: Example[] = [
  {
    id: "namaste",
    title: "Namaste",
    ne: "नमस्ते",
    file: "namaste.lip",
    description: "Variables, a function and a yadi / natra. The program that started Lipi.",
    code: `rakha naam = "Ram";
bhana("Namaste, " + naam);

kaam jod(a, b) {
  farkau a + b;
}

yadi (jod(2, 3) > 4) {
  bhana("thulo!");
} natra {
  bhana("sano");
}
`,
    output: "Namaste, Ram\nthulo!",
  },
  {
    id: "fibonacci",
    title: "Fibonacci",
    ne: "फिबोनाची",
    file: "fibonacci.lip",
    description: "A function that calls itself. Each number is the sum of the two before it.",
    code: `// fib(n) = fib(n - 1) + fib(n - 2)
kaam fib(n) {
  yadi (n < 2) {
    farkau n;
  }
  farkau fib(n - 1) + fib(n - 2);
}

bhana("fib(10) =", fib(10));
bhana("fib(20) =", fib(20));
`,
    output: "fib(10) = 55\nfib(20) = 6765",
  },
  {
    id: "closures",
    title: "Closures",
    ne: "गुणक",
    file: "gunak.lip",
    description: "A kaam can return another kaam, and the inner one remembers x.",
    code: `kaam gunak(x) {
  farkau kaam(y) {
    farkau x * y;
  };
}

rakha dobar = gunak(2);
rakha tebar = gunak(3);

bhana("dobar(21) =", dobar(21));
bhana("tebar(7)  =", tebar(7));
`,
    output: "dobar(21) = 42\ntebar(7)  = 21",
  },
  {
    id: "grade",
    title: "Exam grade",
    ne: "नतिजा",
    file: "natija.lip",
    description: "Nested yadi / natra to turn marks into a division.",
    code: `kaam natija(ank) {
  yadi (ank > 79) {
    farkau "Distinction";
  } natra {
    yadi (ank > 59) {
      farkau "First division";
    } natra {
      yadi (ank > 39) {
        farkau "Pass";
      } natra {
        farkau "Fail";
      }
    }
  }
}

bhana("Sita:", natija(88));
bhana("Hari:", natija(64));
bhana("Gita:", natija(35));
`,
    output: "Sita: Distinction\nHari: First division\nGita: Fail",
  },
  {
    id: "countdown",
    title: "Countdown",
    ne: "गन्ती",
    file: "ganti.lip",
    description: "Lipi has no loops yet, so recursion does the counting.",
    code: `kaam ganti(n) {
  yadi (n > 0) {
    bhana(n);
    ganti(n - 1);
  } natra {
    bhana("Suru!");
  }
}

ganti(5);
`,
    output: "5\n4\n3\n2\n1\nSuru!",
  },
  {
    id: "higher-order",
    title: "Functions as values",
    ne: "दुई पटक",
    file: "dui-patak.lip",
    description: "Pass a kaam into another kaam, like any other value.",
    code: `kaam duiPatak(f, x) {
  farkau f(f(x));
}

rakha thapDas = kaam(n) { farkau n + 10; };
rakha dobar = kaam(n) { farkau n * 2; };

bhana(duiPatak(thapDas, 5));
bhana(duiPatak(dobar, 5));
`,
    output: "25\n20",
  },
  {
    id: "devanagari",
    title: "Devanagari strings",
    ne: "देवनागरी",
    file: "devanagari.lip",
    description: "Strings can hold Devanagari. lambai counts letters, not bytes.",
    code: `rakha naam = "सीता";
bhana("नमस्ते, " + naam + "!");

rakha shahar = "काठमाडौं";
bhana(shahar, "ma", lambai(shahar), "akshar chhan");
`,
    output: "नमस्ते, सीता!\nकाठमाडौं ma 8 akshar chhan",
  },
  {
    id: "galti",
    title: "Galti demo",
    ne: "गल्ती",
    file: "galti.lip",
    description: "A program with a mistake on purpose. Run it, read the Galti, then fix it.",
    code: `rakha naam = "Sita";
rakha umer = 21;
bhana("Namaste, " + naam)

// Line 3 is missing a ';'. Add it and run again.
// Then try this line:  bhana(naam + umer);
`,
    output: "Galti: line 3 ma semicolon birsiyo",
  },
];

export function getExample(id: string | null | undefined): Example | undefined {
  return EXAMPLES.find((e) => e.id === id);
}
