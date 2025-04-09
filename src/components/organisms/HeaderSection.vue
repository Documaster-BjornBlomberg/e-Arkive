<template>
  <header class="header-section">
    <h1>e-Arkive</h1>
    <nav class="nav-links">
      <router-link to="/list" class="nav-link">Dokumentlista</router-link>
      <router-link to="/upload" class="nav-link">Ladda upp</router-link>
    </nav>
    <button @click="toggleTheme" class="theme-toggle-btn">
      {{ isDarkTheme ? 'Ljust Tema' : 'Mörkt Tema' }}
    </button>
  </header>
</template>

<script setup>
import { ref, onMounted } from 'vue';

const isDarkTheme = ref(true); // Default to dark theme

const toggleTheme = () => {
  isDarkTheme.value = !isDarkTheme.value;
  const theme = isDarkTheme.value ? 'dark' : 'light';
  document.documentElement.setAttribute('data-theme', theme);
  
  // Save preference to localStorage
  localStorage.setItem('theme', theme);
};

// Set theme on load, prioritizing saved preference if available
onMounted(() => {
  // Check for saved theme preference
  const savedTheme = localStorage.getItem('theme') || 'dark'; // Default to dark if none saved
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
  background-color: var(--header-bg);
  color: var(--text-color);
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

.theme-toggle-btn {
  background-color: var(--button-bg);
  color: var(--button-text);
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.theme-toggle-btn:hover {
  background-color: var(--button-hover-bg);
}
</style>