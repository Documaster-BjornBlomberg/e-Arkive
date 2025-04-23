<template>
  <header class="header-section">
    <h1>{{ title || 'e-Arkive' }}</h1>
    <nav v-if="isAuthenticated" class="nav-links">
      <router-link to="/list" class="nav-link" active-class="active">
        <span class="material-icons">list</span>
        <span>Dokumentlista</span>
      </router-link>
      <router-link to="/upload" class="nav-link" active-class="active">
        <span class="material-icons">upload</span>
        <span>Ladda upp</span>
      </router-link>
      <router-link to="/settings" class="nav-link" active-class="active">
        <span class="material-icons">settings</span>
        <span>Inställningar</span>
      </router-link>
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
  padding: 15px 20px;
  background-color: var(--surface-color);
  color: var(--text-color);
  border-bottom: 1px solid var(--border-color);
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
}

.nav-links {
  display: flex;
  gap: 10px;
  margin: 0 20px;
}

.nav-link {
  text-decoration: none;
  color: var(--text-color);
  font-weight: bold;
  padding: 10px 15px;
  border-radius: 4px 4px 0 0;
  transition: all 0.3s;
  border-bottom: 3px solid transparent;
  display: flex;
  align-items: center;
  gap: 6px;
  position: relative;
  margin: 0 1px;
}

.nav-link:hover {
  background-color: rgba(0, 0, 0, 0.05);
  color: var(--button-bg);
}

.nav-link.active {
  color: var(--button-bg);
  border-bottom-color: var(--button-bg);
  background-color: rgba(0, 0, 0, 0.03);
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