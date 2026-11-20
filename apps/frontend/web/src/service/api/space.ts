import { z } from 'zod';

export const SpaceSchema = z.object({
  id: z.uuid(),
  name: z.string(),
  emoji: z.preprocess(
    (value) => (value === '' ? null : value),
    z.emoji().nullable(),
  ),
});

export type Space = z.infer<typeof SpaceSchema>;
