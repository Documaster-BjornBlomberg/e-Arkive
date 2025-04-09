<template>
  <div class="file-list-container">
    <div class="file-list-header">
      <h2>Dokumentlista</h2>
      <div class="search-container">
        <SearchInput v-model="searchQuery" placeholder="Sök dokument..." />
      </div>
    </div>
    
    <div class="file-list">
      <FileCard 
        v-for="file in filteredFiles" 
        :key="file.id"
        :file="file"
        :isSelected="selectedFileId === file.id"
        :isExpanded="viewMode === 'expandable' && expandedFileId === file.id"
        :viewMode="viewMode"
        @click="$emit('select-file', file)"
        @toggle-expand="$emit('toggle-expand', file.id)"
      />
      
      <div v-if="filteredFiles.length === 0 && !isLoading" class="empty-state">
        <p>Inga dokument hittades</p>
        <small v-if="searchQuery">Prova att söka med andra sökord</small>
      </div>
      
      <div v-if="isLoading" class="loading-state">
        Laddar dokumentlista...
      </div>
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
}

.file-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.file-list-header h2 {
  margin: 0;
}

.search-container {
  width: 250px;
}

.file-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-y: auto;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px;
  border: 1px dashed #ccc;
  border-radius: 6px;
  color: #666;
}

.empty-state p {
  margin: 0 0 8px 0;
  font-weight: 500;
}

.loading-state {
  display: flex;
  justify-content: center;
  padding: 40px;
  color: #666;
  font-style: italic;
}
</style>