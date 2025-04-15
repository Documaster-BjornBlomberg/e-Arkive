<template>
  <MainLayout>
    <div class="list-page">
      <div class="list-page-content">
        <!-- Folder panel with node/folder structure -->
        <FolderPanel v-if="isNodePanelOpen" 
                    :is-loading="nodeLoading" 
                    :nodes="rootNodes"
                    :expanded-node-ids="expandedNodeIds"
                    :selected-node-id="selectedNodeId"
                    @toggle-expand="handleNodeExpand" 
                    @select-node="handleNodeSelect"
                    @add-node="handleAddNode"
                    @rename-node="handleRenameNode"
                    @delete-node="handleDeleteNode" />
        
        <!-- Main content area with file list -->
        <div class="file-content-area">
          <div class="list-header">
            <div class="list-header-content">
              <h2>{{ currentNodeName }}</h2>
              <div class="list-actions">
                <div class="view-mode-buttons">
                  <Button 
                    v-for="mode in viewModes"
                    :key="mode.id"
                    :icon="mode.icon"
                    :variant="viewMode === mode.id ? 'primary' : 'default'"
                    @click="viewMode = mode.id"
                    :title="mode.label">
                    {{ mode.label }}
                  </Button>
                </div>
                <Button 
                  icon="refresh"
                  @click="handleRefresh"
                  variant="default"
                  title="Uppdatera">
                  Uppdatera
                </Button>
              </div>
            </div>
          </div>
          
          <FileList 
            :files="nodeFiles" 
            :is-loading="loading" 
            :selected-file-id="selectedFileId"
            :expanded-file-id="expandedFileId" 
            :view-mode="viewMode"
            @select-file="handleFileSelect" 
            @toggle-expand="handleFileExpand"
            @update:viewMode="updateViewMode" />
        </div>
        
        <!-- Detail side panel for selected file -->
        <FileDetailSidepanel 
          v-if="selectedFile && showSidepanel" 
          :file="selectedFile"
          :is-open="!!selectedFile"
          :is-editing="isEditingMetadata"
          :editing-metadata="editingMetadata"
          :active-metadata-tab="activeMetadataTab"
          :metadata-search="metadataSearch"
          @close="handleCloseSidepanel"
          @download="handleDownloadFile"
          @edit-metadata="handleEditMetadata"
          @save-metadata="handleSaveMetadata"
          @cancel-edit="handleCancelEdit"
          @delete-file="handleDeleteFile"
          @delete-metadata="handleDeleteMetadata"
          @update:activeTab="activeMetadataTab = $event"
          @update:metadataSearch="metadataSearch = $event"
          @update:editingMetadata="editingMetadata = $event" />
      </div>
    </div>
  </MainLayout>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue';
import MainLayout from '../templates/MainLayout.vue';
import Button from '../atoms/Button.vue';
import FileList from '../organisms/FileList.vue';
import FileDetailSidepanel from '../organisms/FileDetailSidepanel.vue';
import FolderPanel from '../organisms/FolderPanel.vue';
import { useGraphQL } from '../../composables/useGraphQL';
import { useNodeHandling } from '../../composables/useNodeHandling';

// Set up GraphQL client
const { loading, error, getFiles, getFileById, updateFileMetadata, deleteFile, deleteFileMetadata } = useGraphQL();

// Setup node handling for folders
const { 
  nodeLoading, 
  nodeError,
  rootNodes, 
  nodeFiles,
  currentNode,
  isNodePanelOpen, 
  selectedNodeId,
  expandedNodeIds,
  fetchRootNodes,
  fetchNodeFiles,
  fetchNodeChildren,
  toggleNodeExpand,
  selectNode,
  addNewNode,
  updateNodeName,
  removeNode,
  initialize
} = useNodeHandling();

const viewModes = [
  { id: 'table', label: 'Tabellvy', icon: 'table_rows' },
  { id: 'grid', label: 'Rutnätsvy', icon: 'grid_view' },
  { id: 'card', label: 'Kortvy', icon: 'view_agenda' },
  { id: 'split', label: 'Delad vy', icon: 'view_sidebar' }
];

// UI state
const viewMode = ref(localStorage.getItem('preferredViewMode') || 'table');
const files = ref([]);
const selectedFileId = ref(null);
const expandedFileId = ref(null);
const selectedFile = ref(null);
const isEditingMetadata = ref(false);
const editingMetadata = ref([]);
const activeMetadataTab = ref('list');
const metadataSearch = ref('');
const statusMessage = ref('');
const statusType = ref(''); // 'success', 'error', 'info'

// Computed property for the current node name
const currentNodeName = computed(() => {
  if (currentNode.value) {
    return currentNode.value.name;
  }
  return 'Alla dokument';
});

// Computed property to determine if sidepanel should be shown
const showSidepanel = computed(() => {
  return viewMode.value !== 'split';
});

// Update view mode and save preference
const updateViewMode = (mode) => {
  viewMode.value = mode;
  localStorage.setItem('preferredViewMode', mode);
};

// Handle file selection
const handleFileSelect = async (fileId) => {
  if (selectedFileId.value === fileId && viewMode.value !== 'split') {
    // Toggle off if already selected, except in split view
    selectedFileId.value = null;
    selectedFile.value = null;
    return;
  }
  
  selectedFileId.value = fileId;
  
  try {
    const fileData = await getFileById(fileId);
    selectedFile.value = fileData;
  } catch (err) {
    console.error('Error fetching file details:', err);
    selectedFile.value = null;
  }
};

// Handle file expand toggle
const handleFileExpand = (fileId) => {
  if (expandedFileId.value === fileId) {
    expandedFileId.value = null;
  } else {
    expandedFileId.value = fileId;
  }
};

// Handle closing the sidepanel
const handleCloseSidepanel = () => {
  if (viewMode.value !== 'split') {
    selectedFileId.value = null;
    selectedFile.value = null;
  }
  isEditingMetadata.value = false;
};

// Handle node selection
const handleNodeSelect = async (nodeId) => {
  try {
    // Use the selectNode function that will fetch files and expand the node
    await selectNode(nodeId);
    
    // Clear any selected file when changing nodes
    selectedFileId.value = null;
    selectedFile.value = null;

    // Show status message for better user feedback
    showStatus(`Visar innehåll från "${currentNode.value?.name || 'mapp'}"`, 'info');
  } catch (error) {
    console.error('Error selecting node:', error);
    showStatus('Kunde inte ladda mapinnehåll', 'error');
  }
};

// Handle node expansion toggle
const handleNodeExpand = async (nodeId) => {
  try {
    await toggleNodeExpand(nodeId);
  } catch (error) {
    console.error('Error toggling node expansion:', error);
  }
};

// Handle node addition
const handleAddNode = async (nodeName, parentId) => {
  try {
    await addNewNode(nodeName, parentId);
    showStatus('Mapp skapad', 'success');
  } catch (error) {
    console.error('Error creating folder:', error);
    showStatus('Kunde inte skapa mapp: ' + error.message, 'error');
  }
};

// Handle node rename
const handleRenameNode = async (nodeId, newName) => {
  try {
    await updateNodeName(nodeId, newName);
    showStatus('Mappnamn uppdaterat', 'success');
  } catch (error) {
    console.error('Error renaming folder:', error);
    showStatus('Kunde inte byta namn på mapp: ' + error.message, 'error');
  }
};

// Handle node deletion
const handleDeleteNode = async (nodeId) => {
  try {
    await removeNode(nodeId);
    showStatus('Mapp borttagen', 'success');
  } catch (error) {
    console.error('Error deleting folder:', error);
    showStatus('Kunde inte ta bort mapp: ' + error.message, 'error');
  }
};

// Show status message
const showStatus = (message, type = 'info') => {
  statusMessage.value = message;
  statusType.value = type;
  
  // Auto hide after 3 seconds
  setTimeout(() => {
    statusMessage.value = '';
  }, 3000);
};

// Handle download file
const handleDownloadFile = (file) => {
  // Implement download logic
  console.log('Downloading file:', file);
};

// Handle edit metadata
const handleEditMetadata = (file) => {
  editingMetadata.value = [...(file.metadata || [])];
  if (editingMetadata.value.length === 0) {
    editingMetadata.value.push({ key: '', value: '' });
  }
  isEditingMetadata.value = true;
};

// Handle save metadata
const handleSaveMetadata = async () => {
  if (!selectedFile.value) return;
  
  try {
    // Filter out empty metadata entries
    const validMetadata = editingMetadata.value.filter(meta => meta.key && meta.value);
    
    await updateFileMetadata(selectedFile.value.id, validMetadata);
    
    // Refresh file data
    const updatedFile = await getFileById(selectedFile.value.id);
    selectedFile.value = updatedFile;
    
    isEditingMetadata.value = false;
    showStatus('Metadata sparad', 'success');
  } catch (err) {
    console.error('Error saving metadata:', err);
    showStatus('Kunde inte spara metadata', 'error');
  }
};

// Handle cancel edit
const handleCancelEdit = () => {
  isEditingMetadata.value = false;
};

// Handle delete file
const handleDeleteFile = async (fileId) => {
  if (confirm('Är du säker på att du vill ta bort denna fil? Denna åtgärd kan inte ångras.')) {
    try {
      await deleteFile(fileId);
      handleCloseSidepanel();
      await fetchNodeFiles(selectedNodeId.value);
      showStatus('Fil borttagen', 'success');
    } catch (err) {
      console.error('Error deleting file:', err);
      showStatus('Kunde inte ta bort fil', 'error');
    }
  }
};

// Handle delete metadata
const handleDeleteMetadata = async (fileId, metadataKey) => {
  try {
    await deleteFileMetadata(fileId, metadataKey);
    
    // Refresh file data
    const updatedFile = await getFileById(fileId);
    selectedFile.value = updatedFile;
    showStatus('Metadata borttagen', 'success');
  } catch (err) {
    console.error('Error deleting metadata:', err);
    showStatus('Kunde inte ta bort metadata', 'error');
  }
};

// Handle refresh
const handleRefresh = async () => {
  try {
    await fetchNodeFiles(selectedNodeId.value);
    showStatus('Innehåll uppdaterat', 'info');
  } catch (error) {
    console.error('Error refreshing content:', error);
    showStatus('Kunde inte uppdatera innehåll', 'error');
  }
};

// Initialize the page
onMounted(async () => {
  await initialize();
});
</script>

<style scoped>
.list-page {
  height: calc(100vh - 70px);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.list-page-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.file-content-area {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.list-header {
  padding: 15px 20px;
  background: var(--surface-color);
  border-bottom: 1px solid var(--border-color);
}

.list-header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.list-header h2 {
  margin: 0;
}

.list-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.view-mode-buttons {
  display: flex;
  gap: 5px;
}

/* ... existing styles ... */
</style>
