<template>
  <div class="node-management">
    <div class="node-selector-header">
      <h3>Node Management</h3>
    </div>
    
    <div class="node-explorer">
      <div v-if="nodesLoading" class="loading">
        Loading nodes...
      </div>
      <div v-else class="node-tree">
        <div class="node-item" v-for="node in rootNodes" :key="node.id">
          <node-tree-item 
            :node="node"
            :expanded-nodes="expandedNodeIds"
            :selected-node-id="selectedNodeId"
            @select="selectNode"
            @toggle-expanded="toggleNodeExpanded"
          />
        </div>
      </div>
    </div>
    
    <div class="node-details" v-if="selectedNode">
      <h4 class="details-header">Node Details: {{ selectedNode.name }}</h4>
      
      <div class="form-group">
        <label for="nodeName">Name</label>
        <input type="text" id="nodeName" v-model="selectedNode.name" />
      </div>
      
      <!-- Ownership Section -->
      <div class="ownership-section">
        <h5>Ownership</h5>
        
        <div class="form-group">
          <label for="ownerUser">Owner User</label>
          <select id="ownerUser" v-model="selectedOwnerUserId">
            <option value="">None</option>
            <option v-for="user in users" :key="user.id" :value="user.id">
              {{ user.name }} (@{{ user.username }})
            </option>
          </select>
        </div>
        
        <div class="form-group">
          <label for="ownerGroup">Owner Group</label>
          <select id="ownerGroup" v-model="selectedOwnerGroupId">
            <option value="">None</option>
            <option v-for="group in groups" :key="group.id" :value="group.id">
              {{ group.name }}
            </option>
          </select>
        </div>
      </div>
      
      <!-- Permissions Section -->
      <node-permissions v-model="nodePermissions" />
      
      <div class="node-actions">
        <button @click="saveNodeChanges" class="save-btn">Save Changes</button>
        <button @click="createChildNode" class="add-btn">Add Child Node</button>
        <button @click="deleteSelectedNode" class="delete-btn">Delete Node</button>
      </div>
    </div>
    
    <!-- Create Node Modal -->
    <modal v-if="showCreateNodeModal" @close="showCreateNodeModal = false">
      <template #header>
        <h3>Create New Node</h3>
      </template>
      
      <template #body>
        <div class="form-group">
          <label for="newNodeName">Node Name</label>
          <input type="text" id="newNodeName" v-model="newNodeName" />
        </div>
        
        <p>This node will be created as a child of: <strong>{{ getParentNodeName() }}</strong></p>
      </template>
      
      <template #footer>
        <button @click="confirmCreateNode" class="save-btn" :disabled="!newNodeName.trim()">Create</button>
        <button @click="showCreateNodeModal = false" class="cancel-btn">Cancel</button>
      </template>
    </modal>
    
    <!-- Delete Confirmation Modal -->
    <modal v-if="showDeleteNodeModal" @close="showDeleteNodeModal = false">
      <template #header>
        <h3>Confirm Delete</h3>
      </template>
      
      <template #body>
        <p>Are you sure you want to delete the node <strong>{{ selectedNode?.name }}</strong>?</p>
        <p class="warning">Warning: This will also delete all child nodes and associated files!</p>
      </template>
      
      <template #footer>
        <button @click="confirmDeleteNode" class="delete-btn">Delete</button>
        <button @click="showDeleteNodeModal = false" class="cancel-btn">Cancel</button>
      </template>
    </modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useNodeHandling } from '../../composables/useNodeHandling';
import { useUserManagement } from '../../composables/useUserManagement';
import NodeTreeItem from '../molecules/NodeTreeItem.vue';
import NodePermissions from '../molecules/NodePermissions.vue';
import Modal from './Modal.vue';

// Composables
const { 
  rootNodes, 
  expandedNodeIds, 
  nodeLoading: nodesLoading,
  fetchRootNodes,
  fetchNodeById,
  addNewNode,
  updateNodeName,
  removeNode
} = useNodeHandling();

const {
  users,
  groups,
  getUsers,
  getGroups,
  setNodePermissions,
  setNodeOwnership
} = useUserManagement();

// Local state
const selectedNodeId = ref<string | null>(null);
const selectedNode = ref<any | null>(null);
const selectedOwnerUserId = ref<string>('');
const selectedOwnerGroupId = ref<string>('');
const nodePermissions = ref<number>(0);
const showCreateNodeModal = ref(false);
const showDeleteNodeModal = ref(false);
const newNodeName = ref('');

// Initialize component
onMounted(async () => {
  await fetchRootNodes();
  await getUsers();
  await getGroups();
});

// Methods
const selectNode = async (nodeId: string) => {
  selectedNodeId.value = nodeId;
  
  // Fetch the complete node details
  const node = await fetchNodeById(nodeId);
  if (node) {
    selectedNode.value = node;
    selectedOwnerUserId.value = node.ownerUserId || '';
    selectedOwnerGroupId.value = node.ownerGroupId || '';
    nodePermissions.value = node.permissions || 0;
  }
};

const toggleNodeExpanded = (nodeId: string) => {
  if (expandedNodeIds.has(nodeId)) {
    expandedNodeIds.delete(nodeId);
  } else {
    expandedNodeIds.add(nodeId);
  }
};

const saveNodeChanges = async () => {
  if (!selectedNode.value) return;
  
  try {
    // Update node name if needed
    await updateNodeName(selectedNode.value.id, selectedNode.value.name);
    
    // Update ownership if changed
    await setNodeOwnership(
      selectedNode.value.id, 
      selectedOwnerUserId.value || undefined, 
      selectedOwnerGroupId.value || undefined
    );
    
    // Update permissions if changed
    await setNodePermissions(selectedNode.value.id, nodePermissions.value);
    
    // Refresh node data
    await fetchRootNodes();
    
    // Re-select the node to get fresh data
    if (selectedNodeId.value) {
      await selectNode(selectedNodeId.value);
    }
  } catch (err) {
    console.error('Error saving node changes:', err);
  }
};

const createChildNode = () => {
  newNodeName.value = '';
  showCreateNodeModal.value = true;
};

const confirmCreateNode = async () => {
  if (!selectedNodeId.value || !newNodeName.value.trim()) return;
  
  try {
    // Create new node as child of selected node
    await addNewNode(newNodeName.value, selectedNodeId.value);
    
    // Refresh nodes
    await fetchRootNodes();
    
    // Ensure parent node is expanded
    if (selectedNodeId.value) {
      expandedNodeIds.add(selectedNodeId.value);
    }
    
    // Close modal
    showCreateNodeModal.value = false;
    newNodeName.value = '';
  } catch (err) {
    console.error('Error creating node:', err);
  }
};

const deleteSelectedNode = () => {
  if (!selectedNode.value) return;
  showDeleteNodeModal.value = true;
};

const confirmDeleteNode = async () => {
  if (!selectedNodeId.value) return;
  
  try {
    await removeNode(selectedNodeId.value);
    
    // Reset selection and refresh nodes
    selectedNodeId.value = null;
    selectedNode.value = null;
    await fetchRootNodes();
    
    // Close modal
    showDeleteNodeModal.value = false;
  } catch (err) {
    console.error('Error deleting node:', err);
  }
};

const getParentNodeName = () => {
  if (!selectedNodeId.value) return 'Root';
  return selectedNode.value?.name || 'Selected Node';
};
</script>

<style scoped>
.node-management {
  display: flex;
  flex-direction: column;
}

.node-selector-header {
  margin-bottom: 16px;
}

.node-explorer {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background-color: var(--surface-color);
  padding: 16px;
  margin-bottom: 24px;
  max-height: 300px;
  overflow-y: auto;
}

.loading {
  padding: 20px;
  text-align: center;
}

.node-tree {
  user-select: none;
}

.node-details {
  background-color: var(--surface-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 16px;
}

.details-header {
  margin-top: 0;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-color);
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-weight: bold;
  color: var(--text-color);
}

.form-group input, .form-group select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background-color: var(--input-background);
  color: var(--text-color);
}

.ownership-section {
  background-color: var(--background-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 16px;
}

.ownership-section h5 {
  margin-top: 0;
  margin-bottom: 16px;
  color: var(--text-color-secondary);
}

.node-actions {
  display: flex;
  gap: 8px;
  margin-top: 16px;
}

.save-btn, .add-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 4px;
  padding: 8px 16px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.save-btn:hover, .add-btn:hover {
  background-color: var(--button-hover-bg);
}

.delete-btn {
  background-color: var(--error-color);
  color: white;
  border: none;
  border-radius: 4px;
  padding: 8px 16px;
  cursor: pointer;
  margin-left: auto;
  transition: background-color 0.2s;
}

.delete-btn:hover {
  opacity: 0.9;
}

.cancel-btn {
  background-color: var(--border-color);
  color: var(--text-color);
  border: none;
  border-radius: 4px;
  padding: 8px 16px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.cancel-btn:hover {
  background-color: var(--hover-color);
}

.warning {
  color: var(--error-color);
  font-weight: bold;
}
</style>
