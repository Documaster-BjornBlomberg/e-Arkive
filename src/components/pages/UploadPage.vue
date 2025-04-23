<template>
  <MainLayout>
    <div class="upload-page">
      <div class="upload-container">
        <h1>Ladda upp dokument</h1>
        
        <FileUploadForm 
          @upload="handleUpload"
          :is-uploading="uploading"
          :available-nodes="availableNodes"
          :selected-node-id="selectedNodeId"
          :expanded-node-ids="expandedNodeIds"
          @node-change="handleNodeChange"
          @toggle-expand="handleNodeExpand" />
        
        <div v-if="statusMessage" 
             :class="['upload-status', { 'success': isSuccess, 'error': !isSuccess }]">
          {{ statusMessage }}
        </div>
        
        <div v-if="uploadedFile" class="upload-result">
          <h3>Dokument uppladdad!</h3>
          <div class="file-preview">
            <FileCardMolecule :file="uploadedFile" />
          </div>
          <div class="upload-actions">
            <Button icon="list" @click="navigateToList">Visa alla dokument</Button>
            <Button primary icon="upload" @click="resetUpload">Ladda upp fler</Button>
          </div>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import MainLayout from '../templates/MainLayout.vue';
import FileUploadForm from '../organisms/FileUploadForm.vue';
import Button from '../atoms/Button.vue';
import FileCardMolecule from '../molecules/FileCardMolecule.vue';
import { useGraphQL } from '../../composables/useGraphQL';
import { useNodeHandling } from '../../composables/useNodeHandling';

const { saveFile } = useGraphQL();
const { 
  getRootNodes, 
  getChildNodes 
} = useGraphQL();

const { 
  fetchRootNodes, 
  fetchNodeChildren, 
  getAllNodesForUpload 
} = useNodeHandling(); 

const router = useRouter();
const uploading = ref(false);
const uploadedFile = ref(null);
const statusMessage = ref('');
const isSuccess = ref(true);
const availableNodes = ref([]);
const selectedNodeId = ref('1'); // Default to root node
const expandedNodeIds = ref(new Set()); // Start with an empty Set so no nodes are expanded

// Handle node expansion toggle
const handleNodeExpand = async (nodeId) => {
  // Create a new Set to avoid directly mutating the ref value
  const newExpandedNodeIds = new Set([...expandedNodeIds.value]);
  
  if (newExpandedNodeIds.has(nodeId)) {
    newExpandedNodeIds.delete(nodeId);
  } else {
    newExpandedNodeIds.add(nodeId);
    
    // Load child nodes if needed
    const node = findNodeById(availableNodes.value, nodeId);
    if (node && (!node.children || node.children.length === 0) && node.hasChildren) {
      await loadNodeChildren(nodeId);
    }
  }
  
  expandedNodeIds.value = newExpandedNodeIds;
};

// Find a node by ID in the node tree
const findNodeById = (nodes, id) => {
  if (!nodes) return null;
  
  for (const node of nodes) {
    if (node.id === id) return node;
    if (node.children) {
      const found = findNodeById(node.children, id);
      if (found) return found;
    }
  }
  return null;
};

// Load children for a specific node and update the node structure
const loadNodeChildren = async (nodeId) => {
  try {
    console.log(`Loading children for node ${nodeId}`);
    const children = await getChildNodes(nodeId);
    
    if (children && children.length > 0) {
      // Mark each child node with hasChildren property and empty children array
      for (const childNode of children) {
        childNode.hasChildren = true; // Assume all nodes might have children until checked
        childNode.children = []; // Initialize empty children array
      }
      
      // Update the node in the tree to include its children
      const node = findNodeById(availableNodes.value, nodeId);
      if (node) {
        node.children = children;
        // Force reactivity update by creating a new array reference
        availableNodes.value = [...availableNodes.value];
        console.log(`Updated node ${nodeId} with children:`, children);
      } else {
        console.warn(`Could not find node ${nodeId} in availableNodes`);
      }
      return children;
    } else {
      console.log(`No children found for node ${nodeId}`);
      // If no children, update the node to indicate it has no children
      const node = findNodeById(availableNodes.value, nodeId);
      if (node) {
        node.hasChildren = false;
        // Force reactivity update
        availableNodes.value = [...availableNodes.value];
      }
      return [];
    }
  } catch (error) {
    console.error(`Error loading children for node ${nodeId}:`, error);
    return [];
  }
};

const handleNodeChange = (nodeId) => {
  selectedNodeId.value = nodeId;
};

const handleUpload = async (fileData) => {
  uploading.value = true;
  statusMessage.value = 'Laddar upp dokument...';
  isSuccess.value = true;
  
  try {
    if (selectedNodeId.value) {
      fileData.nodeId = selectedNodeId.value;
    }
    
    const result = await saveFile(fileData);
    uploadedFile.value = result;
    statusMessage.value = 'Dokumentet har laddats upp!';
  } catch (error) {
    console.error('Upload error:', error);
    statusMessage.value = `Fel vid uppladdning: ${error.message}`;
    isSuccess.value = false;
  } finally {
    uploading.value = false;
  }
};

// Initialize and load the nodes
const loadAvailableNodes = async () => {
  try {
    // First, fetch the root nodes directly
    const rootNodesResult = await getRootNodes();
    if (rootNodesResult && rootNodesResult.length > 0) {
      // For each root node, check if it has children by querying for its children
      for (const node of rootNodesResult) {
        // Mark that this node might have children (will be loaded when expanded)
        node.hasChildren = true;
        node.children = []; // Initialize empty children array
      }
      
      availableNodes.value = rootNodesResult;
      console.log("Root nodes loaded:", rootNodesResult);
    } else {
      console.warn("No root nodes returned from getRootNodes");
      availableNodes.value = [];
    }
  } catch (error) {
    console.error('Error loading root nodes:', error);
    availableNodes.value = [];
  }
};

// Reset upload form
const resetUpload = () => {
  uploadedFile.value = null;
  statusMessage.value = '';
};

// Navigate to list page
const navigateToList = () => {
  router.push('/list');
};

// Initialize when component is mounted
onMounted(async () => {
  await loadAvailableNodes();
});
</script>

<style scoped>
.upload-page {
  display: flex;
  justify-content: center;
  padding: 40px 20px;
  max-width: 100%;
}

.upload-container {
  width: 100%;
  max-width: 800px;
}

h1 {
  margin-bottom: 30px;
  color: var(--text-color);
}

.upload-status {
  margin-top: 20px;
  padding: 15px;
  border-radius: 4px;
  font-weight: 500;
}

.upload-status.success {
  background-color: var(--success-background);
  color: var(--success-color);
}

.upload-status.error {
  background-color: var(--error-background);
  color: var(--error-color);
}

.upload-result {
  margin-top: 30px;
  padding: 20px;
  background-color: var(--surface-color);
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.upload-result h3 {
  color: var(--text-color);
  margin-bottom: 15px;
}

.file-preview {
  margin: 15px 0;
}

.upload-actions {
  display: flex;
  gap: 15px;
  margin-top: 20px;
}
</style>
