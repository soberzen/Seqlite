import { Button } from '@qingmo/shadcn/components/button';
import { FolderPlusIcon, PlusIcon } from 'lucide-react';

import { useOpenSpaceDialog } from '@/context/space-dialog-context';
import { WorkspacePage } from '@/layouts/workspace/workspace-layout';

import { SpaceIndexHeader } from './headers/space-index-header';

export default function SpaceIndexPage() {
  const open = useOpenSpaceDialog();
  return (
    <WorkspacePage header={<SpaceIndexHeader />}>
      <div className='w-full h-full flex flex-col items-center justify-center'>
        <div className='w-16 h-16 flex items-center justify-center rounded-[18px] bg-accent'>
          <FolderPlusIcon
            className='text-foreground'
            size={28}
          />
        </div>
        <span className='text-[13px] text-muted-foreground mt-4.5'>
          欢迎使用轻墨文档
        </span>
        <h2 className='font-semibold text-foreground mt-4.5 text-3xl'>
          创建你的第一个空间
        </h2>
        <p className='mt-4.5 text-muted-foreground text-sm text-center leading-[1.8]'>
          把相关文档放在一起，让知识井然有序。
          <br />
          你可以独自记录，也可以邀请成员一起协作。
        </p>
        <div className='flex flex-col items-center gap-4 mt-4.5 pt-2'>
          <Button
            size='lg'
            className='w-30 h-11 rouned-[7px]'
            onClick={() => open({ mode: 'create' })}
          >
            <PlusIcon />
            <span>创建空间</span>
          </Button>
          <span className='text-muted-foreground text-[12px]'>
            创建空间 → 新建文档 → 开始记录
          </span>
        </div>
        <div className='flex flex-col items-center mt-4.5 pt-6 gap-2'>
          <div className='w-12 h-px bg-border' />
          <span className='text-muted-foreground text-[12px]'>
            收到团队邀请？通过邀请链接即可加入空间。
          </span>
        </div>
      </div>
    </WorkspacePage>
  );
}
