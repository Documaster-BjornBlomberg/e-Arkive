<template>
  <div class="settings-page">
    <HeaderSection title="Settings" />
    
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
</template>

<script setup>
import { ref } from 'vue';
import HeaderSection from '../organisms/HeaderSection.vue';
import StatusMessage from '../atoms/StatusMessage.vue';
import { useTheme } from '../../composables/useTheme';

const { currentTheme, saveTheme } = useTheme();
const statusMessage = ref('');
const statusType = ref('');

const setTheme = async (theme) => {
  try {
    const success = await saveTheme(theme);
    
    if (success) {
      statusMessage.value = 'Theme settings saved successfully!';
      statusType.value = 'success';
    } else {
      throw new Error('Failed to save theme');
    }
    
    // Clear status message after 3 seconds
    setTimeout(() => {
      statusMessage.value = '';
    }, 3000);
    
  } catch (error) {
    console.error('Error saving theme setting:', error);
    statusMessage.value = 'Failed to save theme settings.';
    statusType.value = 'error';
  }
};
</script>

<style scoped>
.settings-page {
  padding: 20px;
}

.settings-container {
  max-width: 800px;
  margin: 20px auto;
  padding: 20px;
  background-color: var(--surface-color);
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

h2 {
  margin-top: 0;
  color: var(--text-color);
  font-size: 1.5rem;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 20px;
}

.settings-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.setting-item {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

label {
  font-weight: bold;
  color: var(--text-color);
}

.theme-selector {
  display: flex;
  gap: 10px;
}

.theme-button {
  padding: 8px 16px;
  border: 2px solid var(--border-color);
  border-radius: 4px;
  background-color: var(--background-color);
  color: var(--text-color);
  cursor: pointer;
  transition: all 0.2s ease;
}

.theme-button.active {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.theme-button:hover:not(.active) {
  background-color: var(--hover-color);
}
</style>