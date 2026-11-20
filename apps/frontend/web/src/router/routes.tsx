import { lazy } from 'react';
import { Outlet, type RouteObject } from 'react-router';

import { SpaceDialogProvider } from '@/context/space-dialog-context';
import { RedirectRoute, ProtectRoute } from '@/layouts/route/route-layout';
import { WorkspaceLayout } from '@/layouts/workspace/workspace-layout';

const Login = lazy(() => import('@/views/login/page'));
const Signup = lazy(() => import('@/views/signup/page'));
const NotFound = lazy(() => import('@/views/not-found/not-found'));
const SpaceIndexPage = lazy(() => import('@/views/space/space-index-page'));

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
    element: (
      <SpaceDialogProvider>
        <WorkspaceLayout>
          <Outlet />
        </WorkspaceLayout>
      </SpaceDialogProvider>
    ),
    children: [
      {
        path: 'recent',
      },
      {
        path: 'favorites',
      },
      {
        path: 'space',
        element: <SpaceIndexPage />,
      },
      {
        path: 'space/:spaceId',
      },
      {
        path: 'wiki/:docId',
      },
    ],
  },
  // {
  //   element: <ProtectRoute />,
  //   children: [

  //   ],
  // },
  { path: '*', element: <NotFound /> },
];
