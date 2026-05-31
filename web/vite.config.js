import { sveltekit } from '@sveltejs/kit/vite'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [sveltekit()],
  server: {
    port: 2050,
    proxy: {
      '/v1': {
        target: 'http://localhost:5000',
        changeOrigin: true,
      },
      '/admin': {
        target: 'http://localhost:2051',
        changeOrigin: true,
      }
    }
  }
})