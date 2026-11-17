export type EmojiPickerEntry = {
  emoji: string;
  label: string;
  keywords: readonly string[];
};

export type EmojiPickerGroup = {
  id: string;
  label: string;
  emojis: readonly EmojiPickerEntry[];
};

const groupLabels: Readonly<Record<string, string>> = {
  'smileys-emotion': '表情与情感',
  'people-body': '人物与身体',
  'animals-nature': '动物与自然',
  'food-drink': '食物与饮料',
  'travel-places': '旅行与地点',
  activities: '活动',
  objects: '物品',
  symbols: '符号',
  flags: '旗帜',
};

let dataPromise: Promise<readonly EmojiPickerGroup[]> | undefined;

async function readEmojiPickerData(): Promise<readonly EmojiPickerGroup[]> {
  const [{ default: entries }, { default: metadata }] = await Promise.all([
    import('emojibase-data/zh/compact.json'),
    import('emojibase-data/meta/groups.json'),
  ]);
  const orderedEntries = [...entries].sort(
    (first, second) => (first.order ?? 0) - (second.order ?? 0),
  );

  return Object.entries(metadata.groups).flatMap(([groupNumber, id]) => {
    const label = groupLabels[id];
    if (label === undefined) return [];

    const group = Number(groupNumber);
    const emojis = orderedEntries
      .filter((entry) => entry.group === group)
      .map((entry) => ({
        emoji: entry.unicode,
        label: entry.label,
        keywords: entry.tags ?? [],
      }));

    return [{ id, label, emojis }];
  });
}

export function loadEmojiPickerData(): Promise<readonly EmojiPickerGroup[]> {
  dataPromise ??= readEmojiPickerData().catch((error: unknown) => {
    dataPromise = undefined;
    throw error;
  });

  return dataPromise;
}
