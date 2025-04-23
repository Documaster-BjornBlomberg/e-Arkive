<template>
  <header class="header-section">
    <h1>e-Arkive</h1>
    <nav v-if="isAuthenticated" class="nav-links">
      <router-link to="/list" class="nav-link">Dokumentlista</router-link>
      <router-link to="/upload" class="nav-link">Ladda upp</router-link>
    </nav>
    <div class="header-actions">
      <button v-if="isAuthenticated" @click="handleLogout" class="logout-btn">
        Logga ut
      </button>
      <button @click="toggleTheme" class="theme-toggle-btn">
        {{ isDarkTheme ? 'Ljust Tema' : 'Mörkt Tema' }}
      </button>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useAuth } from '../../composables/useAuth';

const { isAuthenticated, logout } = useAuth();
const isDarkTheme = ref(true);

const toggleTheme = () => {
  isDarkTheme.value = !isDarkTheme.value;
  const theme = isDarkTheme.value ? 'dark' : 'light';
  document.documentElement.setAttribute('data-theme', theme);
  localStorage.setItem('theme', theme);
};

const handleLogout = () => {
  logout();
};

onMounted(() => {
  const savedTheme = localStorage.getItem('theme') || 'dark';
  isDarkTheme.value = savedTheme === 'dark';
  document.documentElement.setAttribute('data-theme', savedTheme);
});
</script>

<style scoped>
.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 20px;
  background-color: var(--surface-color);
  color: var(--text-color);
  border-bottom: 1px solid var(--border-color);
}

.nav-links {
  display: flex;
  gap: 15px;
}

.nav-link {
  text-decoration: none;
  color: var(--text-color);
  font-weight: bold;
  transition: color 0.3s;
}

.nav-link:hover {
  color: var(--button-hover-bg);
}

.header-actions {
  display: flex;
  gap: 10px;
}

.theme-toggle-btn, .logout-btn {
  background-color: var(--button-bg);
  color: var(--button-text);
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.theme-toggle-btn:hover, .logout-btn:hover {
  background-color: var(--button-hover-bg);
}
</style>