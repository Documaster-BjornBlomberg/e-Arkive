import { createApp } from 'vue';
import App from './App.vue';
import { createRouter, createWebHistory } from 'vue-router';
import './style.css';

// Initialize theme before app mounts
const initializeTheme = () => {
  const savedTheme = localStorage.getItem('theme') || 'dark'; // Default to dark theme
  document.documentElement.setAttribute('data-theme', savedTheme);
};

// Apply theme immediately
initializeTheme();

const routes = [
  { path: '/upload', component: () => import('./components/pages/UploadPage.vue') },
  { path: '/list', component: () => import('./components/pages/ListPage.vue') },
  { path: '/', redirect: '/list' },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

const app = createApp(App);
app.use(router);
app.mount('#app');