import { LumaSpin } from '@qingmo/shadcn/components/luma-spin';
import { cn } from '@qingmo/shadcn/lib/utils';
import { type PropsWithChildren, type ReactElement } from 'react';
import { Navigate, Outlet, useLocation } from 'react-router';

import { useAuthStore } from '@/stores/auth/use-auth';

export const PageLoading = ({
  children,
  className,
}: PropsWithChildren<{ className?: string }>) => {
  return (
    <div
      className={cn(
        'flex h-full w-full items-center justify-center bg-background flex-col gap-2',
        className,
      )}
    >
      <LumaSpin />
      {children}
    </div>
  );
};

export function ProtectRoute() {
  const loading = useAuthStore((state) => state.loading);
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const location = useLocation();

  if (loading) {
    return <PageLoading />;
  }

  if (!isAuthenticated) {
    const redirect = `${location.pathname}${location.search}`;
    return (
      <Navigate
        to={`/login?redirect=${encodeURIComponent(redirect)}`}
        replace
      />
    );
  }

  return <Outlet />;
}

export function RedirectRoute({ children }: { children: ReactElement }) {
  const loading = useAuthStore((state) => state.loading);
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  const location = useLocation();

  if (loading) return <PageLoading />;

  if (isAuthenticated) {
    const params = new URLSearchParams(location.search);
    const redirectPath = params.get('redirect') || '/';
    return (
      <Navigate
        to={redirectPath}
        replace
      />
    );
  }
  return children;
}
