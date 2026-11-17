import { LumaSpin } from '@qingmo/shadcn/components/luma-spin';
import { type ReactElement } from 'react';
import { Navigate, Outlet, useLocation } from 'react-router';

import { useAuthStore } from '@/stores/auth/use-auth';

const PageLoading = () => {
  return (
    <div className='flex min-h-screen w-full items-center justify-center bg-background'>
      <LumaSpin />
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
