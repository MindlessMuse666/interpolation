import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { aliases, mdi } from 'vuetify/iconsets/mdi'
import '@mdi/font/css/materialdesignicons.css'

const remRamTheme = {
  dark: false,
  colors: {
    background: '#F8F6F9',
    surface: '#FFFFFF',
    primary: '#E88CA5',
    secondary: '#5E9DC8',
    error: '#DC5A5A',
    success: '#3BAB7B',
    warning: '#EA9A4A',
    'on-background': '#1E293B',
    'on-surface': '#1E293B',
    'on-primary': '#FFFFFF',
    'on-secondary': '#FFFFFF',
    'on-error': '#FFFFFF',
    'on-success': '#FFFFFF',
    'on-warning': '#1E293B',
    'border-color': '#E2E8F0',
    'text-primary': '#1E293B',
    'text-secondary': '#64748B'
  },
  variables: {
    'border-radius-root': '8px',
    'font-family-root': '"Inter", sans-serif',
    'font-family-mono': '"JetBrains Mono", monospace'
  }
}

export default createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'remRamTheme',
    themes: {
      remRamTheme
    }
  },
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: {
      mdi,
    },
  },
})
