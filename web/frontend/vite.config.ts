import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { readFileSync } from 'fs'
import { join } from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    {
      name: 'docusaurus-proxy',
      configureServer(server) {
        server.middlewares.use((req, res, next) => {
          // If the request is for /docs/ or any path under /docs/, serve Docusaurus
          if (req.url?.startsWith('/docs/')) {
            // Remove /docs prefix for file lookup
            const filePath = req.url.replace('/docs', '')
            const fullPath = join(__dirname, 'public', 'docs', filePath === '/' ? 'index.html' : filePath)
            
            try {
              // Try to serve the file directly
              const content = readFileSync(fullPath)
              const ext = fullPath.split('.').pop()
              const contentType = 
                ext === 'html' ? 'text/html' :
                ext === 'css' ? 'text/css' :
                ext === 'js' ? 'application/javascript' :
                ext === 'json' ? 'application/json' :
                ext === 'svg' ? 'image/svg+xml' :
                ext === 'png' ? 'image/png' :
                ext === 'jpg' || ext === 'jpeg' ? 'image/jpeg' :
                'text/plain'
              
              res.setHeader('Content-Type', contentType)
              res.end(content)
              return
            } catch (err) {
              // If file doesn't exist, try index.html (for SPA routing)
              if (filePath !== '/index.html') {
                try {
                  const indexContent = readFileSync(join(__dirname, 'public', 'docs', 'index.html'))
                  res.setHeader('Content-Type', 'text/html')
                  res.end(indexContent)
                  return
                } catch (indexErr) {
                  // Fall through to next middleware
                }
              }
            }
          }
          next()
        })
      }
    }
  ],
  server: {
    fs: {
      strict: false,
    },
  },
  publicDir: 'public',
})
