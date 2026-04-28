import { createRouter, createWebHistory } from 'vue-router'
import LearningView from '../views/LearningView.vue'
import SandboxView from '../views/SandboxView.vue'

const routes = [
  {
    path: '/',
    redirect: '/learning'
  },
  {
    path: '/learning',
    name: 'learning',
    component: LearningView
  },
  {
    path: '/sandbox',
    name: 'sandbox',
    component: SandboxView
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
