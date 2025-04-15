<template>
  <div class="view-selector">
    <div class="selector-buttons">
      <Button 
        v-for="mode in viewModes"
        :key="mode.id"
        :variant="modelValue === mode.id ? 'primary' : 'default'"
        :title="mode.label + (mode.shortcut ? ` (${mode.shortcut})` : '')"
        @click="$emit('update:modelValue', mode.id)"
        class="selector-button">
        <span class="material-icons">{{ mode.icon }}</span>
        <span class="button-label">{{ mode.label }}</span>
      </Button>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted } from 'vue';
import Button from '../atoms/Button.vue';

const props = defineProps({
  modelValue: {
    type: String,
    required: true
  }
});

const emit = defineEmits(['update:modelValue']);

const viewModes = [
  { id: 'table', label: 'Tabellvy', icon: 'table_rows', shortcut: 'Alt+1' },
  { id: 'grid', label: 'Rutnätsvy', icon: 'grid_view', shortcut: 'Alt+2' },
  { id: 'card', label: 'Kortvy', icon: 'view_agenda', shortcut: 'Alt+3' },
  { id: 'split', label: 'Delad vy', icon: 'view_sidebar', shortcut: 'Alt+4' },
  { id: 'miller', label: 'Mappvy', icon: 'folder_open', shortcut: 'Alt+5' }
];

// Handle keyboard shortcuts
const handleKeydown = (event) => {
  if (event.altKey && event.key >= '1' && event.key <= '5') {
    const index = parseInt(event.key) - 1;
    if (index < viewModes.length) {
      event.preventDefault();
      emit('update:modelValue', viewModes[index].id);
    }
  }
};

onMounted(() => {
  window.addEventListener('keydown', handleKeydown);
});

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown);
});
</script>

<style scoped>
.view-selector {
  display: flex;
  align-items: center;
  gap: 10px;
}

.selector-buttons {
  display: flex;
  gap: 5px;
  background-color: var(--surface-color);
  padding: 4px;
  border-radius: 6px;
  border: 1px solid var(--border-color);
}

.selector-button {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 4px;
}

.material-icons {
  font-size: 18px;
}

.button-label {
  font-size: 0.9rem;
}

@media (max-width: 768px) {
  .button-label {
    display: none;
  }
  
  .selector-button {
    padding: 6px;
  }
}
</style>