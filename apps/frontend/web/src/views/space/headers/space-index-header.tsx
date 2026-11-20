import { HouseIcon } from 'lucide-react';

export function SpaceIndexHeader() {
  return (
    <div className='flex items-center gap-1.75'>
      <HouseIcon
        size={18}
        className='text-muted-foreground'
      />
      <span className='text-foreground text-sm'>开始使用</span>
    </div>
  );
}
