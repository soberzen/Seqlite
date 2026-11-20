import { useSidebar } from '@qingmo/shadcn/components/sidebar';
import {
  useRef,
  useState,
  type KeyboardEvent,
  type PointerEvent,
  type RefObject,
} from 'react';

const MIN_WIDTH = 220;
const MAX_WIDTH = 380;
const DRAG_START_DISTANCE = 5;
const COLLAPSE_AT_WIDTH = MIN_WIDTH * 0.5;
const EXPAND_DRAG_DISTANCE = MIN_WIDTH * 0.2;
const KEYBOARD_STEP = 16;

type WorkspaceSidebarResizeHandleProps = {
  sidebarRootRef: RefObject<HTMLDivElement | null>;
  defaultWidth: number;
  collapsedWidth: number;
};

type ResizeDrag = {
  pointerId: number;
  startX: number;
  startY: number;
  startWidth: number;
  width: number;
  hasMoved: boolean;
  open: boolean;
  expandAt: number;
};

export function WorkspaceSidebarResizeHandle({
  sidebarRootRef,
  defaultWidth,
  collapsedWidth,
}: WorkspaceSidebarResizeHandleProps) {
  const { open, isMobile, setOpen, toggleSidebar } = useSidebar();
  const [width, setWidth] = useState(defaultWidth);
  const drag = useRef<ResizeDrag | null>(null);
  const visibleWidth = open ? width : collapsedWidth;

  function applyWidth(nextWidth: number) {
    const clampedWidth = Math.round(
      Math.min(MAX_WIDTH, Math.max(MIN_WIDTH, nextWidth)),
    );
    sidebarRootRef.current?.style.setProperty(
      '--sidebar-width',
      `${clampedWidth}px`,
    );
    return clampedWidth;
  }

  function setOpenDuringDrag(current: ResizeDrag, nextOpen: boolean) {
    current.open = nextOpen;
    sidebarRootRef.current?.setAttribute(
      'data-sidebar-resizing',
      'transitioning',
    );
    setOpen(nextOpen);
  }

  function startDrag(event: PointerEvent<HTMLDivElement>) {
    if (event.button !== 0 || drag.current) return;

    event.preventDefault();
    event.currentTarget.setPointerCapture(event.pointerId);
    drag.current = {
      pointerId: event.pointerId,
      startX: event.clientX,
      startY: event.clientY,
      startWidth: visibleWidth,
      width,
      hasMoved: false,
      open,
      expandAt: collapsedWidth + EXPAND_DRAG_DISTANCE,
    };
  }

  function moveDrag(event: PointerEvent<HTMLDivElement>) {
    const current = drag.current;
    if (isMobile || current?.pointerId !== event.pointerId) return;

    if (!current.hasMoved) {
      if (isWithinClickDistance(event, current)) return;
      current.hasMoved = true;
      sidebarRootRef.current?.setAttribute('data-sidebar-resizing', 'dragging');
    }

    const requestedWidth = current.startWidth + event.clientX - current.startX;
    if (current.open && requestedWidth <= COLLAPSE_AT_WIDTH) {
      current.width = applyWidth(MIN_WIDTH);
      // Require a deliberate reversal before re-opening in the same gesture.
      current.expandAt = COLLAPSE_AT_WIDTH + EXPAND_DRAG_DISTANCE;
      setOpenDuringDrag(current, false);
      return;
    }

    if (!current.open) {
      if (requestedWidth < current.expandAt) return;
      setOpenDuringDrag(current, true);
    }

    current.width = applyWidth(requestedWidth);
  }

  function finishDrag(event: PointerEvent<HTMLDivElement>) {
    const current = drag.current;
    if (current?.pointerId !== event.pointerId) return;

    drag.current = null;
    sidebarRootRef.current?.removeAttribute('data-sidebar-resizing');
    setWidth(current.width);

    if (
      event.type === 'pointerup' &&
      !isMobile &&
      !current.hasMoved &&
      isWithinClickDistance(event, current)
    ) {
      toggleSidebar();
    }
  }

  function resizeWithKeyboard(event: KeyboardEvent<HTMLDivElement>) {
    if (drag.current) return;

    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      if (!event.repeat) toggleSidebar();
      return;
    }

    const nextWidths: Record<string, number | undefined> = {
      ArrowLeft: visibleWidth - KEYBOARD_STEP,
      ArrowRight: visibleWidth + KEYBOARD_STEP,
      Home: MIN_WIDTH,
      End: MAX_WIDTH,
    };
    const nextWidth = nextWidths[event.key];
    if (typeof nextWidth !== 'number') return;

    event.preventDefault();
    if (nextWidth < collapsedWidth) return;
    setWidth(applyWidth(nextWidth));
    if (!open) setOpen(true);
  }

  return (
    <div
      hidden={isMobile}
      className='sticky top-0 z-20 h-svh w-0 shrink-0 max-md:hidden'
    >
      <div
        role='separator'
        tabIndex={0}
        aria-label='调整侧栏宽度'
        aria-controls='workspace-sidebar'
        aria-orientation='vertical'
        aria-valuemin={collapsedWidth}
        aria-valuemax={MAX_WIDTH}
        aria-valuenow={visibleWidth}
        className='absolute inset-y-0 left-0 w-4 -translate-x-1/2 cursor-col-resize touch-none select-none outline-none after:pointer-events-none after:absolute after:inset-y-0 after:left-1/2 after:w-px after:-translate-x-1/2 hover:after:bg-ring focus-visible:after:bg-ring'
        onPointerDown={startDrag}
        onPointerMove={moveDrag}
        onPointerUp={finishDrag}
        onPointerCancel={finishDrag}
        onLostPointerCapture={finishDrag}
        onKeyDown={resizeWithKeyboard}
      />
    </div>
  );
}

function isWithinClickDistance(
  event: PointerEvent<HTMLDivElement>,
  drag: ResizeDrag,
) {
  return (
    Math.abs(event.clientX - drag.startX) <= DRAG_START_DISTANCE &&
    Math.abs(event.clientY - drag.startY) <= DRAG_START_DISTANCE
  );
}
