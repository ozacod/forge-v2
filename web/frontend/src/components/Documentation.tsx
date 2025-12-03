import { useEffect, useRef, useCallback } from 'react';
import { useTheme } from '../contexts/ThemeContext';
import { DocSearch } from './DocSearch';

export function Documentation() {
  const { theme } = useTheme();
  const iframeRef = useRef<HTMLIFrameElement>(null);

  useEffect(() => {
    // Update Docusaurus theme when theme changes
    const iframe = iframeRef.current;
    if (iframe?.contentDocument) {
      const iframeRoot = iframe.contentDocument.documentElement;
      iframeRoot.setAttribute('data-theme', theme);
    }
  }, [theme]);

  const handleNavigate = useCallback((url: string) => {
    const iframe = iframeRef.current;
    if (iframe) {
      // Navigate the iframe to the search result URL
      iframe.src = url;
    }
  }, []);

  return (
    <div className="w-full h-[calc(100vh-64px)] min-h-[600px] flex flex-col">
      {/* Search bar */}
      <div className={`px-4 py-3 border-b ${
        theme === 'dark' 
          ? 'bg-black/40 border-white/10' 
          : 'bg-gray-50 border-gray-200'
      }`}>
        <div className="max-w-md">
          <DocSearch onNavigate={handleNavigate} />
        </div>
      </div>
      
      {/* Docs iframe */}
      <iframe
        ref={iframeRef}
        src="/docs/"
        className="w-full flex-1 border-0"
        title="Docusaurus Documentation"
        style={{
          display: 'block',
          margin: 0,
          padding: 0,
          width: '100%',
        }}
      />
    </div>
  );
}
