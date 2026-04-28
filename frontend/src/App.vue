<template>
  <v-app>
    <v-app-bar color="background" elevation="0" border="b" class="glass-header">
      <v-container class="d-flex align-center pa-0 px-4 fill-height" fluid>
        <v-avatar size="32" class="mr-3">
          <v-img src="/favicon.svg" alt="Logo" />
        </v-avatar>
        <v-app-bar-title class="font-weight-bold text-h6">Interpolation</v-app-bar-title>
        
        <v-spacer></v-spacer>
        
        <div class="d-none d-sm-flex gap-2">
          <v-btn 
            to="/learning" 
            variant="text" 
            class="nav-btn"
            :active-class="'nav-btn-active'"
          >Обучение</v-btn>
          <v-btn 
            to="/sandbox" 
            variant="text" 
            class="nav-btn"
            :active-class="'nav-btn-active'"
          >Песочница</v-btn>
        </div>

        <!-- Mobile Menu -->
        <v-menu class="d-flex d-sm-none">
          <template v-slot:activator="{ props }">
            <v-btn icon="mdi-menu" variant="text" v-bind="props" class="d-sm-none"></v-btn>
          </template>
          <v-list>
            <v-list-item to="/learning" prepend-icon="mdi-school">
              <v-list-item-title>Обучение</v-list-item-title>
            </v-list-item>
            <v-list-item to="/sandbox" prepend-icon="mdi-test-tube">
              <v-list-item-title>Песочница</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-container>
    </v-app-bar>

    <v-main>
      <router-view></router-view>
    </v-main>

    <v-footer app border height="48" color="surface">
      <v-container class="pa-0 px-4 fill-height d-flex justify-space-between" fluid>
        <span class="text-caption text-text-secondary">© 2026 Учебный проект по интерполяции</span>
        <span class="text-caption text-text-secondary d-none d-sm-inline">Численные методы</span>
      </v-container>
    </v-footer>
  </v-app>
</template>

<script setup>
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&family=JetBrains+Mono:ital,wght@0,100..800;1,100..800&display=swap');

:root {
  --color-ram-pink: #E88CA5;
  --color-rem-blue: #5E9DC8;
  --color-ram-pink-light: #F5C6D3;
  --color-rem-blue-light: #C3E0F2;
  --shadow-card: 0 4px 16px rgba(30, 41, 59, 0.07), 0 1px 3px rgba(30, 41, 59, 0.04);
  --font-family: 'Inter', sans-serif;
  --font-mono: 'JetBrains Mono', monospace;
}

body {
  font-family: var(--font-family);
  background-color: #F8F6F9;
}

.v-application {
  font-family: var(--font-family) !important;
}

.math-font {
  font-family: var(--font-mono) !important;
}

.glass-header {
  backdrop-filter: blur(12px);
  background-color: rgba(248, 246, 249, 0.85) !important;
}

.nav-btn {
  font-weight: 600 !important;
  text-transform: none !important;
  letter-spacing: 0.02em !important;
  border-radius: 8px !important;
}

.nav-btn-active {
  color: var(--color-ram-pink) !important;
  position: relative;
}

.nav-btn-active::after {
  content: '';
  position: absolute;
  bottom: 4px;
  left: 8px;
  right: 8px;
  height: 2px;
  background-color: var(--color-ram-pink);
  border-radius: 2px;
}

.v-card {
  border-radius: 12px !important;
  box-shadow: var(--shadow-card) !important;
  border: 1px solid #E2E8F0 !important;
}

/* Скрытие стрелок у числовых инпутов */
input::-webkit-outer-spin-button,
input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
input[type=number] {
  -moz-appearance: textfield;
}

.gap-2 {
  gap: 8px;
}
</style>
