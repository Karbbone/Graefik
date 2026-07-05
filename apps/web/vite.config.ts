import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

// En conteneur, docker-compose fournit VITE_PROXY_TARGET=http://api:8080.
// En natif, on retombe sur localhost.
const proxyTarget = process.env.VITE_PROXY_TARGET ?? 'http://localhost:8080'

// https://vite.dev/config/ — https://vitest.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    host: true,
    port: 5173,
    strictPort: true,
    // Polling : necessaire pour le HMR sur bind mount Windows/Docker.
    watch: { usePolling: true },
    proxy: {
      '/api': {
        target: proxyTarget,
        changeOrigin: true,
      },
    },
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    css: false,
    coverage: {
      provider: 'v8',
      include: ['src/**/*.{ts,tsx}'],
      exclude: ['src/**/*.test.{ts,tsx}', 'src/test/**', 'src/**/index.ts'],
    },
  },
})
