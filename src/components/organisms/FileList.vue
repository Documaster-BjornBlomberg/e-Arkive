<template>
  <div class="file-list-container">
    <div class="file-list-search">
      <SearchInput v-model="searchQuery" placeholder="Sök i dokument..." />
    </div>
    <div class="file-list">
      <div v-if="isLoading" class="loading-state">
        Laddar dokumentlista...
      </div>
      <div v-else-if="!filteredFiles.length" class="empty-state">
        Inga dokument hittades.
      </div>
      <template v-else>
        <FileCard
          v-for="file in filteredFiles"
          :key="file.id"
          :file="file"
          :is-selected="selectedFileId === file.id"
          :is-expanded="expandedFileId === file.id"
          :view-mode="viewMode"
          @select="$emit('select-file', file.id)"
          @toggle-expand="$emit('toggle-expand', file.id)" />
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import SearchInput from '../molecules/SearchInput.vue';
import FileCard from './FileCard.vue';

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
    default: 'sidepanel'
  }
});

defineEmits(['select-file', 'toggle-expand']);

const searchQuery = ref('');

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

.file-list-search {
  padding: 15px 20px;
  border-bottom: 1px solid var(--border-color);
}

.file-list {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  grid-gap: 20px;
}

.loading-state, .empty-state {
  grid-column: 1 / -1;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
  color: var(--text-color-secondary);
  font-size: 16px;
}

/* List view styles */
.file-list.list-view {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
</style>