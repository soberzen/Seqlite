import {
  Avatar,
  AvatarImage,
  AvatarFallback,
} from '@qingmo/shadcn/components/avatar';
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupAction,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarFooter,
  useSidebar,
} from '@qingmo/shadcn/components/sidebar';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@qingmo/shadcn/components/tooltip';
import {
  PanelLeftCloseIcon,
  PanelLeftOpenIcon,
  SearchIcon,
  CommandIcon,
  Clock3Icon,
  StarIcon,
  PlusIcon,
  SlidersHorizontalIcon,
  Trash2Icon,
  UserPlusIcon,
  Settings2Icon,
} from 'lucide-react';
import { Link } from 'react-router';

import { useOpenSpaceDialog } from '@/context/space-dialog-context';
import { useAuthStore } from '@/stores/auth/use-auth';
import { useIsMac } from '@/stores/os/os';
import { getAvatarText } from '@/utils/avatar';

import ActionButton from './action-button';

function SidebarTrigger() {
  const { open, toggleSidebar } = useSidebar();
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <ActionButton
          label={open ? '关闭' : '展开' + '侧边栏'}
          icon={
            open ? (
              <PanelLeftCloseIcon size={18} />
            ) : (
              <PanelLeftOpenIcon size={18} />
            )
          }
          onClick={toggleSidebar}
        />
      </TooltipTrigger>
      <TooltipContent>{open ? '关闭' : '展开'}侧边栏</TooltipContent>
    </Tooltip>
  );
}

type SidebarSearchProps = {
  onClick?: () => void;
};

function SidebarSearch({ onClick }: SidebarSearchProps) {
  const isMac = useIsMac();
  return (
    <button
      type='button'
      aria-label='搜索文档'
      onClick={onClick}
      className='flex justify-between items-center h-9 w-full px-2.5 mt-3 rounded-[7px] text-muted-foreground bg-sidebar-search cursor-pointer hover:bg-control-active active:bg-selection group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:mt-2'
    >
      <div className='flex items-center gap-2 h-full'>
        <SearchIcon size={17} />
        <span className='text-sm group-data-[collapsible=icon]:hidden'>
          搜索文档
        </span>
      </div>
      <kbd className='ml-auto inline-flex shrink-0 items-center gap-1 whitespace-nowrap text-[11px] leading-4 text-muted-foreground group-data-[collapsible=icon]:hidden'>
        {isMac ? (
          <CommandIcon
            size={12}
            aria-label='Command'
          />
        ) : (
          <span>Ctrl</span>
        )}
        <span>K</span>
      </kbd>
    </button>
  );
}

function SidebarRecent() {
  return (
    <Link
      to='/workspace/recent'
      aria-label='最近打开'
      className='w-full h-9 mt-4.5 px-3 flex items-center gap-2.25 rounded-[7px] hover:bg-control-active active:bg-selection group-data-[collapsible=icon]:mt-4.5 group-data-[collapsible=icon]:justify-center'
    >
      <Clock3Icon
        size={18}
        className='text-muted-foreground'
      />
      <span className='text-foreground text-sm group-data-[collapsible=icon]:hidden'>
        最近打开
      </span>
    </Link>
  );
}

function SidebarFavorites() {
  return (
    <Link
      to='/workspace/favorites'
      aria-label='我的收藏'
      className='w-full h-9 px-3 flex items-center gap-2.25 rounded-[7px] hover:bg-control-active active:bg-selection group-data-[collapsible=icon]:justify-center'
    >
      <StarIcon
        size={18}
        className='text-muted-foreground'
      />
      <span className='text-foreground text-sm group-data-[collapsible=icon]:hidden'>
        我的收藏
      </span>
    </Link>
  );
}

export function WorkspaceSidebar() {
  const openSpaceDialog = useOpenSpaceDialog();
  return (
    <Sidebar
      id='workspace-sidebar'
      collapsible='icon'
      className='px-3 pt-6 pb-4 bg-sidebar'
    >
      <div className='flex-1'>
        <SidebarHeader className='flex justify-between min-h-8 px-2 items-center py-0'>
          <div className='flex justify-between w-full min-h-8 items-center group-data-[collapsible=icon]:flex-col group-data-[collapsible=icon]:items-center group-data-[collapsible=icon]:gap-3'>
            <div className='flex items-center gap-2.5'>
              <img
                width={30}
                height={30}
                src='/logo.png'
                alt='轻墨文档'
                className='object-contain'
              />
              <h1 className='text-[19px] font-semibold overflow-hidden whitespace-nowrap group-data-[collapsible=icon]:hidden'>
                轻墨文档
              </h1>
            </div>
            <SidebarTrigger />
          </div>
        </SidebarHeader>
        <SidebarSearch />
        <SidebarRecent />
        <SidebarFavorites />
        <SidebarContent className='w-full mt-7'>
          <SidebarGroup className='h-7 w-full flex-row items-center justify-between px-3 py-0 text-muted-foreground group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0'>
            <SidebarGroupLabel className='h-full px-0 group-data-[collapsible=icon]:hidden'>
              空间
            </SidebarGroupLabel>
            <div className='flex items-center gap-1 group-data-[collapsible=icon]:gap-0 '>
              <SidebarGroupAction
                asChild
                className='static size-6 shrink-0 group-data-[collapsible=icon]:inline-flex'
              >
                <ActionButton
                  label='管理空间'
                  className='text-subtle'
                  icon={<SlidersHorizontalIcon size={14} />}
                />
              </SidebarGroupAction>

              <SidebarGroupAction
                asChild
                className='static size-6 shrink-0 group-data-[collapsible=icon]:inline-flex'
              >
                <ActionButton
                  className='text-muted-foreground'
                  label='创建空间'
                  icon={<PlusIcon size={16} />}
                  onClick={() => {
                    openSpaceDialog({ mode: 'create' });
                  }}
                />
              </SidebarGroupAction>
            </div>
          </SidebarGroup>
        </SidebarContent>
      </div>
      <WorkspaceSidebarFooter />
    </Sidebar>
  );
}

function SidebarInvitation() {
  return (
    <button
      type='button'
      aria-label='邀请协作'
      className='w-full h-9 px-3 flex items-center gap-2.25 rounded-[7px] hover:bg-control-active active:bg-selection group-data-[collapsible=icon]:justify-center'
    >
      <UserPlusIcon
        size={18}
        className='text-muted-foreground'
      />
      <span className='text-muted-foreground text-[13px] group-data-[collapsible=icon]:hidden'>
        邀请协作
      </span>
    </button>
  );
}

function SidebarTrash() {
  return (
    <button
      type='button'
      aria-label='回收站'
      className='w-full h-9 px-3 flex items-center gap-2.25 rounded-[7px] hover:bg-control-active active:bg-selection group-data-[collapsible=icon]:justify-center'
    >
      <Trash2Icon
        size={18}
        className='text-muted-foreground'
      />
      <span className='text-muted-foreground text-[13px] group-data-[collapsible=icon]:hidden'>
        回收站
      </span>
    </button>
  );
}

function SidebarUser() {
  const user = useAuthStore((state) => state.user);
  if (!user) return null;

  return (
    <button
      type='button'
      aria-label='个人信息'
      className='w-full h-11 px-2 flex justify-between items-center mt-3.75 rounded-[7px] hover:bg-control-active active:bg-selection group-data-[collapsible=icon]:mt-3'
    >
      <div className='h-full flex items-center gap-2.5'>
        <Avatar className='text-secondary w-7.5 h-7.5 after:w-7.5 after:h-7.5 after:border-0'>
          {user.avatarURL && (
            <AvatarImage
              src={user.avatarURL}
              alt={user.userName}
            />
          )}
          <AvatarFallback className='bg-avatar-neutral'>
            {getAvatarText(user.userName)}
          </AvatarFallback>
        </Avatar>
        <div className='flex flex-col justify-between items-start group-data-[collapsible=icon]:hidden'>
          <span className='text-foreground font-medium'>{user.userName}</span>
          <span className='text-muted-foreground text-[10px]'>
            {user.email}
          </span>
        </div>
      </div>
      <Settings2Icon
        size={18}
        className='text-muted-foreground group-data-[collapsible=icon]:hidden'
      />
    </button>
  );
}

function WorkspaceSidebarFooter() {
  return (
    <SidebarFooter className='p-0'>
      <SidebarInvitation />
      <SidebarTrash />
      <SidebarUser />
    </SidebarFooter>
  );
}
