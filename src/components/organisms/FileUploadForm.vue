<template>
  <div class="file-upload-form">
    <div class="upload-layout">
      <!-- Folder panel on the left -->
      <FolderPanel
        :nodes="availableNodes"
        :is-loading="isNodeLoading"
        :expanded-node-ids="localExpandedNodeIds"
        :selected-node-id="localSelectedNodeId"
        @toggle-expand="handleNodeExpand"
        @select-node="handleNodeSelect"
        @add-node="handleAddNode"
        @rename-node="handleRenameNode"
        @delete-node="handleDeleteNode"
      />

      <!-- Upload form on the right -->
      <div class="upload-content">
        <!-- Selected folder path -->
        <div class="selected-folder-path" v-if="currentNodePath.length > 0">
          <span class="material-icons">folder</span>
          <div class="path-segments">
            <span v-for="(segment, index) in currentNodePath" :key="index">
              {{ segment }}
              <span v-if="index < currentNodePath.length - 1" class="path-separator">/</span>
            </span>
          </div>
        </div>

        <form @submit.prevent="handleSubmit">
          <!-- File Input Area -->
          <div class="file-drop-area" @dragover.prevent="fileDragOver = true" @dragleave.prevent="fileDragOver = false" @drop.prevent="onFileDrop">
            <div class="file-drop-content" :class="{ 'drag-over': fileDragOver, 'has-file': selectedFile }">
              <div v-if="!selectedFile">
                <span class="material-icons icon-large">cloud_upload</span>
                <p>Dra och släpp filen här eller</p>
                <div class="file-button-wrapper">
                  <Button @click="triggerFileInput" variant="primary">Välj fil</Button>
                  <input 
                    type="file" 
                    ref="fileInputRef" 
                    class="hidden-file-input"
                    @change="onFileSelect"
                    :accept="allowedFileTypes"
                  />
                </div>
              </div>
              <div v-else class="selected-file-preview">
                <FileIcon :fileName="selectedFile.name" :fileType="fileType" />
                <div class="file-info">
                  <div class="file-name">{{ selectedFile.name }}</div>
                  <div class="file-size">{{ formatFileSize(selectedFile.size) }}</div>
                  <div class="selected-folder" v-if="currentNodeName">
                    <span class="material-icons">folder</span>
                    {{ currentNodeName }}
                  </div>
                </div>
                <button type="button" class="remove-file-btn" @click="removeFile">
                  <span class="material-icons">close</span>
                </button>
              </div>
            </div>
          </div>

          <!-- Metadata Section -->
          <div class="metadata-section">
            <h3>Metadata</h3>
            <p class="section-description">Lägg till relevant information om dokumentet</p>

            <div v-for="(meta, index) in metadata" :key="index" class="metadata-field">
              <div class="metadata-inputs">
                <InputField 
                  label="Nyckel" 
                  v-model="meta.key" 
                  placeholder="Nyckel"
                  class="metadata-key" />
                <InputField 
                  label="Värde" 
                  v-model="meta.value" 
                  placeholder="Värde"
                  class="metadata-value" />
              </div>
              <button type="button" class="remove-metadata-btn" @click="removeMetadata(index)">
                <span class="material-icons">delete_outline</span>
              </button>
            </div>

            <button type="button" class="add-metadata-btn" @click="addMetadata">
              <span class="material-icons">add</span>
              Lägg till metadata
            </button>
          </div>

          <div class="form-actions">
            <Button 
              type="submit" 
              primary 
              :disabled="!formValid || isUploading"
              :loading="isUploading">
              {{ isUploading ? 'Laddar upp...' : 'Ladda upp dokument' }}
            </Button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import Button from '../atoms/Button.vue';
import FileIcon from '../atoms/FileIcon.vue';
import InputField from '../atoms/InputField.vue';
import FolderPanel from '../organisms/FolderPanel.vue';

const props = defineProps({
  isUploading: {
    type: Boolean,
    default: false
  },
  availableNodes: {
    type: Array,
    default: () => []
  },
  selectedNodeId: {
    type: String,
    default: '1' // Default to root node
  },
  expandedNodeIds: {
    type: Object, // Set object
    default: () => new Set()
  }
});

const emit = defineEmits(['upload', 'node-change', 'toggle-expand']);

// Local state to manage node selection and expansion
const localSelectedNodeId = ref(props.selectedNodeId);
const localExpandedNodeIds = ref(new Set()); // Initialize as empty set
const isNodeLoading = ref(false);

// Watch for changes in props
watch(() => props.selectedNodeId, (newValue) => {
  localSelectedNodeId.value = newValue;
});

watch(() => props.expandedNodeIds, (newValue) => {
  // Only update if the sets are actually different
  const newSet = new Set([...newValue]);
  const localSet = new Set([...localExpandedNodeIds.value]);
  
  if ([...newSet].sort().join(',') !== [...localSet].sort().join(',')) {
    localExpandedNodeIds.value = newSet;
  }
}, { deep: true });

const selectedFile = ref(null);
const fileDragOver = ref(false);
const fileType = ref('');
const fileInputRef = ref(null);
const metadata = ref([
  { key: '', value: '' }
]);

// List of allowed file types
const allowedFileTypes = '.pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.txt,.jpg,.jpeg,.png,.gif';

// Track form validity
const formValid = computed(() => {
  return selectedFile.value !== null;
});

// Computed properties for folder display
const currentNodeName = computed(() => {
  const selectedNode = findNodeById(props.availableNodes, localSelectedNodeId.value);
  return selectedNode?.name || '';
});

// Computed path to current node
const currentNodePath = computed(() => {
  const path = [];
  let currentNode = findNodeById(props.availableNodes, localSelectedNodeId.value);
  
  while (currentNode) {
    path.unshift(currentNode.name);
    if (!currentNode.parentId) break;
    currentNode = findNodeById(props.availableNodes, currentNode.parentId);
  }
  
  return path;
});

// Helper to find node by ID
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

// Node handling methods
const handleNodeExpand = (nodeId) => {
  console.log(`FileUploadForm: Node ${nodeId} toggle expand clicked`);
  
  // Create a new Set to avoid directly mutating the ref value
  const newExpandedNodeIds = new Set([...localExpandedNodeIds.value]);
  
  if (newExpandedNodeIds.has(nodeId)) {
    newExpandedNodeIds.delete(nodeId);
  } else {
    newExpandedNodeIds.add(nodeId);
  }
  
  localExpandedNodeIds.value = newExpandedNodeIds;
  
  // Important: Emit the toggle-expand event to the parent component to load children if needed
  emit('toggle-expand', nodeId);
};

const handleNodeSelect = (nodeId) => {
  localSelectedNodeId.value = nodeId;
  emit('node-change', nodeId);
};

const handleAddNode = (nodeName, parentId) => {
  // Implementation left to parent component 
};

const handleRenameNode = (nodeId, newName) => {
  // Implementation left to parent component
};

const handleDeleteNode = (nodeId) => {
  // Implementation left to parent component
};

// File type detection helper
const getFileType = (filename) => {
  const extension = filename.split('.').pop().toLowerCase();
  const imageTypes = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'svg', 'webp'];
  const documentTypes = ['pdf', 'doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx', 'txt'];

  if (imageTypes.includes(extension)) return 'image';
  if (documentTypes.includes(extension)) return 'document';
  return 'file';
};

// Trigger file input click
const triggerFileInput = () => {
  if (fileInputRef.value) {
    fileInputRef.value.click();
  }
};

// Handle file selection from input
const onFileSelect = (event) => {
  const file = event.target.files[0];
  if (file) {
    selectedFile.value = file;
    fileType.value = getFileType(file.name);
  }
};

// Handle file drop event
const onFileDrop = (event) => {
  fileDragOver.value = false;
  const file = event.dataTransfer.files[0];
  if (file) {
    selectedFile.value = file;
    fileType.value = getFileType(file.name);
  }
};

// Remove selected file
const removeFile = () => {
  selectedFile.value = null;
  fileType.value = '';
};

// Format file size for display
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes';
  
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

// Add empty metadata field
const addMetadata = () => {
  metadata.value.push({ key: '', value: '' });
};

// Remove metadata field
const removeMetadata = (index) => {
  metadata.value.splice(index, 1);
  if (metadata.value.length === 0) {
    addMetadata(); // Always have at least one metadata field
  }
};

// Form submission handler
const handleSubmit = async () => {
  if (!formValid.value || props.isUploading) return;

  const file = selectedFile.value;
  
  // Read file as base64
  const reader = new FileReader();
  reader.readAsDataURL(file);
  
  reader.onload = () => {
    // Extract the base64 data part (remove data:mime/type;base64, prefix)
    const base64Data = reader.result.split(',')[1];
    
    // Filter out empty metadata entries
    const validMetadata = metadata.value.filter(meta => 
      meta.key.trim() !== '' && meta.value.trim() !== ''
    );

    // Prepare file data for upload
    const fileData = {
      name: file.name,
      size: file.size,
      contentType: file.type,
      fileData: base64Data,
      metadata: validMetadata,
      nodeId: localSelectedNodeId.value
    };

    // Emit the upload event with the file data
    emit('upload', fileData);
  };
  
  reader.onerror = error => {
    console.error('Error reading file:', error);
  };
};
</script>

<style scoped>
.file-upload-form {
  background: var(--surface-color);
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.file-drop-area {
  margin-bottom: 30px;
}

.file-drop-content {
  border: 2px dashed var(--border-color);
  border-radius: 8px;
  padding: 40px 20px;
  text-align: center;
  color: var(--text-color-secondary);
  transition: all 0.3s ease;
}

.file-drop-content.drag-over {
  border-color: var(--primary-color);
  background-color: var(--primary-color-light-transparent);
}

.file-drop-content.has-file {
  border-style: solid;
  background-color: var(--surface-color);
  padding: 20px;
}

.icon-large {
  font-size: 48px;
  margin-bottom: 15px;
  color: var(--text-color-secondary);
}

.file-button-wrapper {
  margin-top: 15px;
  display: inline-block;
}

.hidden-file-input {
  display: none;
}

.selected-file-preview {
  display: flex;
  align-items: center;
  gap: 15px;
}

.file-info {
  flex: 1;
  text-align: left;
}

.file-name {
  font-weight: 500;
  margin-bottom: 5px;
  color: var(--text-color);
}

.file-size {
  color: var(--text-color-secondary);
  font-size: 14px;
}

.remove-file-btn {
  background: none;
  border: none;
  color: var(--danger-color);
  cursor: pointer;
  padding: 5px;
}

.metadata-section {
  margin-top: 30px;
  margin-bottom: 30px;
}

.section-description {
  color: var(--text-color-secondary);
  margin-bottom: 20px;
}

.metadata-field {
  display: flex;
  align-items: flex-start;
  margin-bottom: 20px;
  width: 100%;
  gap: 15px;
}

.metadata-inputs {
  display: flex;
  flex: 1;
  width: 100%;
  gap: 15px;
  flex-wrap: wrap; /* Allow inputs to wrap on small screens */
}

.metadata-key {
  flex: 1;
  min-width: 180px; /* Ensure reasonable width before wrapping */
}

.metadata-value {
  flex: 1;
  min-width: 180px; /* Ensure reasonable width before wrapping */
}

.remove-metadata-btn {
  background: none;
  border: none;
  color: var(--danger-color);
  cursor: pointer;
  padding: 8px;
  margin-top: 24px; /* Align with input fields */
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 36px;
  height: 36px;
  border-radius: 50%;
  transition: background-color 0.2s ease;
}

.remove-metadata-btn:hover {
  background-color: rgba(231, 76, 60, 0.1);
}

.add-metadata-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  background: none;
  border: 1px dashed var(--border-color);
  color: var(--primary-color);
  padding: 10px;
  width: 100%;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 10px;
  transition: all 0.2s ease;
}

.add-metadata-btn:hover {
  background-color: var(--primary-color-light-transparent);
}

.form-actions {
  margin-top: 30px;
  display: flex;
  justify-content: flex-end;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  color: var(--text-color);
  font-weight: 500;
}

.form-select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background-color: var(--input-background);
  color: var(--text-color);
  font-size: 14px;
}

.form-hint {
  display: block;
  margin-top: 5px;
  color: var(--text-color-secondary);
  font-size: 12px;
}

/* Make sure Material Icons font is used correctly */
.material-icons {
  font-family: 'Material Icons';
  font-weight: normal;
  font-style: normal;
  font-size: 24px;
  line-height: 1;
  letter-spacing: normal;
  text-transform: none;
  display: inline-block;
  white-space: nowrap;
  word-wrap: normal;
  direction: ltr;
  font-feature-settings: 'liga';
  -webkit-font-feature-settings: 'liga';
  -webkit-font-smoothing: antialiased;
}

.upload-layout {
  display: flex;
  gap: 20px;
  min-height: 600px;
}

.node-tree {
  margin-left: 12px; /* Add margin for child nodes */
}

.node-connector {
  position: absolute;
  left: -12px;
  top: 50%;
  width: 12px;
  height: 1px;
  background-color: var(--border-color);
}

.node-connector::before {
  content: '';
  position: absolute;
  left: 0;
  top: -12px;
  width: 1px;
  height: 24px;
  background-color: var(--border-color);
}

.folder-node {
  position: relative;
  padding-left: 12px;
}

.selected-node-path {
  font-size: 0.85em;
  color: var(--text-color-secondary);
  margin-top: 4px;
}

.selected-node-path .material-icons {
  font-size: 16px;
  vertical-align: middle;
  margin: 0 4px;
}

.upload-content {
  flex: 1;
  padding: 20px;
  background: var(--surface-color);
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.selected-folder {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 4px;
  color: var(--text-color-secondary);
  font-size: 0.9em;
}

.selected-folder .material-icons {
  font-size: 16px;
  color: var(--folder-color);
}

.selected-folder-path {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background-color: var(--background-color);
  border-radius: 4px;
  margin-bottom: 20px;
  color: var(--text-color);
}

.path-segments {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.9rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.path-separator {
  color: var(--text-color-secondary);
  margin: 0 4px;
}
</style>