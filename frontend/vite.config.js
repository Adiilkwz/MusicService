import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // Login/Register and all /api routes go to api-gateway (not streaming-service :8081)
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    }
  }
})