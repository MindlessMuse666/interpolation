<template>
  <v-app>
    <v-app-bar color="primary">
      <v-app-bar-title>Численные методы: Интерполяция</v-app-bar-title>
      <v-spacer></v-spacer>
      <v-btn to="/learning" variant="text">Обучение</v-btn>
      <v-btn to="/sandbox" variant="text">Песочница</v-btn>
    </v-app-bar>

    <v-main>
      <router-view></router-view>
    </v-main>

    <v-footer app border>
      <span class="text-caption">© 2026 Учебный проект по интерполяции</span>
    </v-footer>
  </v-app>
</template>

<script setup>
</script>

<style>
.v-main {
  background-color: #f5f5f5;
}
</style>
