<script setup>
import FileIcon from '../atoms/FileIcon.vue'

defineProps({
  file: {
    type: Object,
    required: true
  }
})
</script>

<template>
  <div class="file-preview">
    <h2>Selected File</h2>
    
    <div class="file-preview-content">
      <FileIcon :fileName="file.name" :fileType="file.type" />
      
      <div class="file-details">
        <p class="file-name">{{ file.name }}</p>
        <p><strong>Size:</strong> {{ Math.round(file.size / 1024) }} KB</p>
        <p><strong>Type:</strong> {{ file.type }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.file-preview {
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 15px;
  margin-bottom: 20px;
  background-color: rgba(128, 128, 128, 0.05);
}

.file-preview h2 {
  margin-top: 0;
  margin-bottom: 15px;
}

.file-preview-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.file-details {
  flex: 1;
}

.file-name {
  font-weight: bold;
  font-size: 1.1em;
  margin-bottom: 5px;
}
</style>