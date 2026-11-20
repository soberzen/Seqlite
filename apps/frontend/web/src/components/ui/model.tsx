import { Button } from '@qingmo/shadcn/components/button';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@qingmo/shadcn/components/dialog';
import { cn } from '@qingmo/shadcn/lib/utils';
import { XIcon } from 'lucide-react';
import {
  useRef,
  type PropsWithChildren,
  type ReactNode,
  type RefObject,
} from 'react';

export type ModelProps = PropsWithChildren<{
  open?: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  footer?: ReactNode;
  initialFocusRef?: RefObject<HTMLElement | null>;
  contentClassName?: string;
  headerClassName?: string;
  footerClassName?: string;
  onClose?: () => void;
}>;

export function Model({
  onOpenChange,
  title,
  footer,
  initialFocusRef,
  children,
  contentClassName,
  headerClassName,
  footerClassName,
  onClose,
  ...dialogProps
}: ModelProps) {
  const returnFocusRef = useRef<HTMLElement>(null);

  function handleOpenChange(nextOpen: boolean) {
    onOpenChange(nextOpen);
    if (!nextOpen) onClose?.();
  }

  return (
    <Dialog
      {...dialogProps}
      onOpenChange={handleOpenChange}
    >
      <DialogContent
        showCloseButton={false}
        className={cn(
          'max-h-[calc(100dvh-2rem)] w-[calc(100%-2rem)] gap-4 overflow-y-auto overscroll-contain rounded-[16px] border border-border bg-surface-subtle p-6 shadow-[0_8px_24px_var(--shadow-dialog)] ring-0 sm:max-w-120 ',
          contentClassName,
        )}
        onOpenAutoFocus={(event) => {
          const activeElement = document.activeElement;
          returnFocusRef.current =
            activeElement instanceof HTMLElement ? activeElement : null;
          const initialFocus = initialFocusRef?.current;
          if (initialFocus) {
            event.preventDefault();
            initialFocus.focus();
          }
        }}
        onCloseAutoFocus={(event) => {
          const previousFocus = returnFocusRef.current;
          if (previousFocus?.isConnected) {
            event.preventDefault();
            previousFocus.focus();
          }
        }}
      >
        <DialogHeader
          className={cn(
            'min-h-8 min-w-0 flex-row items-center pr-10',
            headerClassName,
          )}
        >
          <DialogTitle className='text-xl/7 font-semibold text-dialog-foreground'>
            {title}
          </DialogTitle>
        </DialogHeader>
        {children}
        {footer ? (
          <DialogFooter
            className={cn(
              'm-0 flex-row flex-wrap items-center justify-end rounded-none border-t-0 bg-transparent p-0',
              footerClassName,
            )}
          >
            {footer}
          </DialogFooter>
        ) : null}
        <DialogClose asChild>
          <Button
            type='button'
            variant='ghost'
            size='icon'
            aria-label='关闭'
            className='absolute top-6 right-6 text-dialog-secondary'
          >
            <XIcon
              aria-hidden='true'
              className='size-4.5'
            />
          </Button>
        </DialogClose>
      </DialogContent>
    </Dialog>
  );
}
