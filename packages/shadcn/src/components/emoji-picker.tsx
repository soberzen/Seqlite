import {
  CarFront,
  Eraser,
  Flag,
  Hand,
  Heart,
  Lightbulb,
  PawPrint,
  Sandwich,
  Smile,
  Volleyball,
  type LucideIcon,
} from 'lucide-react';
import {
  useId,
  useRef,
  useState,
  type KeyboardEvent,
  type ReactElement,
} from 'react';

import {
  loadEmojiPickerData,
  type EmojiPickerEntry,
  type EmojiPickerGroup,
} from '../lib/emoji-picker-data';
import {
  readRecentEmojis,
  rememberRecentEmoji,
} from '../lib/emoji-picker-recents';
import { cn } from '../lib/utils';
import { Button } from './button';
import { Input } from './input';
import { Popover, PopoverContent, PopoverTrigger } from './popover';

const categoryIcons: Record<string, LucideIcon> = {
  'smileys-emotion': Smile,
  'people-body': Hand,
  'animals-nature': PawPrint,
  'food-drink': Sandwich,
  'travel-places': CarFront,
  activities: Volleyball,
  objects: Lightbulb,
  symbols: Heart,
  flags: Flag,
};

// Keep the keyboard's row movement in sync with the visual grid.
const columns = 8;
const arrowOffsets: Readonly<Record<string, number>> = {
  ArrowLeft: -1,
  ArrowRight: 1,
  ArrowUp: -columns,
  ArrowDown: columns,
};

type DataState =
  | { status: 'idle' | 'loading' | 'error' }
  | { status: 'ready'; groups: readonly EmojiPickerGroup[] };

export type EmojiPickerProps = {
  children: ReactElement;
  modal?: boolean;
  onSelect: (emoji: string) => void;
  onClear?: () => void;
  value?: string | null;
};

function EmojiGrid({
  group,
  value,
  onSelect,
}: {
  group: EmojiPickerGroup;
  value: EmojiPickerProps['value'];
  onSelect: (emoji: string) => void;
}) {
  const headingId = useId();
  const [focusedEmoji, setFocusedEmoji] = useState<string>();
  const tabStop = group.emojis.some(({ emoji }) => emoji === focusedEmoji)
    ? focusedEmoji
    : group.emojis[0]?.emoji;

  function handleGridKeyDown(event: KeyboardEvent<HTMLUListElement>) {
    const offset = arrowOffsets[event.key];
    if (offset === undefined && event.key !== 'Home' && event.key !== 'End') {
      return;
    }

    const buttons = Array.from(event.currentTarget.querySelectorAll('button'));
    const index = buttons.findIndex((button) => button === event.target);
    if (index < 0) return;

    event.preventDefault();
    if (event.key === 'Home') buttons[0]?.focus();
    else if (event.key === 'End') buttons.at(-1)?.focus();
    else if (offset !== undefined) buttons[index + offset]?.focus();
  }

  return (
    <section aria-labelledby={headingId}>
      <h3
        id={headingId}
        className='sticky top-0 z-10 bg-popover pt-2 pb-1.5 text-[0.8125rem] font-medium text-muted-foreground'
      >
        {group.label}
      </h3>
      <ul
        className='grid grid-cols-8 gap-0.5'
        onKeyDown={handleGridKeyDown}
      >
        {group.emojis.map(({ emoji, label }) => (
          <li key={emoji}>
            <button
              type='button'
              aria-label={label}
              aria-pressed={value === emoji}
              title={label}
              tabIndex={emoji === tabStop ? 0 : -1}
              data-emoji-option=''
              className='flex aspect-square w-full scroll-mt-7 items-center justify-center rounded-lg text-2xl leading-none select-none font-sans hover:bg-muted focus-visible:bg-muted focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none motion-safe:transition-colors aria-pressed:bg-accent'
              onFocus={() => setFocusedEmoji(emoji)}
              onClick={() => onSelect(emoji)}
            >
              {emoji}
            </button>
          </li>
        ))}
      </ul>
    </section>
  );
}

function matchesQuery(
  { emoji, label, keywords }: EmojiPickerEntry,
  query: string,
) {
  const text = [emoji, label, ...keywords]
    .join(' ')
    .normalize('NFKC')
    .toLowerCase();
  return text.includes(query);
}

export function EmojiPicker({
  children,
  modal = false,
  onSelect,
  onClear,
  value,
}: EmojiPickerProps) {
  const [open, setOpen] = useState(false);
  const [data, setData] = useState<DataState>({ status: 'idle' });
  const [query, setQuery] = useState('');
  const [category, setCategory] = useState<string>();
  const [recentEmojis, setRecentEmojis] = useState<string[]>([]);
  const resultsRef = useRef<HTMLDivElement>(null);

  const groups = data.status === 'ready' ? data.groups : [];
  const activeGroup = groups.find(({ id }) => id === category) ?? groups[0];
  const normalizedQuery = query.trim().normalize('NFKC').toLowerCase();
  const visibleGroups: EmojiPickerGroup[] = normalizedQuery
    ? groups
        .map((group) => ({
          ...group,
          emojis: group.emojis.filter((entry) =>
            matchesQuery(entry, normalizedQuery),
          ),
        }))
        .filter((group) => group.emojis.length > 0)
    : activeGroup
      ? [activeGroup]
      : [];

  if (!normalizedQuery && recentEmojis.length > 0) {
    const entries = groups.flatMap(({ emojis }) => emojis);
    const emojis = recentEmojis.flatMap((emoji) => {
      const entry = entries.find((option) => option.emoji === emoji);
      return entry ? [entry] : [];
    });

    if (emojis.length > 0) {
      visibleGroups.unshift({ id: 'recent', label: '最近使用', emojis });
    }
  }

  function loadData() {
    setData({ status: 'loading' });
    void loadEmojiPickerData().then(
      (groups) => setData({ status: 'ready', groups }),
      () => setData({ status: 'error' }),
    );
  }

  function handleOpenChange(nextOpen: boolean) {
    setOpen(nextOpen);
    if (nextOpen) {
      setRecentEmojis(readRecentEmojis());
      if (data.status === 'idle' || data.status === 'error') loadData();
    } else {
      setQuery('');
    }
  }

  function updateQuery(nextQuery: string) {
    setQuery(nextQuery);
    if (resultsRef.current) resultsRef.current.scrollTop = 0;
  }

  function selectEmoji(emoji: string) {
    rememberRecentEmoji(emoji);
    onSelect(emoji);
    handleOpenChange(false);
  }

  function handleSearchKeyDown(event: KeyboardEvent<HTMLInputElement>) {
    if (event.nativeEvent.isComposing) return;
    if (event.key === 'ArrowDown') {
      const firstOption = resultsRef.current?.querySelector<HTMLButtonElement>(
        '[data-emoji-option]',
      );
      if (firstOption) {
        event.preventDefault();
        firstOption.focus();
      }
    } else if (event.key === 'Enter' && normalizedQuery) {
      const firstMatch = visibleGroups[0]?.emojis[0];
      if (firstMatch) {
        event.preventDefault();
        selectEmoji(firstMatch.emoji);
      }
    }
  }

  return (
    <Popover
      modal={modal}
      open={open}
      onOpenChange={handleOpenChange}
    >
      <PopoverTrigger asChild>{children}</PopoverTrigger>
      <PopoverContent
        aria-label='选择表情'
        collisionPadding={12}
        sideOffset={8}
        className='max-h-[min(26rem,var(--radix-popover-content-available-height))] w-[23rem] max-w-(--radix-popover-content-available-width) gap-0 overflow-hidden rounded-2xl p-0 shadow-lg motion-reduce:animate-none'
      >
        <div className='shrink-0 px-3 pt-3'>
          <Input
            type='search'
            aria-label='搜索表情'
            placeholder='搜索表情…'
            value={query}
            onChange={(event) => updateQuery(event.target.value)}
            onKeyDown={handleSearchKeyDown}
            className='h-9 px-3'
          />
          {groups.length > 0 && (
            <div
              role='group'
              aria-label='表情分类'
              className='mt-2 grid grid-cols-9 gap-0.5 pb-1'
            >
              {groups.map((group) => {
                const Icon = categoryIcons[group.id] ?? Smile;
                const active = !normalizedQuery && activeGroup?.id === group.id;

                return (
                  <button
                    key={group.id}
                    type='button'
                    aria-label={group.label}
                    aria-pressed={active}
                    title={group.label}
                    className={cn(
                      'mx-auto grid aspect-square w-full max-w-8 place-items-center rounded-full text-muted-foreground outline-none hover:bg-muted hover:text-foreground focus-visible:ring-2 focus-visible:ring-ring',
                      active && 'bg-muted text-foreground',
                    )}
                    onClick={() => {
                      setCategory(group.id);
                      updateQuery('');
                    }}
                  >
                    <Icon
                      aria-hidden='true'
                      className='size-[1.0625rem]'
                      strokeWidth={1.8}
                    />
                  </button>
                );
              })}
            </div>
          )}
        </div>

        <div
          ref={resultsRef}
          aria-busy={data.status === 'loading'}
          className='h-72 min-h-0 overflow-y-auto overscroll-contain px-3 pb-3'
        >
          {data.status === 'error' ? (
            <div className='grid h-full content-center justify-items-center gap-2 text-center text-muted-foreground'>
              <p role='status'>表情加载失败</p>
              <Button
                type='button'
                size='sm'
                variant='outline'
                onClick={loadData}
              >
                重试
              </Button>
            </div>
          ) : data.status !== 'ready' ? (
            <p
              role='status'
              className='grid h-full place-items-center text-muted-foreground'
            >
              正在加载表情…
            </p>
          ) : visibleGroups.length === 0 ? (
            <p
              role='status'
              className='grid h-full place-items-center text-muted-foreground'
            >
              没有找到匹配的表情
            </p>
          ) : (
            visibleGroups.map((group) => (
              <EmojiGrid
                key={group.id}
                group={group}
                value={value}
                onSelect={selectEmoji}
              />
            ))
          )}
        </div>
        {value && onClear && (
          <div className='flex shrink-0 justify-end border-t px-3 py-1.5'>
            <Button
              type='button'
              variant='ghost'
              size='sm'
              className='text-muted-foreground'
              onClick={() => {
                onClear();
                handleOpenChange(false);
              }}
            >
              <Eraser aria-hidden='true' />
              清除表情
            </Button>
          </div>
        )}
      </PopoverContent>
    </Popover>
  );
}
