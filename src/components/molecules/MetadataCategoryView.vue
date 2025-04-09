<template>
  <div class="metadata-categories">
    <div v-if="metadata && metadata.length > 0">
      <div v-for="(category, categoryName) in categorizedMetadata" :key="categoryName" class="category-section">
        <h4 class="category-name">{{ categoryName }}</h4>
        <div class="category-items">
          <div v-for="(meta, index) in category" :key="index" class="metadata-item">
            <div class="meta-key">{{ meta.key.split('.').pop() }}:</div>
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
      </div>
    </div>
    <div v-else class="no-metadata">
      Ingen metadata har lagts till ännu.
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  metadata: {
    type: Array,
    required: true
  }
});

defineEmits(['delete-metadata']);

// Organize metadata by categories
const categorizedMetadata = computed(() => {
  if (!props.metadata || props.metadata.length === 0) return {};
  
  const categories = {};
  
  props.metadata.forEach(meta => {
    // Extract category from key (assuming keys with dots like 'system.filename')
    // If no dot, use 'General' as category
    const parts = meta.key.split('.');
    const category = parts.length > 1 ? parts[0] : 'General';
    
    if (!categories[category]) {
      categories[category] = [];
    }
    
    categories[category].push(meta);
  });
  
  return categories;
});
</script>

<style scoped>
.metadata-categories {
  max-height: 300px;
  overflow-y: auto;
}

.category-section {
  margin-bottom: 15px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  overflow: hidden;
}

.category-name {
  background-color: rgba(128, 128, 128, 0.1);
  padding: 8px 10px;
  margin: 0;
  font-size: 1rem;
  text-transform: capitalize;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-color);
}

.category-items {
  display: flex;
  flex-direction: column;
}

.metadata-item {
  display: flex;
  border-bottom: 1px solid var(--border-color);
  padding: 8px;
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
  top: 5px;
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