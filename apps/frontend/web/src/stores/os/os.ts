import { create } from 'zustand';

const isMac =
  typeof navigator !== 'undefined' &&
  navigator.userAgent.toLowerCase().includes('mac');

export const useIsMacStore = create(() => ({
  isMac,
}));

export const useIsMac = () => useIsMacStore((state) => state.isMac);
