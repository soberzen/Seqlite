import { z } from 'zod';

import { get } from '../base';
import type { HTTPResponse } from '../fetch';

export const UserSchema = z.object({
  userName: z.string(),
  email: z.email(),
  avatarURL: z.preprocess(
    (value) => (value === '' ? null : value),
    z.url().nullable(),
  ),
});

export type User = z.infer<typeof UserSchema>;
export function getProfile() {
  return get<HTTPResponse<User>>('/user/me');
}
