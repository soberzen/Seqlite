import { create } from 'zustand';

import { logout as logoutRequest } from '@/service/api/auth';
import { getProfile, UserSchema } from '@/service/api/user';
import type { User } from '@/service/api/user';
import { removeToken } from '@/utils/auth';

type AuthState = {
  loading: boolean;
  isAuthenticated: boolean;
  user: User | null;
  setUser: (user: User) => void;
  refreshUser: () => Promise<void>;
  logout: () => Promise<void>;
};

export const useAuthStore = create<AuthState>((set) => ({
  loading: false,
  isAuthenticated: false,
  user: null,
  setUser: (user: User) => set({ user, isAuthenticated: true }),
  refreshUser: async () => {
    set({ loading: true });
    try {
      const { data } = await getProfile();
      const res = UserSchema.safeParse(data);
      if (!res.success) {
        set({ user: null, isAuthenticated: false });
        return;
      }
      set({ user: res.data, isAuthenticated: true });
    } catch {
      set({ user: null, isAuthenticated: false });
    } finally {
      set({ loading: false });
    }
  },
  logout: async () => {
    try {
      await logoutRequest();
    } finally {
      removeToken();
      set({ user: null, isAuthenticated: false, loading: false });
    }
  },
}));
