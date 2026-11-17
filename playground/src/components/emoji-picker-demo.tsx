import { Button } from '@qingmo/shadcn/components/button';
import { EmojiPicker } from '@qingmo/shadcn/components/emoji-picker';
import { SmilePlus } from 'lucide-react';
import { useState } from 'react';

export function EmojiPickerDemo() {
  const [emoji, setEmoji] = useState<string | null>('📚');

  return (
    <div className='mt-4 flex items-center gap-3'>
      <EmojiPicker
        value={emoji}
        onSelect={setEmoji}
        onClear={() => setEmoji(null)}
      >
        <Button
          type='button'
          variant='outline'
          size='icon-lg'
          className='size-11 rounded-xl text-2xl'
          aria-label='更换空间图标'
        >
          {emoji ? (
            <span aria-hidden='true'>{emoji}</span>
          ) : (
            <SmilePlus
              aria-hidden='true'
              className='size-5 text-muted-foreground'
            />
          )}
        </Button>
      </EmojiPicker>
      <div>
        <p className='text-sm font-medium'>轻墨团队</p>
        <p className='text-xs text-muted-foreground'>工作空间</p>
      </div>
      <span
        role='status'
        className='sr-only'
      >
        {emoji ? `当前图标为 ${emoji}` : '尚未设置空间图标'}
      </span>
    </div>
  );
}
