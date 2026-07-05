import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// En conteneur, docker-compose fournit VITE_PROXY_TARGET=http://api:8080.
// En natif, on retombe sur localhost.
const proxyTarget = process.env.VITE_PROXY_TARGET ?? 'http://localhost:8080'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
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
})
