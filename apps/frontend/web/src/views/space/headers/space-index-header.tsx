import { FolderClosedIcon } from 'lucide-react';

export function SpaceIndexHeader() {
  return (
    <nav
      aria-label='面包屑'
      className='flex items-center gap-2.5 text-sm'
    >
      <FolderClosedIcon
        size={16}
        className='text-subtle'
        aria-hidden
      />
      <span className='text-foreground'>空间</span>
    </nav>
  );
}
