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

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue';
import MainLayout from '../templates/MainLayout.vue';
import Button from '../atoms/Button.vue';
import FileList from '../organisms/FileList.vue';
import FileDetailSidepanel from '../organisms/FileDetailSidepanel.vue';
import FolderPanel from '../organisms/FolderPanel.vue';
import { useGraphQL } from '../../composables/useGraphQL';
import { useNodeHandling } from '../../composables/useNodeHandling';
import { File, Metadata, MetadataInput } from '../../types/schema';
import { ViewMode, ViewModeOption, StatusType } from '../../types/ui';

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
  fetchNodeById,
  toggleNodeExpand,
  selectNode,
  addNewNode,
  updateNodeName,
  removeNode,
  initialize
} = useNodeHandling();

const viewModes: ViewModeOption[] = [
  { id: 'table', label: 'Tabellvy', icon: 'table_rows' },
  { id: 'grid', label: 'Rutnätsvy', icon: 'grid_view' },
  { id: 'card', label: 'Kortvy', icon: 'view_agenda' },
  { id: 'split', label: 'Delad vy', icon: 'view_sidebar' }
];

// UI state
const viewMode = ref<ViewMode>(localStorage.getItem('preferredViewMode') as ViewMode || 'table');
const files = ref<File[]>([]);
const selectedFileId = ref<string | undefined>(undefined);
const expandedFileId = ref<string | undefined>(undefined);
const selectedFile = ref<File | null>(null);
const isEditingMetadata = ref<boolean>(false);
const editingMetadata = ref<MetadataInput[]>([]);
const activeMetadataTab = ref<string>('list');
const metadataSearch = ref<string>('');
const statusMessage = ref<string>('');
const statusType = ref<StatusType | ''>(''); // 'success', 'error', 'info'

// Computed property for the current node name
const currentNodeName = computed(() => {
  const currentNodeTyped = currentNode.value as { name?: string } | null;
  if (currentNodeTyped && currentNodeTyped.name) {
    return currentNodeTyped.name;
  }
  return 'Alla dokument';
});

// Computed property to determine if sidepanel should be shown
const showSidepanel = computed(() => {
  return viewMode.value !== 'split';
});

// Update view mode and save preference
const updateViewMode = (mode: ViewMode): void => {
  viewMode.value = mode;
  localStorage.setItem('preferredViewMode', mode);
};

// Handle file selection
const handleFileSelect = async (fileId: string): Promise<void> => {
  if (selectedFileId.value === fileId && viewMode.value !== 'split') {
    // Toggle off if already selected, except in split view
    selectedFileId.value = undefined;
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
const handleFileExpand = (fileId: string): void => {
  if (expandedFileId.value === fileId) {
    expandedFileId.value = undefined;
  } else {
    expandedFileId.value = fileId;
  }
};

// Handle closing the sidepanel
const handleCloseSidepanel = (): void => {
  if (viewMode.value !== 'split') {
    selectedFileId.value = undefined;
    selectedFile.value = null;
  }
  isEditingMetadata.value = false;
};

// Handle node selection
const handleNodeSelect = async (nodeId: string): Promise<void> => {
  try {
    // Use the selectNode function that will fetch files and expand the node
    await selectNode(nodeId);
    
    // Clear any selected file when changing nodes
    selectedFileId.value = undefined;
    selectedFile.value = null;
    
    // Show status message for better user feedback
    const currentNodeTyped = currentNode.value as { name?: string } | null;
    const nodeName = currentNodeTyped && currentNodeTyped.name ? currentNodeTyped.name : 'mapp';
    showStatus(`Visar innehåll från "${nodeName}"`, 'info');
  } catch (error) {
    console.error('Error selecting node:', error);
    showStatus('Kunde inte ladda mapinnehåll', 'error');
  }
};

// Handle node expansion toggle
const handleNodeExpand = async (nodeId: string): Promise<void> => {
  try {
    await toggleNodeExpand(nodeId);
  } catch (error) {
    console.error('Error toggling node expansion:', error);
  }
};

// Handle node addition
const handleAddNode = async (nodeName: string, parentId?: string): Promise<void> => {
  try {
    await addNewNode(nodeName, parentId);
    showStatus('Mapp skapad', 'success');
  } catch (error: any) {
    console.error('Error creating folder:', error);
    showStatus('Kunde inte skapa mapp: ' + error.message, 'error');
  }
};

// Handle node rename
const handleRenameNode = async (nodeId: string, newName: string): Promise<void> => {
  try {
    await updateNodeName(nodeId, newName);
    showStatus('Mappnamn uppdaterat', 'success');
  } catch (error: any) {
    console.error('Error renaming folder:', error);
    showStatus('Kunde inte byta namn på mapp: ' + error.message, 'error');
  }
};

// Handle node deletion
const handleDeleteNode = async (nodeId: string): Promise<void> => {
  try {
    await removeNode(nodeId);
    showStatus('Mapp borttagen', 'success');
  } catch (error: any) {
    console.error('Error deleting folder:', error);
    showStatus('Kunde inte ta bort mapp: ' + error.message, 'error');
  }
};

// Show status message
const showStatus = (message: string, type: StatusType = 'info') => {
  statusMessage.value = message;
  statusType.value = type;
  
  // Auto hide after 3 seconds
  setTimeout(() => {
    statusMessage.value = '';
  }, 3000);
};

// Handle download file
const handleDownloadFile = (file: File): void => {
  // Implement download logic
  console.log('Downloading file:', file);
};

// Handle edit metadata
const handleEditMetadata = (file: File): void => {
  editingMetadata.value = [...(file.metadata || []).map(m => ({key: m.key, value: m.value}))];
  if (editingMetadata.value.length === 0) {
    editingMetadata.value.push({ key: '', value: '' });
  }
  isEditingMetadata.value = true;
};

// Handle save metadata
const handleSaveMetadata = async (): Promise<void> => {
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
const handleCancelEdit = (): void => {
  isEditingMetadata.value = false;
};

// Handle delete file
const handleDeleteFile = async (fileId: string): Promise<void> => {
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
const handleDeleteMetadata = async (fileId: string, metadataKey: string | string[]): Promise<void> => {
  try {
    const keysToDelete = Array.isArray(metadataKey) ? metadataKey : [metadataKey];
    await deleteFileMetadata(fileId, keysToDelete);
    
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
const handleRefresh = async (): Promise<void> => {
  try {
    await fetchNodeFiles(selectedNodeId.value);
    showStatus('Innehåll uppdaterat', 'info');
  } catch (error) {
    console.error('Error refreshing content:', error);
    showStatus('Kunde inte uppdatera innehåll', 'error');
  }
};

// Initialize the page
onMounted(async (): Promise<void> => {
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
