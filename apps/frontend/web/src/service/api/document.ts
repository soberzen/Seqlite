import { z } from 'zod';

import { SpaceSchema } from './space';

export const SpaceDocumentPageSchema = z.object({
  items: z.array(
    z.object({
      id: z.uuid,
      title: z.string,
      emoji: SpaceSchema.shape.emoji,
    }),
  ),
  page: z.number().int().positive(),
  pageSize: z.number().int().positive(),
  hasMore: z.boolean(),
});

export type SpaceDocumentPage = z.infer<typeof SpaceDocumentPageSchema>;
