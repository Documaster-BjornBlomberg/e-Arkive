<template>
  <header class="header-section">
    <h1>{{ title || 'e-Arkive' }}</h1>
    <nav v-if="isAuthenticated" class="nav-links">
      <router-link to="/list" class="nav-link">Dokumentlista</router-link>
      <router-link to="/upload" class="nav-link">Ladda upp</router-link>
      <router-link to="/settings" class="nav-link">Inställningar</router-link>
    </nav>
    <div class="header-actions">
      <button class="theme-toggle-btn" @click="toggleTheme" :title="currentTheme === 'light' ? 'Växla till mörkt läge' : 'Växla till ljust läge'">
        <span v-if="currentTheme === 'light'" class="material-icons">dark_mode</span>
        <span v-else class="material-icons">light_mode</span>
      </button>
      <button v-if="isAuthenticated" @click="handleLogout" class="logout-btn">
        Logga ut
      </button>
    </div>
  </header>
</template>

<script setup>
import { useAuth } from '../../composables/useAuth';
import { useTheme } from '../../composables/useTheme';

defineProps({
  title: {
    type: String,
    default: null
  }
});

const { isAuthenticated } = useAuth();
const { currentTheme, toggleTheme } = useTheme();

const emit = defineEmits(['logout']);

const handleLogout = () => {
  emit('logout');
};
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
  align-items: center;
}

.logout-btn {
  background-color: var(--button-bg);
  color: var(--button-text);
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.logout-btn:hover {
  background-color: var(--button-hover-bg);
}

.theme-toggle-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--text-color);
  border-radius: 50%;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.3s;
}

.theme-toggle-btn:hover {
  background-color: var(--hover-color);
}

.material-icons {
  font-size: 22px;
}
</style>