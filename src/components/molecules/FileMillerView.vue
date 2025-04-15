<template>
  <div class="miller-columns">
    <div class="miller-scroll-container" ref="scrollContainer">
      <!-- Root folder column -->
      <div class="miller-column">
        <div class="column-header">
          <h3>Root</h3>
        </div>
        <div class="column-content">
          <div 
            v-for="node in rootNodes" 
            :key="node.id"
            class="miller-item"
            :class="{ 
              selected: selectedNodeId === node.id,
              'has-children': node.children?.length || node.hasChildren
            }"
            @click="selectNode(node)">
            <span class="material-icons folder-icon">folder</span>
            <span class="item-name">{{ node.name }}</span>
            <span 
              v-if="node.children?.length || node.hasChildren"
              class="material-icons nav-icon">
              chevron_right
            </span>
          </div>
        </div>
      </div>

      <!-- Dynamic columns for selected path -->
      <template v-for="(column, index) in activeColumns" :key="column.nodeId">
        <div class="miller-column">
          <div class="column-header">
            <h3>{{ column.name }}</h3>
          </div>
          <div class="column-content">
            <!-- Show child folders -->
            <template v-if="column.children">
              <div 
                v-for="node in column.children" 
                :key="node.id"
                class="miller-item"
                :class="{ 
                  selected: selectedNodeId === node.id,
                  'has-children': node.children?.length || node.hasChildren
                }"
                @click="selectNode(node, index + 1)">
                <span class="material-icons folder-icon">folder</span>
                <span class="item-name">{{ node.name }}</span>
                <span 
                  v-if="node.children?.length || node.hasChildren"
                  class="material-icons nav-icon">
                  chevron_right
                </span>
              </div>
            </template>
          </div>
        </div>

        <!-- Files column when a folder is selected -->
        <div v-if="index === activeColumns.length - 1 && column.files" class="miller-column files-column">
          <div class="column-header">
            <h3>Files</h3>
          </div>
          <div class="column-content">
            <div 
              v-for="file in column.files" 
              :key="file.id"
              class="miller-item file-item"
              :class="{ selected: selectedFileId === file.id }"
              @click="selectFile(file)">
              <FileIcon :fileName="file.name" :fileType="file.contentType" />
              <span class="item-name">{{ file.name }}</span>
              <MetadataBadge 
                :count="file.metadata?.length || 0"
                class="metadata-badge" />
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue';
import FileIcon from '../atoms/FileIcon.vue';
import MetadataBadge from '../atoms/MetadataBadge.vue';

const props = defineProps({
  rootNodes: {
    type: Array,
    required: true
  },
  selectedNodeId: {
    type: String,
    default: null
  },
  selectedFileId: {
    type: String,
    default: null
  }
});

const emit = defineEmits(['select-node', 'select-file']);

const scrollContainer = ref(null);
const activeColumns = ref([]);

// Handle node selection
const selectNode = async (node, columnIndex) => {
  emit('select-node', node);
  
  // Trim columns after the selected one
  if (typeof columnIndex === 'number') {
    activeColumns.value = activeColumns.value.slice(0, columnIndex);
  } else {
    activeColumns.value = [];
  }
  
  // Add new column for selected node
  activeColumns.value.push({
    nodeId: node.id,
    name: node.name,
    children: node.children || [],
    files: node.files || []
  });
  
  // Scroll to show new column
  await scrollToEnd();
};

const selectFile = (file) => {
  emit('select-file', file);
};

// Scroll handling
const scrollToEnd = async () => {
  await nextTick();
  if (scrollContainer.value) {
    scrollContainer.value.scrollLeft = scrollContainer.value.scrollWidth;
  }
};

// Watch for external node selection
watch(() => props.selectedNodeId, async (newId) => {
  if (!newId) {
    activeColumns.value = [];
    return;
  }
  
  // Find node in current columns
  const columnIndex = activeColumns.value.findIndex(col => col.nodeId === newId);
  if (columnIndex === -1) {
    // TODO: Handle case when selected node is not in current path
    // This would require traversing the tree to find the path
    return;
  }
  
  // Trim columns after selected one
  activeColumns.value = activeColumns.value.slice(0, columnIndex + 1);
  await scrollToEnd();
});
</script>

<style scoped>
.miller-columns {
  height: 100%;
  background-color: var(--surface-color);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  overflow: hidden;
}

.miller-scroll-container {
  display: flex;
  overflow-x: auto;
  height: 100%;
  scroll-behavior: smooth;
}

.miller-column {
  min-width: 250px;
  max-width: 250px;
  height: 100%;
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  background-color: var(--surface-color);
}

.miller-column:last-child {
  border-right: none;
}

.column-header {
  padding: 12px;
  border-bottom: 1px solid var(--border-color);
  background-color: var(--header-bg);
}

.column-header h3 {
  margin: 0;
  font-size: 1rem;
  color: var(--text-color);
}

.column-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.miller-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  cursor: pointer;
  color: var(--text-color);
  gap: 8px;
  position: relative;
}

.miller-item:hover {
  background-color: var(--hover-color);
}

.miller-item.selected {
  background-color: var(--selection-color);
}

.folder-icon {
  color: var(--folder-color);
  font-size: 20px;
}

.item-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.nav-icon {
  font-size: 18px;
  color: var(--text-color-secondary);
}

.has-children {
  padding-right: 28px;
}

.file-item {
  padding-right: 12px;
}

.metadata-badge {
  flex-shrink: 0;
}

/* Scrollbar styling */
.miller-scroll-container::-webkit-scrollbar,
.column-content::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.miller-scroll-container::-webkit-scrollbar-track,
.column-content::-webkit-scrollbar-track {
  background: var(--surface-color);
}

.miller-scroll-container::-webkit-scrollbar-thumb,
.column-content::-webkit-scrollbar-thumb {
  background-color: var(--border-color);
  border-radius: 4px;
}

.miller-scroll-container::-webkit-scrollbar-thumb:hover,
.column-content::-webkit-scrollbar-thumb:hover {
  background-color: var(--text-color-secondary);
}
</style>