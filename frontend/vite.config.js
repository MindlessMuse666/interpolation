import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  build: {
    chunkSizeWarningLimit: 1500,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (!id.includes('node_modules')) return
          if (id.includes('/katex/')) return 'katex'
          if (id.includes('/chart.js/') || id.includes('/vue-chartjs/')) return 'chart'
          if (id.includes('/vuetify/') || id.includes('/vue-router/') || id.includes('/vue/')) return 'vue'
          if (id.includes('/lodash/')) return 'lodash'
          return 'vendor'
        }
      }
    }
  }
})
