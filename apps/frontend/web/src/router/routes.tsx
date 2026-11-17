import { lazy } from 'react';
import { type RouteObject } from 'react-router';

import { RedirectRoute, ProtectRoute } from '@/layouts/route/route-layout';

const Login = lazy(() => import('@/views/login/page'));
const Signup = lazy(() => import('@/views/signup/page'));
const NotFound = lazy(() => import('@/views/not-found/not-found'));

export const routes: RouteObject[] = [
  {
    path: '/login',
    element: (
      <RedirectRoute>
        <Login />
      </RedirectRoute>
    ),
  },
  {
    path: '/signup',
    element: <Signup />,
  },
  {
    path: '/',
    element: <ProtectRoute />,
  },
  { path: '*', element: <NotFound /> },
];
