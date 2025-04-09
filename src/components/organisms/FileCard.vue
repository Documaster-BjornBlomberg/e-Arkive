<template>
  <div 
    class="file-card"
    :class="{ 'selected': isSelected, 'expanded': isExpanded }"
  >
    <div class="file-header" @click="$emit('click', file)">
      <FileIcon :fileName="file.name" :fileType="file.contentType" />
      <h3 class="file-name">{{ file.name }}</h3>
      
      <!-- Add expand button for expandable view mode -->
      <button
        v-if="expandable"
        class="expand-button"
        @click.stop="$emit('toggle-expand', file.id)"
        :title="isExpanded ? 'Minimera' : 'Expandera'"
      >
        {{ isExpanded ? '▲' : '▼' }}
      </button>
    </div>
    
    <div class="file-info-summary" @click="$emit('click', file)">
      <div class="file-size">{{ formatFileSize(file.size) }}</div>
      <div class="file-type">{{ file.contentType }}</div>
    </div>
    
    <div class="file-metadata-summary" @click="$emit('click', file)">
      <span class="metadata-count">{{ file.metadata ? file.metadata.length : 0 }} metadatafält</span>
    </div>
    
    <!-- Expanded content when in expanded mode -->
    <div v-if="isExpanded" class="expanded-content">
      <hr class="divider" />
      <div class="expanded-section">
        <h4>Metadata</h4>
        <div v-if="file.metadata && file.metadata.length > 0" class="metadata-list">
          <div v-for="(meta, index) in file.metadata" :key="index" class="metadata-item">
            <strong>{{ meta.key }}</strong>: {{ meta.value }}
          </div>
        </div>
        <div v-else class="no-metadata">
          Ingen metadata tillgänglig
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import FileIcon from '../atoms/FileIcon.vue';

const props = defineProps({
  file: {
    type: Object,
    required: true
  },
  isSelected: {
    type: Boolean,
    default: false
  },
  isExpanded: {
    type: Boolean,
    default: false
  },
  viewMode: {
    type: String,
    default: 'sidepanel'
  }
});

defineEmits(['click', 'toggle-expand']);

// Compute whether this card should be expandable based on view mode
const expandable = computed(() => props.viewMode === 'expandable');

// Formatera filstorlek
const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
}
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

.file-card.expanded {
  transform: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border-color: var(--button-bg);
  transition: all 0.3s;
}

.file-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
  position: relative;
}

.file-name {
  margin: 0;
  font-size: 1.1rem;
  word-break: break-word;
  flex: 1;
}

.expand-button {
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 5px 8px;
  color: var(--text-color);
  opacity: 0.7;
  margin-left: auto;
  border-radius: 50%;
}

.expand-button:hover {
  background-color: rgba(128, 128, 128, 0.1);
  opacity: 1;
}

.file-info-summary {
  display: flex;
  gap: 15px;
  margin-bottom: 10px;
  font-size: 0.9rem;
}

.file-metadata-summary {
  font-size: 0.85rem;
  color: var(--text-color);
  opacity: 0.7;
}

.divider {
  margin: 15px 0;
  border: none;
  border-top: 1px solid var(--border-color);
}

.expanded-content {
  margin-top: 10px;
}

.expanded-section h4 {
  margin-top: 0;
  margin-bottom: 10px;
  color: var(--text-color);
}

.metadata-list {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.metadata-item {
  font-size: 0.9rem;
  padding: 5px 0;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-color);
}

.metadata-item:last-child {
  border-bottom: none;
}

.no-metadata {
  font-style: italic;
  color: var(--text-color);
  opacity: 0.6;
  font-size: 0.9rem;
}
</style>