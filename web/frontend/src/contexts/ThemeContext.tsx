import { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';

type Theme = 'light' | 'dark';

interface ThemeContextType {
  theme: Theme;
  toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setTheme] = useState<Theme>(() => {
    // Check localStorage first
    const savedTheme = localStorage.getItem('theme') as Theme;
    if (savedTheme) {
      return savedTheme;
    }
    // Check system preference
    if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
      return 'dark';
    }
    return 'dark'; // Default to dark
  });

  useEffect(() => {
    const root = document.documentElement;
    root.classList.remove('dark', 'light');
    root.classList.add(theme);
    localStorage.setItem('theme', theme);
    
    // Also update Docusaurus iframe if it exists
    const updateIframeTheme = () => {
      const iframe = document.querySelector('iframe[src*="/docs/"]') as HTMLIFrameElement;
      if (iframe?.contentDocument) {
        const iframeRoot = iframe.contentDocument.documentElement;
        iframeRoot.setAttribute('data-theme', theme);
      }
    };
    
    updateIframeTheme();
    // Also listen for iframe load events
    const iframe = document.querySelector('iframe[src*="/docs/"]') as HTMLIFrameElement;
    if (iframe) {
      iframe.addEventListener('load', updateIframeTheme);
      return () => iframe.removeEventListener('load', updateIframeTheme);
    }
  }, [theme]);

  const toggleTheme = () => {
    setTheme(prev => prev === 'dark' ? 'light' : 'dark');
  };

  return (
    <ThemeContext.Provider value={{ theme, toggleTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (context === undefined) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
}

