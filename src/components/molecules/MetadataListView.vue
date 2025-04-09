<template>
  <div class="metadata-list">
    <div v-if="metadata && metadata.length > 0" class="list-container">
      <div v-for="(meta, index) in metadata" :key="index" class="metadata-item">
        <div class="meta-key">{{ meta.key }}:</div>
        <div class="meta-value">{{ meta.value }}</div>
        <button 
          class="delete-meta-btn" 
          @click="$emit('delete-metadata', meta.key)"
          title="Ta bort metadata"
        >
          ×
        </button>
      </div>
    </div>
    <div v-else class="no-metadata">
      Ingen metadata har lagts till ännu.
    </div>
  </div>
</template>

<script setup>
defineProps({
  metadata: {
    type: Array,
    required: true
  }
});

defineEmits(['delete-metadata']);
</script>

<style scoped>
.metadata-list {
  max-height: 300px;
  overflow-y: auto;
}

.list-container {
  border: 1px solid var(--border-color);
  border-radius: 4px;
  overflow: hidden;
}

.metadata-item {
  display: flex;
  padding: 8px 12px;
  border-bottom: 1px solid var(--border-color);
  position: relative;
  color: var(--text-color);
}

.metadata-item:last-child {
  border-bottom: none;
}

.metadata-item:hover {
  background-color: rgba(128, 128, 128, 0.1);
}

.meta-key {
  font-weight: bold;
  width: 35%;
  word-break: break-word;
  padding-right: 10px;
  color: var(--text-color);
}

.meta-value {
  width: 65%;
  word-break: break-word;
  padding-right: 20px;
  color: var(--text-color);
}

.delete-meta-btn {
  position: absolute;
  right: 5px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  color: #e74c3c;
  opacity: 0.6;
  padding: 0 5px;
}

.delete-meta-btn:hover {
  opacity: 1;
}

.no-metadata {
  padding: 20px;
  text-align: center;
  color: var(--text-color);
  opacity: 0.7;
}
</style>