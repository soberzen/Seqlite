import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from '@qingmo/shadcn/components/avatar';
import { Button } from '@qingmo/shadcn/components/button';
import {
  ChevronRightIcon,
  FolderClosedIcon,
  FolderPlusIcon,
  FolderXIcon,
  PlusIcon,
} from 'lucide-react';
import { useShallow } from 'zustand/shallow';

import { useOpenSpaceDialog } from '@/context/space-dialog-context';
import { PageLoading } from '@/layouts/route/route-layout';
import { WorkspacePage } from '@/layouts/workspace/workspace-layout';
import { useAuthStore } from '@/stores/auth/use-auth';
import {
  useSpaceListStore,
  type SpaceListItem,
} from '@/stores/space/use-space-list';
import { getAvatarText } from '@/utils/avatar';

import { SpaceIndexHeader } from './headers/space-index-header';

function SpaceEmpty() {
  const open = useOpenSpaceDialog();

  return (
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
  );
}

function SpaceError({ onRetry }: { onRetry: () => void }) {
  return (
    <div className='flex h-full w-full flex-col items-center justify-center gap-4 text-center'>
      <div className='flex size-16 items-center justify-center rounded-[18px] bg-accent'>
        <FolderXIcon
          className='text-muted-foreground'
          size={28}
          aria-hidden
        />
      </div>

      <h2 className='text-lg font-semibold text-foreground'>空间加载失败</h2>

      <Button
        type='button'
        size='lg'
        onClick={onRetry}
      >
        重试
      </Button>
    </div>
  );
}

function SpaceList({ spaces }: { spaces: SpaceListItem[] }) {
  const open = useOpenSpaceDialog();
  const user = useAuthStore((state) => state.user);
  return (
    <div className='flex w-full flex-col gap-6 px-4 pt-6 pb-7 sm:px-7 lg:px-14 lg:pt-11.5'>
      <div className='flex justify-between'>
        <h2 className='text-foreground text-2xl font-semibold'>全部空间</h2>
        <Button
          size='lg'
          className='w-30 h-11 rouned-[7px]'
          onClick={() => open({ mode: 'create' })}
        >
          <PlusIcon />
          <span>创建空间</span>
        </Button>
      </div>
      <ul
        aria-label='空间列表'
        className='grid grid-cols-[repeat(auto-fit,minmax(min(100%,240px),1fr))] gap-4'
      >
        {spaces.map((space) => (
          <li
            key={space.id}
            className='group flex h-[236px] min-w-0 flex-col rounded-[10px] border border-border bg-background p-[18px] transition-colors duration-150 hover:border-icon-faint hover:bg-surface-subtle motion-reduce:transition-none'
          >
            <span
              aria-hidden
              className='mb-4 flex h-10 shrink-0 items-center justify-between'
            >
              <span className='flex size-8 items-center justify-center text-2xl'>
                {space.emoji || <FolderClosedIcon size={24} />}
              </span>
              <ChevronRightIcon
                size={14}
                className='text-muted-foreground transition duration-150 group-hover:text-foreground motion-safe:group-hover:translate-x-0.5 motion-reduce:transition-none'
              />
            </span>
            <h3
              className='line-clamp-2 text-base font-semibold text-foreground wrap-anywhere'
              title={space.name}
            >
              {space.name}
            </h3>
            {space.description?.trim() && (
              <p
                className='mt-2 line-clamp-2 text-sm text-muted-foreground wrap-anywhere'
                title={space.description}
              >
                {space.description}
              </p>
            )}
            <div className='mt-auto flex min-w-0 shrink-0 items-center gap-1.5 pt-2 text-xs text-muted-foreground'>
              {user && (
                <>
                  <Avatar
                    size='sm'
                    aria-hidden
                  >
                    {user.avatarURL && (
                      <AvatarImage
                        src={user.avatarURL}
                        alt=''
                      />
                    )}
                    <AvatarFallback className='bg-avatar-neutral'>
                      {getAvatarText(user.userName)}
                    </AvatarFallback>
                  </Avatar>
                  <span className='truncate'>{user.userName}</span>
                </>
              )}
              {space.updatedAtLabel && (
                <span className='ml-auto shrink-0'>{space.updatedAtLabel}</span>
              )}
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}

function SpaceIndexContent() {
  const [spaces, status, fetchSpaceList] = useSpaceListStore(
    useShallow((state) => [state.spaces, state.status, state.fetchSpaceList]),
  );

  if (status === 'loading') {
    return (
      <PageLoading className='h-full w-full gap-3'>
        <span className='text-sm text-muted-foreground'>正在加载空间…</span>
      </PageLoading>
    );
  }

  if (status === 'error') {
    return <SpaceError onRetry={fetchSpaceList} />;
  }

  return spaces.length === 0 ? <SpaceEmpty /> : <SpaceList spaces={spaces} />;
}

export default function SpaceIndexPage() {
  return (
    <WorkspacePage header={<SpaceIndexHeader />}>
      <SpaceIndexContent />
    </WorkspacePage>
  );
}
