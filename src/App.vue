<script setup lang="ts">
import { onMounted } from 'vue';
import { RouterView, useRouter, RouteLocationNormalized } from 'vue-router';
import { useAuth } from './composables/useAuth';
import { useTheme } from './composables/useTheme';

const router = useRouter();
const { isAuthenticated, checkAuth } = useAuth();
const { loadTheme } = useTheme();

// Navigation guard setup
onMounted(async () => {
  // Check authentication on app mount
  await checkAuth();
  
  // Load user theme preference
  loadTheme();
  
  // Add navigation guard
  router.beforeEach((to: RouteLocationNormalized, from: RouteLocationNormalized, next) => {
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
