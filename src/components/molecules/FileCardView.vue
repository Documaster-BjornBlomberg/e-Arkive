<template>
  <div class="file-card-view" :class="{ 'split-view': viewMode === 'split' }">
    <!-- File list section -->
    <div class="file-list-section">
      <div 
        v-for="file in files" 
        :key="file.id"
        class="file-card"
        :class="{ 
          selected: selectedFileId === file.id,
          expanded: expandedFileId === file.id 
        }"
        @click="$emit('select', file.id)">
        
        <div class="card-header">
          <FileIcon :fileName="file.name" :fileType="file.contentType" />
          <div class="card-title">
            <h3 class="file-name" :title="file.name">{{ file.name }}</h3>
            <div class="file-meta">
              <span>{{ formatFileSize(file.size) }}</span>
              <span>·</span>
              <span>{{ new Date(file.createdAt).toLocaleDateString() }}</span>
            </div>
          </div>
          <button 
            v-if="viewMode !== 'split'"
            class="expand-button"
            @click.stop="$emit('toggle-expand', file.id)"
            :title="expandedFileId === file.id ? 'Minimera' : 'Expandera'">
            <span class="material-icons">
              {{ expandedFileId === file.id ? 'expand_less' : 'expand_more' }}
            </span>
          </button>
        </div>
        
        <!-- Metadata preview -->
        <div v-if="shouldShowMetadata(file)" class="metadata-preview">
          <div v-if="file.metadata?.length" class="metadata-list">
            <div 
              v-for="(meta, index) in file.metadata" 
              :key="index"
              class="metadata-item">
              <span class="metadata-key">{{ meta.key }}:</span>
              <span class="metadata-value">{{ meta.value }}</span>
            </div>
          </div>
          <div v-else class="no-metadata">
            Ingen metadata tillgänglig
          </div>
        </div>
      </div>
    </div>
    
    <!-- Preview section for split view -->
    <div v-if="viewMode === 'split'" class="preview-section">
      <template v-if="selectedFile">
        <div class="preview-header">
          <h2>{{ selectedFile.name }}</h2>
          <div class="preview-actions">
            <Button 
              icon="download"
              variant="info"
              @click="$emit('download', selectedFile.id)">
              Ladda ner
            </Button>
            <Button 
              icon="edit"
              variant="primary"
              @click="$emit('edit', selectedFile.id)">
              Redigera metadata
            </Button>
          </div>
        </div>
        
        <div class="preview-content">
          <div class="preview-details">
            <div class="detail-row">
              <span class="detail-label">Storlek:</span>
              <span>{{ formatFileSize(selectedFile.size) }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Typ:</span>
              <span>{{ selectedFile.contentType }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Skapad:</span>
              <span>{{ new Date(selectedFile.createdAt).toLocaleString() }}</span>
            </div>
          </div>
          
          <div class="preview-metadata">
            <h3>Metadata</h3>
            <div v-if="selectedFile.metadata?.length" class="metadata-list">
              <div 
                v-for="(meta, index) in selectedFile.metadata" 
                :key="index"
                class="metadata-item">
                <span class="metadata-key">{{ meta.key }}:</span>
                <span class="metadata-value">{{ meta.value }}</span>
              </div>
            </div>
            <div v-else class="no-metadata">
              Ingen metadata tillgänglig
            </div>
          </div>
        </div>
      </template>
      <div v-else class="no-selection">
        Välj ett dokument för att se detaljer
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import FileIcon from '../atoms/FileIcon.vue';
import Button from '../atoms/Button.vue';

const props = defineProps({
  files: {
    type: Array,
    required: true
  },
  selectedFileId: {
    type: String,
    default: null
  },
  expandedFileId: {
    type: String,
    default: null
  },
  viewMode: {
    type: String,
    default: 'card'
  }
});

defineEmits(['select', 'toggle-expand', 'download', 'edit']);

const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
};

const selectedFile = computed(() => {
  return props.files.find(f => f.id === props.selectedFileId);
});

const shouldShowMetadata = (file) => {
  return props.viewMode !== 'split' && props.expandedFileId === file.id;
};
</script>

<style scoped>
.file-card-view {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.file-card-view.split-view {
  flex-direction: row;
}

.file-list-section {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.split-view .file-list-section {
  flex: 0 0 400px;
  border-right: 1px solid var(--border-color);
}

.file-card {
  background-color: var(--surface-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.file-card:hover {
  border-color: var(--primary-color);
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.file-card.selected {
  border-color: var(--primary-color);
  background-color: var(--selection-color);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.card-title {
  flex: 1;
  min-width: 0;
}

.file-name {
  margin: 0;
  font-size: 1rem;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-meta {
  font-size: 0.8rem;
  color: var(--text-color-secondary);
  margin-top: 4px;
  display: flex;
  gap: 5px;
}

.expand-button {
  background: none;
  border: none;
  padding: 4px;
  cursor: pointer;
  color: var(--text-color-secondary);
  border-radius: 4px;
}

.expand-button:hover {
  background-color: var(--hover-color);
  color: var(--text-color);
}

.metadata-preview {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid var(--border-color);
}

.metadata-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.metadata-item {
  display: flex;
  gap: 8px;
  font-size: 0.9rem;
}

.metadata-key {
  color: var(--text-color-secondary);
  min-width: 120px;
}

.metadata-value {
  flex: 1;
  word-break: break-word;
}

.no-metadata {
  color: var(--text-color-secondary);
  font-style: italic;
  font-size: 0.9rem;
}

/* Preview Section Styles */
.preview-section {
  flex: 1;
  padding: 20px;
  background-color: var(--surface-color);
  overflow-y: auto;
}

.preview-header {
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid var(--border-color);
}

.preview-header h2 {
  margin: 0 0 15px 0;
  word-break: break-word;
}

.preview-actions {
  display: flex;
  gap: 10px;
}

.preview-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.preview-details {
  background-color: var(--surface-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 15px;
}

.detail-row {
  display: flex;
  gap: 10px;
  margin-bottom: 8px;
}

.detail-row:last-child {
  margin-bottom: 0;
}

.detail-label {
  color: var(--text-color-secondary);
  min-width: 100px;
}

.preview-metadata {
  background-color: var(--surface-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 15px;
}

.preview-metadata h3 {
  margin: 0 0 15px 0;
  font-size: 1.1rem;
}

.no-selection {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--text-color-secondary);
  font-style: italic;
}
</style>