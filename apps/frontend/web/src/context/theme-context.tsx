import {
  createContext,
  useContext,
  useEffect,
  useLayoutEffect,
  useMemo,
  useState,
  useSyncExternalStore,
  type PropsWithChildren,
} from 'react';

type Theme = 'light' | 'dark' | 'system';
type ResolvedTheme = Omit<Theme, 'system'>;

const ThemeContext = createContext<{
  setTheme: (theme: Theme) => void;
} | null>(null);

const darkQuery = '(prefers-color-scheme: dark)';

function subscribeToSystemTheme(onChange: () => void) {
  const query = window.matchMedia(darkQuery);
  query.addEventListener('change', onChange);
  return () => query.removeEventListener('change', onChange);
}

function getSystemDark() {
  return window.matchMedia(darkQuery).matches;
}

export function ThemeProvider({ children }: PropsWithChildren) {
  const [theme, setTheme] = useState<Theme>(() => {
    const saved = window.localStorage.getItem('theme');
    return saved === 'light' || saved === 'dark' ? saved : 'system';
  });

  const isSystemDark = useSyncExternalStore(
    subscribeToSystemTheme,
    getSystemDark,
    () => false,
  );

  const resolvedTheme: ResolvedTheme = useMemo(
    () => (theme === 'system' ? (isSystemDark ? 'dark' : 'light') : theme),
    [theme, isSystemDark],
  );

  useLayoutEffect(() => {
    document.documentElement.classList.toggle('dark', resolvedTheme === 'dark');
  }, [resolvedTheme]);

  useEffect(() => {
    window.localStorage.setItem('theme', theme);
  }, [theme]);

  return (
    <ThemeContext.Provider value={{ setTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (!context) throw new Error('useTheme must be used within ThemeProvider');
  return context;
}
