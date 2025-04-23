<template>
  <div class="main-layout">
    <HeaderSection @logout="handleLogout" />
    <main class="main-content">
      <slot></slot>
    </main>
  </div>
</template>

<script setup>
import { onMounted } from 'vue';
import { useAuth } from '../../composables/useAuth';
import HeaderSection from '../organisms/HeaderSection.vue';

const { checkAuth, logout } = useAuth();

const handleLogout = () => {
  logout();
};

onMounted(() => {
  checkAuth();
});
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  padding: 20px;
  background-color: var(--background-color);
}
</style>