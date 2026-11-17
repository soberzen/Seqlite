import { z } from 'zod';

import { post } from '../base';
import type { AccessTokenResponse, HTTPResponse } from '../fetch';
import type { User } from './user';

export const LoginSchema = z.object({
  email: z.email('请输入有效的邮箱地址'),
  password: z.string().min(1, '密码不能为空'),
});

export const RegisterSchema = z.object({
  email: z.email('请输入有效的邮箱地址'),
  password: z.string().min(6, '密码长度至少为6位'),
  userName: z.string().min(1, '用户名不能为空'),
});

export type LoginDto = z.infer<typeof LoginSchema>;
export type RegisterDto = z.infer<typeof RegisterSchema>;

export function login(data: LoginDto) {
  return post<AccessTokenResponse>('/auth/login', {
    body: data,
  });
}

export function register(data: RegisterDto) {
  return post<HTTPResponse<User>>('/auth/register', { body: data });
}

export function refresh() {
  return post<AccessTokenResponse>('/auth/refresh');
}

export function logout() {
  return post<HTTPResponse<null>>('/auth/logout');
}
