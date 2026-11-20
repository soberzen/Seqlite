import { z } from 'zod';

import { get, post } from '../base';

export const SpaceSchema = z.object({
  id: z.uuid(),
  name: z.string(),
  emoji: z.preprocess(
    (value) => (value === '' ? null : value),
    z.emoji().nullable(),
  ),
  description: z.string(),
});

export const spaceNameLimit = 80;
export const spaceDescriptionLimit = 200;
export const CreateSpaceSchema = z.object({
  name: z
    .string()
    .trim()
    .min(1, '请输入空间名称')
    .max(spaceNameLimit, `空间名称不能超过 ${spaceNameLimit} 个字符`),
  description: z
    .string()
    .max(spaceDescriptionLimit, `描述不能超过 ${spaceDescriptionLimit} 个字符`),
  emoji: z.emoji().nullable(),
});

export type Space = z.infer<typeof SpaceSchema>;
export type CreateSpaceDto = z.infer<typeof CreateSpaceSchema>;

export function createSpace(data: CreateSpaceDto) {
  return post<Space>('/api/space', { body: data });
}

export function listSpace() {
  return get<Space[]>('/api/space');
}
