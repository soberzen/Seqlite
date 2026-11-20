import { Toaster } from '@qingmo/shadcn/components/sonner';
import { TooltipProvider } from '@qingmo/shadcn/components/tooltip';
import { RouterProvider } from 'react-router';

import { ThemeProvider } from './context/theme-context';
import router from './router';

export default function App() {
  return (
    <ThemeProvider>
      <TooltipProvider>
        <RouterProvider router={router} />
        <Toaster />
      </TooltipProvider>
    </ThemeProvider>
  );
}
