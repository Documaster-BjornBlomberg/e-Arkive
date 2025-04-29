import { createApp } from 'vue'
import { createRouter, createWebHistory, Router, RouteRecordRaw } from 'vue-router'
import App from './App.vue'
import './style.css'

import UploadPage from './components/pages/UploadPage.vue'
import ListPage from './components/pages/ListPage.vue'
import LoginPage from './components/pages/LoginPage.vue'
import SettingsPage from './components/pages/SettingsPage.vue'
import AdminPage from './components/pages/AdminPage.vue'

// Define routes
const routes: RouteRecordRaw[] = [
  { 
    path: '/upload', 
    component: UploadPage,
    meta: { requiresAuth: true }
  },
  { 
    path: '/list', 
    component: ListPage,
    meta: { requiresAuth: true }
  },
  { 
    path: '/login', 
    component: LoginPage,
    meta: { requiresAuth: false }
  },
  { 
    path: '/settings', 
    component: SettingsPage,
    meta: { requiresAuth: true }
  },
  { 
    path: '/admin', 
    component: AdminPage,
    meta: { requiresAuth: true, adminOnly: true }
  },
  { path: '/', redirect: '/list' },
]

// Create router
const router: Router = createRouter({
  history: createWebHistory(),
  routes,
})

// Create and mount the app
const app = createApp(App)

// Provide the router to the app
app.use(router)

// Mount the app
app.mount('#app')

// Export router for use in composables
export { router }
