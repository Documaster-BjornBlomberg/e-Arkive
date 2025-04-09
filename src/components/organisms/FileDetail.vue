<template>
  <div class="file-detail" v-if="file">
    <div class="detail-header">
      <h2>{{ file.name }}</h2>
    </div>
    
    <div class="file-properties">
      <div class="property-row">
        <span class="property-label">Storlek:</span>
        <span class="property-value">{{ formatFileSize(file.size) }}</span>
      </div>
      <div class="property-row">
        <span class="property-label">Typ:</span>
        <span class="property-value">{{ file.contentType }}</span>
      </div>
      <div class="property-row">
        <span class="property-label">Skapad:</span>
        <span class="property-value">{{ new Date(file.createdAt).toLocaleString() }}</span>
      </div>
    </div>
    
    <div class="action-buttons">
      <Button 
        variant="download" 
        @click="$emit('download-file', file.id)"
        title="Ladda ner">
        Ladda ner
      </Button>
      <Button 
        variant="edit" 
        @click="$emit('edit-metadata', file)"
        title="Redigera metadata">
        Redigera metadata
      </Button>
      <Button 
        variant="danger" 
        @click="$emit('delete-file', file.id)"
        title="Ta bort fil">
        Ta bort
      </Button>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue';
import Button from '../atoms/Button.vue';

const props = defineProps({
  file: {
    type: Object,
    required: true
  }
});

defineEmits(['download-file', 'edit-metadata', 'delete-file']);

// Formatera filstorlek
const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
};
</script>

<style scoped>
.file-detail {
  padding: 15px;
}

.detail-header {
  margin-bottom: 20px;
}

.detail-header h2 {
  margin: 0;
  word-break: break-word;
}

.file-properties {
  margin-bottom: 25px;
}

.property-row {
  display: flex;
  margin-bottom: 10px;
}

.property-label {
  width: 80px;
  font-weight: 500;
  color: #555;
}

.property-value {
  flex: 1;
}

.action-buttons {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.action-buttons button {
  flex: 1;
  min-width: 120px;
}
</style>