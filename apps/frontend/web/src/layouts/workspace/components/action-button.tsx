import { cn } from '@qingmo/shadcn/lib/utils';
import type { ComponentPropsWithRef, ReactNode } from 'react';

type ActionButtonProps = ComponentPropsWithRef<'button'> & {
  icon: ReactNode;
  label: string;
  action?: () => void;
};

export default function ActionButton({
  icon,
  label,
  action,
  onClick,
  className,
  ...buttonProps
}: ActionButtonProps) {
  return (
    <button
      type='button'
      {...buttonProps}
      aria-label={label}
      className={cn(
        'inline-flex size-6 shrink-0 cursor-pointer items-center justify-center rounded-sm text-muted-foreground transition-colors hover:bg-border focus-visible:bg-border',
        className,
      )}
      onClick={(event) => {
        onClick?.(event);
        if (!event.defaultPrevented) action?.();
      }}
    >
      {icon}
    </button>
  );
}
