<template>
  <div class="folder-panel">
    <div class="folder-panel-header">
      <div class="folder-panel-title">
        <span class="material-icons folder-panel-icon">folder</span>
        <h3>Mappar</h3>
      </div>
      <div class="folder-panel-actions">
        <span class="material-icons document-list-icon" title="Document List">description</span>
      </div>
    </div>
    <div class="folder-panel-content">
      <div v-if="isLoading" class="loading-state">
        <span>Laddar mappar...</span>
      </div>
      <div v-else-if="!nodes || nodes.length === 0" class="empty-state">
        <span>Inga mappar hittades</span>
      </div>
      <div v-else class="node-tree">
        <!-- Root Folder Node Component with Recursive Children -->
        <NodeTree 
          :nodes="nodes" 
          :expanded-node-ids="expandedNodeIds"
          :selected-node-id="selectedNodeId"
          @toggle-expand="handleToggleNode" 
          @select-node="handleSelectNode"
          @contextmenu="handleContextMenu" />
      </div>
    </div>

    <!-- Context Menu -->
    <div v-if="contextMenu.show" 
         class="context-menu" 
         :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px' }">
      <div class="context-menu-item" @click="handleAddNode">
        <span class="material-icons context-menu-icon">add</span>
        <span>Ny Mapp</span>
      </div>
      <div class="context-menu-item" @click="handleRenameNode" v-if="canRename">
        <span class="material-icons context-menu-icon">edit</span>
        <span>Byt namn</span>
      </div>
      <div class="context-menu-item" @click="handleDeleteNode" v-if="canDelete">
        <span class="material-icons context-menu-icon">delete</span>
        <span>Ta bort</span>
      </div>
    </div>

    <!-- Add node dialog -->
    <div v-if="showAddNodeDialog" class="dialog-backdrop" @click="cancelAddNode">
      <div class="dialog" @click.stop>
        <h3>Lägg till ny mapp</h3>
        <div class="dialog-content">
          <label for="nodeName">Mappnamn:</label>
          <input 
            id="nodeName" 
            v-model="newNodeName" 
            type="text" 
            class="input-field" 
            @keyup.enter="submitAddNode"
            ref="nodeNameInput"
            placeholder="Ange mappnamn..." />
        </div>
        <div class="dialog-actions">
          <Button @click="cancelAddNode">Avbryt</Button>
          <Button primary @click="submitAddNode">Skapa</Button>
        </div>
      </div>
    </div>
    
    <!-- Rename node dialog -->
    <div v-if="showRenameNodeDialog" class="dialog-backdrop" @click="cancelRenameNode">
      <div class="dialog" @click.stop>
        <h3>Byt namn på mapp</h3>
        <div class="dialog-content">
          <label for="renameNodeField">Nytt mappnamn:</label>
          <input 
            id="renameNodeField" 
            v-model="renameNodeName" 
            type="text" 
            class="input-field" 
            @keyup.enter="submitRenameNode"
            ref="renameNodeInput"
            placeholder="Ange nytt mappnamn..." />
        </div>
        <div class="dialog-actions">
          <Button @click="cancelRenameNode">Avbryt</Button>
          <Button primary @click="submitRenameNode">Spara</Button>
        </div>
      </div>
    </div>
    
    <!-- Delete node confirmation dialog -->
    <div v-if="showDeleteNodeDialog" class="dialog-backdrop" @click="cancelDeleteNode">
      <div class="dialog" @click.stop>
        <h3>Ta bort mapp</h3>
        <div class="dialog-content">
          <p>Är du säker på att du vill ta bort mappen "{{ nodeToDeleteName }}"?</p>
          <p class="warning-text">Detta kommer också ta bort alla undermappar och dokument!</p>
        </div>
        <div class="dialog-actions">
          <Button @click="cancelDeleteNode">Avbryt</Button>
          <Button variant="danger" @click="submitDeleteNode">Ta bort</Button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, nextTick, onMounted, onUnmounted, defineComponent, markRaw } from 'vue';
import FolderNode from '../molecules/FolderNode.vue';
import Button from '../atoms/Button.vue';

// Define a new recursive NodeTree component to handle hierarchical display
const NodeTree = defineComponent({
  name: 'NodeTree',
  components: {
    FolderNode, // Make sure FolderNode is available within the component
  },
  props: {
    nodes: {
      type: Array,
      default: () => []
    },
    expandedNodeIds: {
      type: Object,
      default: () => new Set()
    },
    selectedNodeId: {
      type: String,
      default: null
    },
    level: {
      type: Number,
      default: 0
    }
  },
  emits: ['toggle-expand', 'select-node', 'contextmenu'],
  setup(props, { emit }) {
    const isNodeExpanded = (nodeId) => {
      return props.expandedNodeIds.has(nodeId);
    };
    
    // Handler for toggle-expand events from FolderNode
    const handleToggleExpand = (node) => {
      emit('toggle-expand', node);
    };
    
    return {
      isNodeExpanded,
      handleToggleExpand
    };
  },
  template: `
    <div class="node-tree-level" :style="{ paddingLeft: level > 0 ? '20px' : '0' }">
      <div v-for="node in nodes" :key="node.id" class="node-tree-item">
        <FolderNode 
          :node="node" 
          :is-expanded="isNodeExpanded(node.id)"
          :is-selected="selectedNodeId === node.id"
          @toggle-expand="handleToggleExpand"
          @select="$emit('select-node', node.id)"
          @contextmenu="$event => $emit('contextmenu', $event, node)" />
        
        <!-- Recursively render children if this node is expanded -->
        <NodeTree 
          v-if="isNodeExpanded(node.id) && node.children && node.children.length > 0" 
          :nodes="node.children"
          :expanded-node-ids="expandedNodeIds"
          :selected-node-id="selectedNodeId" 
          :level="level + 1"
          @toggle-expand="$emit('toggle-expand', $event)"
          @select-node="$emit('select-node', $event)"
          @contextmenu="($event, childNode) => $emit('contextmenu', $event, childNode)" />
      </div>
    </div>
  `
});

// Register the NodeTree component
const components = {
  NodeTree: markRaw(NodeTree),
  Button,
};

const props = defineProps({
  nodes: {
    type: Array,
    default: () => []
  },
  expandedNodeIds: {
    type: Object,
    default: () => new Set()
  },
  selectedNodeId: {
    type: String,
    default: null
  },
  isLoading: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['toggle-expand', 'select-node', 'add-node', 'rename-node', 'delete-node']);

// Context menu state
const contextMenu = reactive({
  show: false,
  x: 0,
  y: 0,
  nodeId: null
});

// Node dialog states
const showAddNodeDialog = ref(false);
const showRenameNodeDialog = ref(false);
const showDeleteNodeDialog = ref(false);
const newNodeName = ref('');
const renameNodeName = ref('');
const nodeNameInput = ref(null);
const renameNodeInput = ref(null);
const activeNodeId = ref(null);
const activeNode = ref(null);
const nodeToDeleteName = ref('');

// Computed props for context menu
const canRename = computed(() => {
  return activeNodeId.value !== '1'; // Don't allow renaming the root node
});

const canDelete = computed(() => {
  return activeNodeId.value !== '1'; // Don't allow deleting the root node
});

// Check if a node is expanded
const isNodeExpanded = (nodeId) => {
  return props.expandedNodeIds.has(nodeId);
};

// Methods
const handleToggleNode = (node) => {
  emit('toggle-expand', node.id);
};

const handleSelectNode = (nodeId) => {
  emit('select-node', nodeId);
};

const handleContextMenu = (event, node) => {
  event.preventDefault();
  contextMenu.show = true;
  contextMenu.x = event.clientX;
  contextMenu.y = event.clientY;
  activeNodeId.value = node.id;
  activeNode.value = node;
};

// Add Node methods
const handleAddNode = () => {
  showAddNodeDialog.value = true;
  nextTick(() => {
    if (nodeNameInput.value) {
      nodeNameInput.value.focus();
    }
  });
};

const submitAddNode = async () => {
  if (newNodeName.value.trim() !== '') {
    emit('add-node', newNodeName.value.trim(), activeNodeId.value);
    newNodeName.value = '';
    showAddNodeDialog.value = false;
  }
};

const cancelAddNode = () => {
  newNodeName.value = '';
  showAddNodeDialog.value = false;
};

// Rename Node methods
const handleRenameNode = () => {
  renameNodeName.value = activeNode.value?.name || '';
  showRenameNodeDialog.value = true;
  nextTick(() => {
    if (renameNodeInput.value) {
      renameNodeInput.value.focus();
    }
  });
};

const submitRenameNode = () => {
  if (renameNodeName.value.trim() !== '') {
    emit('rename-node', activeNodeId.value, renameNodeName.value.trim());
    renameNodeName.value = '';
    showRenameNodeDialog.value = false;
  }
};

const cancelRenameNode = () => {
  renameNodeName.value = '';
  showRenameNodeDialog.value = false;
};

// Delete Node methods
const handleDeleteNode = () => {
  nodeToDeleteName.value = activeNode.value?.name || '';
  showDeleteNodeDialog.value = true;
};

const submitDeleteNode = () => {
  emit('delete-node', activeNodeId.value);
  showDeleteNodeDialog.value = false;
};

const cancelDeleteNode = () => {
  showDeleteNodeDialog.value = false;
};

// Close context menu on click outside
const handleClickOutside = () => {
  if (contextMenu.show) {
    contextMenu.show = false;
  }
};

// Lifecycle hooks
onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script>

<style scoped>
.folder-panel {
  width: 280px;
  height: 100%;
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  background-color: var(--sidebar-background);
  color: var(--text-color);
  transition: background-color 0.3s, color 0.3s;
}

.folder-panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid var(--border-color);
}

.folder-panel-title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.folder-panel-icon {
  color: var(--primary-color);
}

.folder-panel h3 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 500;
}

.folder-panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
}

.node-tree {
  display: flex;
  flex-direction: column;
}

.node-tree-item {
  position: relative;
  margin-bottom: 2px;
}

/* Removed the ::before and ::after pseudo-elements for lines */

.node-tree-level {
  width: 100%;
  /* Padding is handled dynamically in the template */
}

/* Styles related to node children indentation if needed, but likely handled by NodeTree component's level prop */
/* .node-children { ... } */

.loading-state, .empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100px;
  color: var(--text-color-secondary);
}

/* Context Menu Styles */
.context-menu {
  position: fixed;
  background-color: var(--background-color);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  min-width: 150px;
}

.context-menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  cursor: pointer;
  color: var(--text-color);
}

.context-menu-item:hover {
  background-color: var(--primary-color-light-transparent);
}

.context-menu-icon {
  font-size: 18px;
}

/* Dialog Styles */
.dialog-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.dialog {
  background-color: var(--background-color);
  border-radius: 8px;
  width: 400px;
  max-width: 90vw;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.dialog h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: var(--text-color);
}

.dialog-content {
  margin-bottom: 20px;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.input-field {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  margin-top: 5px;
}

.warning-text {
  color: var(--danger-color);
  font-weight: 500;
}
</style>
