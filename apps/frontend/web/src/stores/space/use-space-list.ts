import { create } from 'zustand';

import { listSpace, type Space } from '@/service/api/space';

export type SpaceListItem = Space & {
  updatedAtLabel?: string;
};

type SpaceListState = {
  spaces: SpaceListItem[];
  fetchSpaceList: () => Promise<void>;
  status: 'loading' | 'success' | 'error';
};

export const useSpaceListStore = create<SpaceListState>((set) => ({
  spaces: [
    {
      id: '12213131231',
      name: '产品设计原则',
      emoji: '☺️',
      description: '共同遵循产品设计原则，确保产品符合用户需求',
      updatedAtLabel: '今天 10:24',
    },
    {
      id: '12213131231',
      name: '产品设计原则',
      emoji: '☺️',
      description: '共同遵循产品设计原则，确保产品符合用户需求',
      updatedAtLabel: '今天 10:24',
    },
    {
      id: '12213131231',
      name: '产品设计原则',
      emoji: '☺️',
      description: '共同遵循产品设计原则，确保产品符合用户需求',
      updatedAtLabel: '今天 10:24',
    },
    {
      id: '12213131231',
      name: '产品设计原则',
      emoji: '☺️',
      description: '共同遵循产品设计原则，确保产品符合用户需求',
      updatedAtLabel: '今天 10:24',
    },
    {
      id: '12213131231',
      name: '产品设计原则',
      emoji: '☺️',
      description: '共同遵循产品设计原则，确保产品符合用户需求',
      updatedAtLabel: '今天 10:24',
    },
    {
      id: '12213131231',
      name: '产品设计原则',
      emoji: '☺️',
      description: '共同遵循产品设计原则，确保产品符合用户需求',
      updatedAtLabel: '今天 10:24',
    },
  ],
  fetchSpaceList: async () => {
    set({ status: 'loading' });
    try {
      const spaces = await listSpace();
      set({ spaces: spaces, status: 'success' });
    } catch (error) {
      set({ status: 'error' });
      throw error;
    }
  },
  status: 'success',
}));
