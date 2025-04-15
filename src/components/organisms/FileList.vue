<template>
  <div class="file-list-container">
    <div class="file-list-header">
      <div class="file-list-search">
        <SearchInput v-model="searchQuery" placeholder="Sök i dokument..." />
      </div>
      <div class="view-mode-selector">
        <button 
          v-for="mode in viewModes" 
          :key="mode.id"
          class="view-mode-button"
          :class="{ active: viewMode === mode.id }"
          @click="$emit('update:viewMode', mode.id)"
          :title="mode.label">
          <span class="material-icons">{{ mode.icon }}</span>
        </button>
      </div>
    </div>
    
    <div class="file-list" :class="[viewMode, { 'transitioning': isTransitioning }]">
      <div v-if="isLoading" class="loading-state">
        Laddar dokumentlista...
      </div>
      <div v-else-if="!filteredFiles.length" class="empty-state">
        Inga dokument hittades.
      </div>
      <template v-else>
        <Transition 
          mode="out-in"
          @before-enter="isTransitioning = true"
          @after-leave="isTransitioning = false">
          <component 
            :is="currentViewComponent"
            :files="filteredFiles"
            :key="viewMode"
            :selected-file-id="selectedFileId"
            :expanded-file-id="expandedFileId"
            :view-mode="viewMode"
            @select="$emit('select-file', $event)"
            @toggle-expand="$emit('toggle-expand', $event)" />
        </Transition>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import SearchInput from '../molecules/SearchInput.vue';
import FileCard from './FileCard.vue';
import FileTableView from '../molecules/FileTableView.vue';
import FileGridView from '../molecules/FileGridView.vue';
import FileCardView from '../molecules/FileCardView.vue';
import FileMillerView from '../molecules/FileMillerView.vue';

const viewModes = [
  { id: 'table', label: 'Tabellvy', icon: 'table_rows' },
  { id: 'grid', label: 'Rutnätsvy', icon: 'grid_view' },
  { id: 'card', label: 'Kortvy', icon: 'view_agenda' },
  { id: 'split', label: 'Delad vy', icon: 'view_sidebar' },
  { id: 'miller', label: 'Mappvy', icon: 'folder_open' }
];

const props = defineProps({
  files: {
    type: Array,
    default: () => []
  },
  selectedFileId: {
    type: String,
    default: null
  },
  expandedFileId: {
    type: String,
    default: null
  },
  isLoading: {
    type: Boolean,
    default: false
  },
  viewMode: {
    type: String,
    default: 'table'
  }
});

defineEmits(['select-file', 'toggle-expand', 'update:viewMode']);

const searchQuery = ref('');
const isTransitioning = ref(false);

const currentViewComponent = computed(() => {
  switch(props.viewMode) {
    case 'table':
      return FileTableView;
    case 'grid':
      return FileGridView;
    case 'miller':
      return FileMillerView;
    case 'split':
    case 'card':
    default:
      return FileCardView;
  }
});

// Filter files based on search query
const filteredFiles = computed(() => {
  if (!searchQuery.value.trim()) return props.files;
  
  const query = searchQuery.value.toLowerCase();
  return props.files.filter(file => 
    file.name.toLowerCase().includes(query) || 
    file.metadata?.some(meta => 
      meta.key.toLowerCase().includes(query) || 
      meta.value.toLowerCase().includes(query)
    )
  );
});
</script>

<style scoped>
.file-list-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.file-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid var(--border-color);
  background-color: var(--surface-color);
}

.file-list-search {
  flex: 1;
  margin-right: 20px;
}

.view-mode-selector {
  display: flex;
  gap: 5px;
}

.view-mode-button {
  background: none;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 6px;
  cursor: pointer;
  color: var(--text-color);
  transition: all 0.2s;
}

.view-mode-button:hover {
  background-color: var(--hover-color);
}

.view-mode-button.active {
  background-color: var(--primary-color-light);
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.file-list {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.file-list.table {
  padding: 0;
}

.file-list.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  grid-gap: 20px;
}

.file-list.card {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.loading-state, .empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
  color: var(--text-color-secondary);
  font-size: 16px;
}

/* View transitions */
.v-enter-active,
.v-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.v-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.v-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.file-list.transitioning {
  pointer-events: none;
}
</style>