<template>
  <div class="folder-panel" :class="{ 'collapsed': !isOpen }">
    <div class="folder-panel-header">
      <h3 v-if="isOpen">Mappar</h3>
      <button class="toggle-button" @click="togglePanel">
        {{ isOpen ? '◀' : '▶' }}
      </button>
    </div>
    
    <div v-if="isOpen" class="folder-panel-content">
      <div v-if="isLoading" class="loading-state">
        Laddar mappstruktur...
      </div>
      <div v-else-if="error" class="error-state">
        Ett fel uppstod: {{ error }}
      </div>
      <div v-else-if="rootNodes.length === 0" class="empty-state">
        Inga mappar hittades
      </div>
      <div v-else class="folder-tree">
        <FolderNode 
          v-for="node in rootNodes" 
          :key="node.id" 
          :node="node"
          :selected-node-id="selectedNodeId"
          @select-node="onSelectNode"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, provide } from 'vue';
import { useNodeHandling } from '../../composables/useNodeHandling';
import FolderNode from '../molecules/FolderNode.vue';

const props = defineProps({
  modelValue: {
    type: String,
    default: null
  }
});

const emit = defineEmits(['update:modelValue']);

const {
  rootNodes,
  nodeLoading: isLoading,
  nodeError: error,
  isNodePanelOpen,
  fetchRootNodes,
  toggleNodePanel
} = useNodeHandling();

const selectedNodeId = ref(props.modelValue);

// Methods
const onSelectNode = (nodeId) => {
  selectedNodeId.value = nodeId;
  emit('update:modelValue', nodeId);
};

const togglePanel = () => {
  toggleNodePanel();
};

// Computed properties
const isOpen = isNodePanelOpen;

// Initialize component
onMounted(async () => {
  await fetchRootNodes();
});
</script>

<style scoped>
.folder-panel {
  min-width: 250px;
  width: 250px;
  height: 100%;
  background-color: var(--background-color);
  border-right: 1px solid var(--border-color);
  transition: all 0.3s ease;
  overflow: hidden;
  position: relative;
}

.folder-panel.collapsed {
  min-width: 40px;
  width: 40px;
}

.folder-panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border-bottom: 1px solid var(--border-color);
}

.folder-panel.collapsed .folder-panel-header {
  justify-content: center;
  padding: 12px 0;
}

.folder-panel-header h3 {
  margin: 0;
  font-size: 1.1rem;
  color: var(--text-color);
}

.toggle-button {
  background: none;
  border: none;
  color: var(--text-color);
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background-color 0.2s ease;
  z-index: 5;
}

.toggle-button:hover {
  background-color: rgba(128, 128, 128, 0.1);
}

.folder-panel-content {
  padding: 12px;
  height: calc(100% - 60px);
  overflow-y: auto;
}

.loading-state, .error-state, .empty-state {
  padding: 20px 10px;
  color: var(--text-color-light);
  text-align: center;
}

.error-state {
  color: #e74c3c;
}

.folder-tree {
  padding-top: 8px;
}
</style>