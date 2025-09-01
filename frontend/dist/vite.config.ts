import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/api': {
        target: 'https://asynx.internal.asynclab.club/',
        changeOrigin: true,
        secure: false
      }
    }
  }
})
