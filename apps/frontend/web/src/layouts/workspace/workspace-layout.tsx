import {
  SidebarInset,
  SidebarProvider,
} from '@qingmo/shadcn/components/sidebar';
import {
  useRef,
  type CSSProperties,
  type PropsWithChildren,
  type ReactNode,
} from 'react';

import { WorkspaceSidebar } from './components/workspace-sidebar';
import { WorkspaceSidebarResizeHandle } from './components/workspace-sidebar-resize-handle';

const SIDEBAR_DEFAULT_WIDTH = 249;
const SIDEBAR_COLLAPSED_WIDTH = 73;

const sidebarStyle = {
  '--sidebar-width': `${SIDEBAR_DEFAULT_WIDTH}px`,
  '--sidebar-width-icon': `${SIDEBAR_COLLAPSED_WIDTH}px`,
} as CSSProperties;

export function WorkspaceLayout({ children }: PropsWithChildren) {
  const sidebarRootRef = useRef<HTMLDivElement>(null);

  return (
    <SidebarProvider
      ref={sidebarRootRef}
      style={sidebarStyle}
      className='data-[sidebar-resizing=dragging]:[&>[data-state=expanded]>:is([data-slot=sidebar-gap],[data-slot=sidebar-container])]:transition-none'
      onTransitionEnd={({ currentTarget, propertyName, target }) => {
        if (
          currentTarget.dataset.sidebarResizing === 'transitioning' &&
          propertyName === 'width' &&
          target instanceof HTMLElement &&
          target.dataset.slot === 'sidebar-gap'
        ) {
          currentTarget.dataset.sidebarResizing = 'dragging';
        }
      }}
    >
      <WorkspaceSidebar />
      <WorkspaceSidebarResizeHandle
        sidebarRootRef={sidebarRootRef}
        defaultWidth={SIDEBAR_DEFAULT_WIDTH}
        collapsedWidth={SIDEBAR_COLLAPSED_WIDTH}
      />
      <SidebarInset className='min-w-0 bg-background'>{children}</SidebarInset>
    </SidebarProvider>
  );
}

type WorkspacePageProps = PropsWithChildren<{ header: ReactNode }>;

export function WorkspacePage({ header, children }: WorkspacePageProps) {
  return (
    <div className='flex min-w-0 flex-1 flex-col bg-background'>
      <header className='w-full h-17 flex border-b border-border items-center px-7'>
        {header}
      </header>
      <main className='w-full min-h-0 flex-1'>{children}</main>
    </div>
  );
}
