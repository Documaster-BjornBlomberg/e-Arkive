<template>
  <div class="main-layout">
    <HeaderSection @logout="handleLogout" />
    <main class="main-content">
      <slot></slot>
    </main>
    <div v-if="themeLoading" class="theme-loading-overlay">
      <div class="theme-loading-spinner"></div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue';
import { useAuth } from '../../composables/useAuth';
import { useTheme } from '../../composables/useTheme';
import HeaderSection from '../organisms/HeaderSection.vue';

const { checkAuth, logout } = useAuth();
const { loadTheme, isLoading: themeLoading } = useTheme();

const handleLogout = async () => {
  await logout();
};

onMounted(async () => {
  await checkAuth();
  loadTheme();
});
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  position: relative;
}

.main-content {
  flex: 1;
  padding: 20px;
  background-color: var(--background-color);
}

.theme-loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(255, 255, 255, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.theme-loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(0, 0, 0, 0.1);
  border-left-color: var(--button-bg);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>