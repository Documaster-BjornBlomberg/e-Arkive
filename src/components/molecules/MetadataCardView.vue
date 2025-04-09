<template>
  <div class="metadata-cards">
    <div v-if="metadata && metadata.length > 0">
      <div v-for="(meta, index) in metadata" :key="index" class="metadata-card">
        <div class="metadata-card-content">
          <div class="metadata-item">
            <div class="metadata-key">{{ formatKey(meta.key) }}:</div>
            <div class="metadata-value">{{ meta.value }}</div>
            <button 
              class="delete-meta-btn" 
              @click="$emit('delete-metadata', meta.key)"
              title="Ta bort metadata"
            >
              ×
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-else class="no-metadata">
      Ingen metadata har lagts till ännu.
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  metadata: {
    type: Array,
    required: true
  }
});

defineEmits(['delete-metadata']);

// Format key for display
const formatKey = (key) => {
  // Convert camelCase or snake_case to Title Case with spaces
  return key
    .replace(/_/g, ' ')
    .replace(/([A-Z])/g, ' $1')
    .replace(/^\w/, c => c.toUpperCase())
    .trim();
};
</script>

<style scoped>
.metadata-cards {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  max-height: 300px;
  overflow-y: auto;
}

.metadata-card {
  flex: 1 1 calc(50% - 10px);
  min-width: 200px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 15px;
  margin-bottom: 5px;
  background-color: var(--background-color);
  position: relative;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.metadata-card-content {
  display: flex;
  flex-direction: column;
}

.metadata-item {
  position: relative;
}

.metadata-key {
  font-weight: bold;
  margin-bottom: 5px;
  color: var(--text-color);
  padding-right: 20px; /* Space for delete button */
}

.metadata-value {
  word-break: break-word;
  color: var(--text-color);
}

.delete-meta-btn {
  position: absolute;
  right: 0;
  top: 0;
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