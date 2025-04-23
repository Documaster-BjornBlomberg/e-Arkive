<script setup>
import { onMounted } from 'vue';
import { RouterView, useRouter } from 'vue-router';
import { useAuth } from './composables/useAuth';

const router = useRouter();
const { isAuthenticated, checkAuth } = useAuth();

// Navigation guard setup
onMounted(() => {
  // Check authentication on app mount
  checkAuth();
  
  // Add navigation guard
  router.beforeEach((to, from, next) => {
    if (to.meta.requiresAuth && !isAuthenticated.value) {
      next('/login');
    } else if (to.path === '/login' && isAuthenticated.value) {
      next('/list');
    } else {
      next();
    }
  });
});
</script>

<template>
  <RouterView />
</template>

<style>
:root {
  --background-color: #ffffff;
  --surface-color: #ffffff;
  --text-color: #333333;
  --text-color-secondary: #666666;
  --border-color: #e0e0e0;
  --button-bg: #4CAF50;
  --button-hover-bg: #45a049;
  --button-text: #ffffff;
  --error-color: #e74c3c;
  --success-color: #2ecc71;
}

[data-theme="dark"] {
  --background-color: #1a1a1a;
  --surface-color: #2d2d2d;
  --text-color: #ffffff;
  --text-color-secondary: #b3b3b3;
  --border-color: #404040;
  --button-bg: #45a049;
  --button-hover-bg: #4CAF50;
}

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen,
    Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  background-color: var(--background-color);
  color: var(--text-color);
}
</style>
