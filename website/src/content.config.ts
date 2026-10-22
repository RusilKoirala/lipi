import { defineCollection } from "astro:content";
import { glob } from "astro/loaders";
import { z } from "astro/zod";

const docs = defineCollection({
  loader: glob({ pattern: "**/*.md", base: "./src/content/docs" }),
  schema: z.object({
    title: z.string(),
    description: z.string(),
    /** Nepali word shown next to the title. */
    ne: z.string().optional(),
    section: z.enum(["Start", "Language", "More"]),
    order: z.number(),
  }),
});

export const collections = { docs };
