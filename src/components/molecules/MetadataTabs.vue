<script setup>
import { defineProps, defineEmits } from 'vue';

const props = defineProps({
  tabs: {
    type: Array,
    required: true
  },
  activeTab: {
    type: String,
    required: true
  }
});

const emit = defineEmits(['tabChange']);

const selectTab = (tabId) => {
  emit('tabChange', tabId);
};
</script>

<template>
  <div class="metadata-tabs">
    <button 
      v-for="tab in tabs" 
      :key="tab.id"
      class="tab-button"
      :class="{ 'active': activeTab === tab.id }"
      @click="selectTab(tab.id)"
    >
      {{ tab.label }}
    </button>
  </div>
</template>

<style scoped>
.metadata-tabs {
  display: flex;
  border-bottom: 1px solid var(--border-color);
}

.tab-button {
  padding: 8px 16px;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 0.9rem;
  border-bottom: 2px solid transparent;
  transition: all 0.2s ease;
}

.tab-button:hover {
  background-color: rgba(0,0,0,0.03);
}

.tab-button.active {
  border-bottom-color: var(--button-bg);
  color: var(--button-bg);
  font-weight: 500;
}
</style>