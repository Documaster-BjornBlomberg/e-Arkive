<template>
  <div class="file-detail-header">
    <div class="sidepanel-header">
      <h2>{{ file.name }}</h2>
    </div>
    
    <div class="file-details">
      <div class="detail-row">
        <span class="detail-label">Storlek:</span>
        <span class="detail-value">{{ formatFileSize(file.size) }}</span>
      </div>
      <div class="detail-row">
        <span class="detail-label">Typ:</span>
        <span class="detail-value">{{ file.contentType }}</span>
      </div>
      <div class="detail-row">
        <span class="detail-label">Skapad:</span>
        <span class="detail-value">{{ new Date(file.createdAt).toLocaleString() }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  file: {
    type: Object,
    required: true
  }
});

// Formatera filstorlek
const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
};
</script>

<style scoped>
.file-detail-header {
  margin-bottom: 20px;
  color: var(--text-color);
}

.sidepanel-header {
  margin-bottom: 15px;
}

.sidepanel-header h2 {
  margin: 0;
  word-break: break-word;
  color: var(--text-color);
}

.file-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-row {
  display: flex;
  gap: 10px;
}

.detail-label {
  font-weight: bold;
  min-width: 80px;
  color: var(--text-color);
}

.detail-value {
  color: var(--text-color);
  opacity: 0.8;
}
</style>