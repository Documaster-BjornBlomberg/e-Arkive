import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import './style.css'

import UploadPage from './components/pages/UploadPage.vue'
import ListPage from './components/pages/ListPage.vue'
import LoginPage from './components/pages/LoginPage.vue'

// Define routes
const routes = [
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
  { path: '/', redirect: '/list' },
]

// Create router
const router = createRouter({
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