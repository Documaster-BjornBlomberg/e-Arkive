<template>
  <MainLayout>
    <div class="settings-page">
      <div class="settings-container">
        <h2>User Settings</h2>
        
        <div class="settings-form">
          <div class="setting-item">
            <label>Theme</label>
            <div class="theme-selector">
              <button 
                class="theme-button" 
                :class="{ 'active': currentTheme === 'light' }"
                @click="setTheme('light')"
              >
                Light
              </button>
              <button 
                class="theme-button" 
                :class="{ 'active': currentTheme === 'dark' }"
                @click="setTheme('dark')"
              >
                Dark
              </button>
            </div>
          </div>
          
          <StatusMessage 
            v-if="statusMessage" 
            :type="statusType" 
            :message="statusMessage" 
          />
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup>
import { ref } from 'vue';
import MainLayout from '../templates/MainLayout.vue';
import StatusMessage from '../atoms/StatusMessage.vue';
import { useTheme } from '../../composables/useTheme';

const { currentTheme, saveTheme } = useTheme();
const statusMessage = ref('');
const statusType = ref('');

const setTheme = async (theme) => {
  try {
    await saveTheme(theme);
    statusMessage.value = `Theme changed to ${theme} mode`;
    statusType.value = 'success';
  } catch (error) {
    console.error('Error saving theme:', error);
    statusMessage.value = 'Error saving theme setting';
    statusType.value = 'error';
  }
  
  // Hide message after 3 seconds
  setTimeout(() => {
    statusMessage.value = '';
  }, 3000);
};
</script>

<style scoped>
.settings-page {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.settings-container {
  background-color: var(--surface-color);
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
}

.settings-form {
  margin-top: 20px;
}

.setting-item {
  margin-bottom: 20px;
}

.setting-item label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
}

.theme-selector {
  display: flex;
  gap: 10px;
}

.theme-button {
  background-color: var(--surface-color);
  border: 1px solid var(--border-color);
  color: var(--text-color);
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.theme-button:hover {
  border-color: var(--button-bg);
}

.theme-button.active {
  background-color: var(--button-bg);
  color: var(--button-text);
  border-color: var(--button-bg);
}
</style>