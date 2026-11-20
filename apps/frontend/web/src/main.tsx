import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';

import './index.css';
import App from './App.tsx';
import { useAuthStore } from './stores/auth/use-auth.ts';
import { getToken } from './utils/auth';

const app = createRoot(document.getElementById('root')!);
// const pathname = window.location.pathname;

// if (getToken() && pathname !== '/login' && pathname !== '/signup')
//   useAuthStore.getState().refreshUser();

app.render(
  <StrictMode>
    <App />
  </StrictMode>,
);
