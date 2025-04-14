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
          @node-change="handleNodeChange" />
        
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
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import MainLayout from '../templates/MainLayout.vue';
import FileUploadForm from '../organisms/FileUploadForm.vue';
import Button from '../atoms/Button.vue';
import FileCardMolecule from '../molecules/FileCardMolecule.vue';
import { useGraphQL } from '../../composables/useGraphQL';
import { useNodeHandling } from '../../composables/useNodeHandling';

const { saveFile } = useGraphQL();
// Updated to use the new function for fetching all nodes
const { getAllNodesForUpload } = useNodeHandling(); 

const router = useRouter();
const uploading = ref(false);
const uploadedFile = ref(null);
const statusMessage = ref('');
const isSuccess = ref(true);
const availableNodes = ref([]);
const selectedNodeId = ref('1'); // Default to root node

// Navigate to the file list page
const navigateToList = () => {
  router.push('/');
};

// Reset the form for another upload
const resetUpload = () => {
  uploadedFile.value = null;
  statusMessage.value = '';
};

// Handle node selection change
const handleNodeChange = (nodeId) => {
  selectedNodeId.value = nodeId;
};

// Handle file upload
const handleUpload = async (fileData) => {
  uploading.value = true;
  statusMessage.value = 'Laddar upp dokument...';
  isSuccess.value = true;
  
  try {
    // Add the selected node ID to the upload data
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

// Initialize available nodes for the dropdown using the recursive fetch
const loadAvailableNodes = async () => {
  try {
    // Fetch all nodes recursively and format them for the dropdown
    availableNodes.value = await getAllNodesForUpload(); 
  } catch (error) {
    console.error('Error loading all nodes for upload:', error);
  }
};

// Load nodes when component is mounted
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
