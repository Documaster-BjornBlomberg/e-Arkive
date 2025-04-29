<template>
  <button 
    class="tab-button" 
    :class="{ active }"
    @click="$emit('click')"
  >
    <span v-if="icon" class="material-icons icon">{{ icon }}</span>
    <slot></slot>
  </button>
</template>

<script setup lang="ts">
defineProps<{
  active?: boolean;
  icon?: string;
}>();

defineEmits(['click']);
</script>

<style scoped>
.tab-button {
  background: none;
  border: none;
  padding: 12px 20px;
  color: var(--text-color-secondary);
  font-size: 15px;
  font-weight: 500;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.tab-button:hover {
  color: var(--text-color);
  background-color: var(--hover-color);
}

.tab-button.active {
  color: var(--primary-color);
  border-bottom-color: var(--primary-color);
  background-color: var(--primary-color-light);
}

.icon {
  font-size: 18px;
}
</style>
