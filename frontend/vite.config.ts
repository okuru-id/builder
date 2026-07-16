import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'
import http from 'http'

// Multi-page: `/` = landing, `/admin/` = admin SPA.
// ponytail: single dev server (5173), single build output. Backend target
// via VITE_BACKEND_HOST / VITE_BACKEND_PORT (see .env.example); defaults 3000.
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const backendHost = env.VITE_BACKEND_HOST || 'localhost'
  const backendPort = Number(env.VITE_BACKEND_PORT || 3000)
  const backendURL = `http://${backendHost}:${backendPort}`

  return {
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
            const parsed = new URL(req.url ?? '/', `http://${req.headers.host}`)
            if (parsed.pathname !== '/') return next()
            const proxyReq = http.request(
              { host: backendHost, port: backendPort, path: req.url, method: req.method, headers: req.headers },
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
          target: backendURL,
          changeOrigin: true,
        },
        '/api': {
          target: backendURL,
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
  }
})
