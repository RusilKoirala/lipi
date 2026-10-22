// @ts-check
import { defineConfig } from "astro/config";

import react from "@astrojs/react";
import tailwindcss from "@tailwindcss/vite";

import { lipiDark, lipiGrammar, lipiLight } from "./src/lib/lipi-grammar.ts";

// https://astro.build/config
export default defineConfig({
  integrations: [react()],

  // Keep spaces between inline elements written on separate lines.
  compressHTML: false,

  redirects: {
    "/docs": "/docs/introduction",
  },

  markdown: {
    shikiConfig: {
      langs: [structuredClone(lipiGrammar)],
      themes: { light: lipiLight, dark: lipiDark },
      defaultColor: false,
    },
  },

  vite: {
    plugins: [tailwindcss()],
  },
});
