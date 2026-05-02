import { defineConfig, type Plugin } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

import { cloudflare } from "@cloudflare/vite-plugin";

function envMetaPlugin(): Plugin {
  const getEnv = () => process.env.VITE_APP_ENV || 'development'
  const getCommit = () => process.env.VITE_APP_COMMIT_HASH || 'unknown'

  const replaceVars = (html: string) =>
    html.replace(/%VITE_APP_ENV%/g, getEnv())
        .replace(/%VITE_APP_COMMIT_HASH%/g, getCommit())

  return {
    name: 'env-meta',
    enforce: 'pre',
    transformIndexHtml: {
      order: 'pre',
      handler(html) {
        return replaceVars(html)
      }
    },
    // Run BEFORE the Cloudflare plugin's middleware (which uses the post-hook pattern).
    // Rewrites SPA routes (e.g. /home) to '/' so Vite serves index.html correctly.
    configureServer(server) {
      server.middlewares.use((req, _res, next) => {
        const url = req.url ?? ''
        if (
          req.method === 'GET' &&
          !url.startsWith('/api') &&
          !url.startsWith('/@') &&
          !url.startsWith('/__') &&
          !url.includes('.')
        ) {
          req.url = '/'
        }
        next()
      })
    },
  }
}

// https://vite.dev/config/
export default defineConfig({
  appType: 'spa',
  plugins: [envMetaPlugin(), react(), tailwindcss(), cloudflare()],
  server: {
    proxy: {
      '/api': {
        target: process.env.API_URL || 'http://localhost:40080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      },
    },
  },
})