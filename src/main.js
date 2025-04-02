import './style.css';
import { createApp } from 'vue'
import App from './App.vue'
import { createRouter, createWebHistory } from 'vue-router'
import FileUpload from './components/FileUpload.vue'
import DocumentList from './components/DocumentList.vue'

const routes = [
  { path: '/upload', component: FileUpload },
  { path: '/list', component: DocumentList },
  { path: '/', redirect: '/upload' }, // Standardrutt
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

const app = createApp(App)
app.use(router)
app.mount('#app')
