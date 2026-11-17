const storageKey = 'qingmo:emoji-picker:recents';
const recentLimit = 8;
let unsavedRecents: string[] = [];

export function readRecentEmojis(): string[] {
  if (typeof window === 'undefined') return [];

  try {
    const stored: unknown = JSON.parse(
      window.localStorage.getItem(storageKey) ?? '[]',
    );
    const emojis = Array.isArray(stored)
      ? stored.filter(
          (value): value is string =>
            typeof value === 'string' && value.trim().length > 0,
        )
      : [];
    return [...new Set([...unsavedRecents, ...emojis])].slice(0, recentLimit);
  } catch {
    return [...unsavedRecents];
  }
}

export function rememberRecentEmoji(emoji: string): string[] {
  const recents = [...new Set([emoji, ...readRecentEmojis()])].slice(
    0,
    recentLimit,
  );
  if (typeof window === 'undefined') return recents;

  try {
    window.localStorage.setItem(storageKey, JSON.stringify(recents));
    unsavedRecents = [];
  } catch {
    unsavedRecents = [...recents];
  }

  return recents;
}
