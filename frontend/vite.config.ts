import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'
import http from 'http'

// Multi-page: `/` = landing, `/admin/` = admin SPA.
// ponytail: single dev server (5173), single build output.
export default defineConfig({
  base: '/',
  plugins: [
    vue(),
    tailwindcss(),
    {
      name: 'rewrite-admin-to-admin-html',
      configureServer(server) {
        // Serve admin SPA shell for any /admin/* path in dev.
        // Skip API routes so the Vite proxy can forward them to the backend.
        server.middlewares.use((req, _res, next) => {
          const url = req.url ?? ''
          if (url.startsWith('/admin/api/')) return next()
          if (url === '/admin' || url.startsWith('/admin/')) {
            req.url = '/admin.html'
          }
          next()
        })
      },
    },
    {
      // ponytail: proxy root '/' to backend so dev matches prod routing
      // (landing_mode custom/empty-fallback logic lives in web.go, not here).
      name: 'proxy-root-to-backend',
      configureServer(server) {
        server.middlewares.use((req, res, next) => {
          if (req.url !== '/') return next()
          const proxyReq = http.request(
            { host: 'localhost', port: 3000, path: '/', method: req.method, headers: req.headers },
            (proxyRes) => {
              res.writeHead(proxyRes.statusCode ?? 200, proxyRes.headers)
              proxyRes.pipe(res)
            },
          )
          proxyReq.on('error', () => next())
          req.pipe(proxyReq)
        })
      },
    },
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
    // ponytail: force single module instance for dnd-kit ecosystem
    dedupe: ['@dnd-kit/abstract', '@dnd-kit/dom', '@dnd-kit/collision', '@dnd-kit/geometry'],
  },
  server: {
    proxy: {
      '/admin/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
      },
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    rollupOptions: {
      input: {
        landing: path.resolve(__dirname, 'index.html'),
        admin: path.resolve(__dirname, 'admin.html'),
      },
    },
  },
})
