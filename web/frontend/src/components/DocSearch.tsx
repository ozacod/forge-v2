import { useState, useEffect, useRef, useCallback } from 'react';
import { useTheme } from '../contexts/ThemeContext';

interface SearchDocument {
  i: number;
  t: string;  // title
  u: string;  // url
  b?: string[]; // breadcrumbs
  h?: string; // heading hash
  p?: number; // parent id
  s?: string; // section title
}

interface SearchSection {
  documents: SearchDocument[];
  index: unknown;
}

interface SearchResult {
  title: string;
  url: string;
  breadcrumbs?: string[];
}

interface DocSearchProps {
  onNavigate: (url: string) => void;
}

export function DocSearch({ onNavigate }: DocSearchProps) {
  const { theme } = useTheme();
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<SearchResult[]>([]);
  const [isOpen, setIsOpen] = useState(false);
  const [documents, setDocuments] = useState<SearchDocument[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedIndex, setSelectedIndex] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  // Load search index on mount
  useEffect(() => {
    setIsLoading(true);
    fetch('/docs/search-index.json')
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch');
        return res.json();
      })
      .then((data: SearchSection[]) => {
        // The search index has multiple sections, each with documents
        const allDocs: SearchDocument[] = [];
        if (Array.isArray(data)) {
          data.forEach((section) => {
            if (section.documents && Array.isArray(section.documents)) {
              allDocs.push(...section.documents);
            }
          });
        }
        console.log('Loaded search documents:', allDocs.length);
        setDocuments(allDocs);
        setIsLoading(false);
      })
      .catch(err => {
        console.error('Failed to load search index:', err);
        setIsLoading(false);
      });
  }, []);

  // Simple search function
  const search = useCallback((searchQuery: string) => {
    if (!searchQuery.trim()) {
      setResults([]);
      return;
    }

    if (documents.length === 0) {
      console.log('No documents loaded yet');
      return;
    }

    const queryLower = searchQuery.toLowerCase().trim();
    const matches: SearchResult[] = [];
    const seen = new Set<string>();

    documents.forEach(doc => {
      const title = doc.t || '';
      // Build URL - add hash if present
      let url = doc.u || '';
      if (doc.h) {
        url = url + doc.h;
      }
      
      const titleLower = title.toLowerCase();
      
      // Check if title contains the query
      if (titleLower.includes(queryLower) && !seen.has(url)) {
        seen.add(url);
        matches.push({
          title,
          url,
          breadcrumbs: doc.b,
        });
      }
    });

    // Sort by relevance (exact match first, then starts with, then contains)
    matches.sort((a, b) => {
      const aTitle = a.title.toLowerCase();
      const bTitle = b.title.toLowerCase();
      
      // Exact match
      if (aTitle === queryLower) return -1;
      if (bTitle === queryLower) return 1;
      
      // Starts with
      if (aTitle.startsWith(queryLower) && !bTitle.startsWith(queryLower)) return -1;
      if (bTitle.startsWith(queryLower) && !aTitle.startsWith(queryLower)) return 1;
      
      // Shorter titles first (more relevant)
      return aTitle.length - bTitle.length;
    });

    console.log('Search results for', searchQuery, ':', matches.length);
    setResults(matches.slice(0, 15));
    setSelectedIndex(0);
  }, [documents]);

  // Search when query changes
  useEffect(() => {
    search(query);
  }, [query, search]);

  // Handle keyboard navigation
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setSelectedIndex(prev => Math.min(prev + 1, results.length - 1));
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setSelectedIndex(prev => Math.max(prev - 1, 0));
    } else if (e.key === 'Enter' && results[selectedIndex]) {
      e.preventDefault();
      handleSelect(results[selectedIndex]);
    } else if (e.key === 'Escape') {
      setIsOpen(false);
      setQuery('');
    }
  };

  const handleSelect = (result: SearchResult) => {
    onNavigate(result.url);
    setIsOpen(false);
    setQuery('');
  };

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Keyboard shortcut to focus search
  useEffect(() => {
    const handleGlobalKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault();
        inputRef.current?.focus();
        setIsOpen(true);
      }
    };
    document.addEventListener('keydown', handleGlobalKeyDown);
    return () => document.removeEventListener('keydown', handleGlobalKeyDown);
  }, []);

  return (
    <div ref={containerRef} className="relative">
      <div className="relative">
        <svg
          className={`absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 ${
            theme === 'dark' ? 'text-gray-400' : 'text-gray-500'
          }`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
          />
        </svg>
        <input
          ref={inputRef}
          type="text"
          placeholder={isLoading ? "Loading search..." : "Search docs... (⌘K)"}
          value={query}
          onChange={(e) => {
            setQuery(e.target.value);
            setIsOpen(true);
          }}
          onFocus={() => setIsOpen(true)}
          onKeyDown={handleKeyDown}
          disabled={isLoading}
          className={`w-full pl-10 pr-4 py-2 rounded-lg text-sm transition-all ${
            theme === 'dark'
              ? 'bg-white/10 border border-white/20 text-white placeholder-gray-400 focus:border-cyan-400 focus:ring-1 focus:ring-cyan-400'
              : 'bg-white border border-gray-300 text-gray-900 placeholder-gray-500 focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500'
          } ${isLoading ? 'opacity-50 cursor-wait' : ''}`}
        />
        <kbd
          className={`absolute right-3 top-1/2 -translate-y-1/2 hidden sm:inline-flex items-center px-1.5 py-0.5 text-xs rounded ${
            theme === 'dark'
              ? 'bg-white/10 text-gray-400 border border-white/20'
              : 'bg-gray-100 text-gray-500 border border-gray-300'
          }`}
        >
          ⌘K
        </kbd>
      </div>

      {/* Results dropdown */}
      {isOpen && results.length > 0 && (
        <div
          className={`absolute top-full left-0 right-0 mt-2 rounded-lg shadow-xl overflow-hidden z-50 max-h-[60vh] overflow-y-auto ${
            theme === 'dark'
              ? 'bg-gray-900 border border-white/10'
              : 'bg-white border border-gray-200'
          }`}
        >
          {results.map((result, index) => (
            <button
              key={`${result.url}-${index}`}
              onClick={() => handleSelect(result)}
              className={`w-full text-left px-4 py-3 flex flex-col gap-1 transition-colors ${
                index === selectedIndex
                  ? theme === 'dark'
                    ? 'bg-cyan-500/20'
                    : 'bg-cyan-50'
                  : theme === 'dark'
                  ? 'hover:bg-white/5'
                  : 'hover:bg-gray-50'
              }`}
            >
              <span
                className={`font-medium ${
                  theme === 'dark' ? 'text-white' : 'text-gray-900'
                }`}
              >
                {result.title}
              </span>
              {result.breadcrumbs && result.breadcrumbs.length > 0 && (
                <span
                  className={`text-xs ${
                    theme === 'dark' ? 'text-gray-400' : 'text-gray-500'
                  }`}
                >
                  {result.breadcrumbs.join(' › ')}
                </span>
              )}
            </button>
          ))}
        </div>
      )}

      {/* No results message */}
      {isOpen && query.trim() && results.length === 0 && !isLoading && (
        <div
          className={`absolute top-full left-0 right-0 mt-2 rounded-lg shadow-xl p-4 text-center z-50 ${
            theme === 'dark'
              ? 'bg-gray-900 border border-white/10 text-gray-400'
              : 'bg-white border border-gray-200 text-gray-500'
          }`}
        >
          No results found for "{query}"
        </div>
      )}
    </div>
  );
}
