<template>
  <div class="sidepanel" :class="{ 'open': isOpen }">
    <div class="sidepanel-content">
      <div class="sidepanel-close">
        <button class="close-button" @click="$emit('close')" title="Stäng">×</button>
      </div>
      
      <template v-if="file">
        <!-- File Detail Header -->
        <FileDetailHeader :file="file" />
        
        <!-- Action Buttons -->
        <FileActionButtons 
          :isEditing="isEditing" 
          @download="$emit('download', file)"
          @edit="$emit('edit-metadata', file)"
          @save="$emit('save-metadata')"
          @cancel="$emit('cancel-edit')"
          @delete="$emit('delete-file', file.id)"
        />
        
        <!-- Metadata Display when not editing -->
        <div v-if="!isEditing" class="metadata-container">
          <MetadataDisplay 
            :metadata="file.metadata" 
            :activeTab="activeMetadataTabLocal"
            :metadataSearch="metadataSearchLocal"
            @update:activeTab="updateActiveTab"
            @update:metadataSearch="updateMetadataSearch"
            @delete-metadata="$emit('delete-metadata', file.id, $event)"
          />
        </div>
        
        <!-- Metadata Editor when editing -->
        <div v-else class="metadata-container">
          <MetadataEditor 
            :modelValue="editingMetadata" 
            @update:modelValue="$emit('update:editingMetadata', $event)" 
            :disabled="false" 
          />
        </div>
      </template>
      
      <div v-else class="loading-placeholder">
        Laddar fil...
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';
import FileDetailHeader from '../molecules/FileDetailHeader.vue';
import FileActionButtons from '../molecules/FileActionButtons.vue';
import MetadataDisplay from './MetadataDisplay.vue';
import MetadataEditor from '../molecules/MetadataEditor.vue';

const props = defineProps({
  file: {
    type: Object,
    default: null
  },
  isOpen: {
    type: Boolean,
    default: false
  },
  isEditing: {
    type: Boolean,
    default: false
  },
  editingMetadata: {
    type: Array,
    default: () => []
  },
  activeMetadataTab: {
    type: String,
    default: 'list'
  },
  metadataSearch: {
    type: String,
    default: ''
  }
});

const emit = defineEmits([
  'close', 
  'download', 
  'edit-metadata', 
  'save-metadata', 
  'cancel-edit', 
  'delete-file',
  'update:activeTab',
  'update:metadataSearch',
  'delete-metadata'
]);

// Local state for metadata display
const activeMetadataTabLocal = ref(props.activeMetadataTab);
const metadataSearchLocal = ref(props.metadataSearch);

// Update handlers
const updateActiveTab = (tabId) => {
  activeMetadataTabLocal.value = tabId;
  emit('update:activeTab', tabId);
};

const updateMetadataSearch = (searchText) => {
  metadataSearchLocal.value = searchText;
  emit('update:metadataSearch', searchText);
};

// Reset local state when file changes
watch(() => props.file?.id, () => {
  activeMetadataTabLocal.value = 'list'; // Reset to default tab when file changes
  metadataSearchLocal.value = '';
});

// Update local state when props change
watch(() => props.activeMetadataTab, (newVal) => {
  activeMetadataTabLocal.value = newVal;
});

watch(() => props.metadataSearch, (newVal) => {
  metadataSearchLocal.value = newVal;
});
</script>

<style scoped>
.sidepanel {
  position: fixed;
  top: 0;
  right: -500px;
  width: 450px;
  height: 100vh;
  background-color: var(--background-color);
  color: var(--text-color);
  box-shadow: -2px 0 10px rgba(0, 0, 0, 0.2);
  z-index: 100;
  transition: right 0.3s ease, background-color 0.3s, color 0.3s;
  overflow-y: auto;
  border-left: 1px solid var(--border-color);
}

.sidepanel.open {
  right: 0;
}

.sidepanel-content {
  padding: 20px;
}

.sidepanel-close {
  text-align: right;
  margin-bottom: 15px;
}

.close-button {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  padding: 5px;
  color: var(--text-color);
}

.close-button:hover {
  color: var(--button-hover-bg);
}

.metadata-container {
  margin-top: 20px;
}

.loading-placeholder {
  display: flex;
  height: 300px;
  align-items: center;
  justify-content: center;
  color: var(--text-color);
  opacity: 0.7;
  font-style: italic;
}
</style>