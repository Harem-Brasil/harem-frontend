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
    configureServer(server) {
      server.middlewares.use((req, res, next) => {
        if (req.url === '/' || req.url === '/index.html') {
          const _end = res.end.bind(res)
          let html = ''
          res.end = function(chunk: any, ...args: any[]) {
            if (chunk) html += chunk.toString()
            return _end(Buffer.from(replaceVars(html)), ...args)
          }
        }
        next()
      })
    },
  }
}

// https://vite.dev/config/
export default defineConfig({
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