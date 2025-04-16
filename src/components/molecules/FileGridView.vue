<template>
  <div class="file-grid">
    <div 
      v-for="file in files" 
      :key="file.id"
      class="grid-item"
      :class="{ selected: selectedFileId === file.id }"
      @click="$emit('select', file.id)">
      
      <div class="grid-item-content">
        <div class="file-icon">
          <FileIcon :fileName="file.name" :fileType="file.contentType" />
        </div>
        <div class="file-info">
          <h3 class="file-name">{{ file.name }}</h3>
          <div class="file-details">
            <span class="file-size">{{ formatFileSize(file.size) }}</span>
            <span class="file-type">{{ file.contentType }}</span>
          </div>
          <MetadataBadge :count="file.metadata?.length || 0" />
        </div>
      </div>

      <!-- Hover preview -->
      <div class="hover-preview">
        <div class="preview-content">
          <div class="preview-header">
            <h4>{{ file.name }}</h4>
          </div>

          <div class="preview-metadata">
            <h5>Metadata</h5>
            <div v-if="file.metadata?.length" class="metadata-preview">
              <div 
                v-for="(meta, index) in file.metadata.slice(0, 3)" 
                :key="index"
                class="metadata-item">
                <span class="metadata-key">{{ meta.key }}:</span>
                <span class="metadata-value">{{ meta.value }}</span>
              </div>
              <div v-if="file.metadata.length > 3" class="metadata-more">
                +{{ file.metadata.length - 3 }} fler fält
              </div>
            </div>
            <div v-else class="no-metadata">
              Ingen metadata tillgänglig
            </div>
          </div>
        </div>

        <div class="preview-actions">
          <button 
            class="preview-action"
            title="Ladda ner"
            @click.stop="$emit('download', file)">
            <span class="material-icons">download</span>
            <span class="action-text">Ladda ner</span>
          </button>
          <button 
            class="preview-action"
            title="Redigera metadata"
            @click.stop="$emit('edit', file)">
            <span class="material-icons">edit</span>
            <span class="action-text">Redigera</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import FileIcon from '../atoms/FileIcon.vue';
import MetadataBadge from '../atoms/MetadataBadge.vue';

const props = defineProps({
  files: {
    type: Array,
    required: true
  },
  selectedFileId: {
    type: String,
    default: null
  }
});

defineEmits(['select', 'download', 'edit']);

const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
};
</script>

<style scoped>
.file-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
  padding: 12px;
}

.grid-item {
  position: relative;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  overflow: hidden;
  background-color: var(--surface-color);
  cursor: pointer;
  transition: all 0.2s ease;
  min-height: 200px;
  max-height: 200px;
  min-width: 200px;
  max-width: 220px;
  display: flex;
  flex-direction: column;
}

.grid-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: var(--primary-color);
}

.grid-item.selected {
  border-color: var(--primary-color);
  background-color: var(--primary-color-light);
}

.grid-item-content {
  padding: 7px;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.file-icon {
  display: flex;
  justify-content: center;
  margin-bottom: 4px;
}

.file-info {
  text-align: center;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.file-name {
  margin: 0;
  font-size: 0.85rem;
  word-break: break-word;
  color: var(--text-color);
  line-height: 1.1;
  max-height: 2.2em;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
}

.file-details {
  font-size: 0.7rem;
  color: var(--text-color-secondary);
  margin: 2px 0;
}

.file-size, .file-type {
  display: inline-block;
  margin: 0 4px;
}

/* Hover preview styles */
.hover-preview {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: var(--surface-color);
  padding: 15px;
  opacity: 0;
  visibility: hidden;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
}

.grid-item:hover .hover-preview {
  opacity: 1;
  visibility: visible;
}

.preview-content {
  flex: 1;
  overflow-y: auto;
  padding-bottom: 50px; /* Minska padding för att ge mer utrymme för knappar */
}

.preview-header {
  margin-bottom: 10px;
}

.preview-header h4 {
  margin: 0;
  font-size: 0.95rem;
  color: var(--text-color);
  word-break: break-word;
}

.preview-actions {
  display: flex;
  gap: 8px;
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 8px;
  background: var(--surface-color);
  border-top: 1px solid var(--border-color);
  z-index: 2;
}

.preview-action {
  flex: 1;
  min-height: 32px;
  background-color: var(--surface-color);
  border: 1px solid var(--border-color);
  padding: 4px 8px;
  cursor: pointer;
  color: var(--text-color);
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.preview-action:hover {
  background-color: var(--hover-color);
  border-color: var(--primary-color);
}

.preview-action .material-icons {
  font-size: 16px;
  margin-right: 4px;
}

.action-text {
  font-size: 0.75rem; /* Minska textstorleken något */
  font-weight: 500;
}

.preview-metadata {
  margin-top: 10px;
}

.preview-metadata h5 {
  margin: 0 0 10px 0;
  font-size: 0.9rem;
  color: var(--text-color);
}

.metadata-preview {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.metadata-item {
  font-size: 0.8rem;
  display: flex;
  gap: 6px;
}

.metadata-key {
  color: var(--text-color-secondary);
  flex-shrink: 0;
}

.metadata-value {
  color: var(--text-color);
  word-break: break-word;
}

.metadata-more {
  font-size: 0.8rem;
  color: var(--text-color-secondary);
  font-style: italic;
  text-align: center;
  padding-top: 4px;
}

.no-metadata {
  font-size: 0.8rem;
  color: var(--text-color-secondary);
  font-style: italic;
  text-align: center;
}
</style>