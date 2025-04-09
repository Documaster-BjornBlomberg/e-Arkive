<template>
  <div 
    class="file-card"
    :class="{ 'selected': isSelected }"
    @click="$emit('click')"
  >
    <h3 class="file-name">{{ file.name }}</h3>
    <div class="file-info-summary">
      <div class="file-size">{{ formatFileSize(file.size) }}</div>
      <div class="file-type">{{ file.contentType }}</div>
    </div>
    <div class="file-metadata-summary">
      <MetadataBadge :count="metadataCount" />
    </div>
  </div>
</template>

<script setup>
import MetadataBadge from '../atoms/MetadataBadge.vue';
import { computed } from 'vue';

const props = defineProps({
  file: {
    type: Object,
    required: true
  },
  isSelected: {
    type: Boolean,
    default: false
  }
});

defineEmits(['click']);

const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
};

const metadataCount = computed(() => {
  return props.file.metadata ? props.file.metadata.length : 0;
});
</script>

<style scoped>
.file-card {
  padding: 15px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  margin-bottom: 10px;
  cursor: pointer;
  transition: all 0.2s;
  background-color: var(--background-color);
}

.file-card:hover {
  border-color: var(--button-bg);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.file-card.selected {
  border-color: var(--button-bg);
  background-color: rgba(76, 175, 80, 0.05);
}

.file-name {
  margin: 0 0 10px 0;
  font-size: 1.1rem;
  word-break: break-word;
}

.file-info-summary {
  display: flex;
  gap: 15px;
  margin-bottom: 10px;
  font-size: 0.9rem;
}

.file-metadata-summary {
  font-size: 0.85rem;
  color: #666;
}
</style>