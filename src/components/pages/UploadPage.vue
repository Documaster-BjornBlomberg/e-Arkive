<template>
  <MainLayout>
    <div class="upload-container">
      <h2>Ladda upp dokument</h2>
      
      <!-- Add file input section -->
      <div class="file-selection">
        <FileInput
          id="document-upload"
          label="Välj dokument att ladda upp"
          @change="onFileChange"
          :disabled="isUploading"
        />
        
        <!-- Show selected file info -->
        <div v-if="file" class="file-info">
          <strong>Vald fil:</strong> {{ file.name }} ({{ formatFileSize(file.size) }})
        </div>
      </div>
      
      <MetadataSection
        :metadata="metadataList"
        :disabled="isUploading"
        @update:metadata="metadataList = $event"
        @add-metadata="addMetadata"
      />
      
      <div class="form-actions">
        <Button
          variant="primary"
          :disabled="!file || isUploading"
          @click="uploadFile"
        >
          {{ isUploading ? 'Laddar upp...' : 'Ladda upp dokument' }}
        </Button>
        <Button
          variant="default"
          :disabled="isUploading"
          @click="resetForm"
        >
          Återställ
        </Button>
      </div>

      <!-- Status message -->
      <div v-if="uploadStatus.show" :class="['status-message', uploadStatus.success ? 'success' : 'error']">
        {{ uploadStatus.message }}
      </div>
    </div>
  </MainLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import MetadataSection from '../organisms/MetadataSection.vue';
import Button from '../atoms/Button.vue';
import MainLayout from '../templates/MainLayout.vue';
import FileInput from '../atoms/FileInput.vue';
import { useGraphQL } from '../../composables/useGraphQL';
import StatusMessage from '../atoms/StatusMessage.vue';

const file = ref(null);
const metadataList = ref([{ key: '', value: '' }]);
const isUploading = ref(false);
const uploadStatus = ref({ show: false, success: false, message: '' });

// Get GraphQL utility
const { saveFile, loading, error } = useGraphQL();

// Handle file selection
const onFileChange = (e) => {
  file.value = e.target.files[0];
  // Reset status message when selecting a new file
  uploadStatus.value = { show: false, success: false, message: '' };
};

const addMetadata = () => {
  metadataList.value.push({ key: '', value: '' });
};

const uploadFile = async () => {
  if (!file.value) return;
  
  isUploading.value = true;
  uploadStatus.value = { show: false, success: false, message: '' };
  
  try {
    // Convert file to base64
    const fileData = await readFileAsBase64(file.value);
    
    // Prepare payload - only include metadata entries that have both key and value
    const validMetadata = metadataList.value.filter(m => m.key && m.value);
    
    const payload = {
      name: file.value.name,
      size: file.value.size,
      contentType: file.value.type,
      fileData: fileData,
      metadata: validMetadata
    };
    
    console.log("Uploading file:", payload.name);
    
    // Upload file using GraphQL mutation
    const result = await saveFile(payload);
    
    if (result && result.id) {
      console.log("File uploaded successfully with ID:", result.id);
      uploadStatus.value = { 
        show: true, 
        success: true, 
        message: `Filen "${result.name}" har laddats upp framgångsrikt.` 
      };
      resetForm();
    } else {
      throw new Error("Filuppladdning misslyckades: Inget svar från servern");
    }
    
  } catch (error) {
    console.error("Error uploading file:", error);
    uploadStatus.value = { 
      show: true, 
      success: false, 
      message: `Filuppladdning misslyckades: ${error.message}` 
    };
  } finally {
    isUploading.value = false;
  }
};

const resetForm = () => {
  file.value = null;
  metadataList.value = [{ key: '', value: '' }];
  
  // Reset file input by clearing the value
  const fileInput = document.getElementById('document-upload');
  if (fileInput) fileInput.value = '';
};

// Helper function to read file as base64
const readFileAsBase64 = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => {
      // Remove data URL prefix (e.g., "data:application/pdf;base64,")
      const base64 = reader.result.split(',')[1];
      resolve(base64);
    };
    reader.onerror = error => reject(error);
  });
};

// Helper function to format file size
const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
};

onMounted(() => {
  // Initialization logic here
});
</script>

<style scoped>
.upload-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background-color: var(--background-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
}

.upload-container h2 {
  margin-top: 0;
  color: var(--text-color);
}

.file-selection {
  margin-bottom: 20px;
}

.file-info {
  margin-top: 10px;
  padding: 10px;
  background-color: rgba(128, 128, 128, 0.1);
  border-radius: 4px;
  color: var(--text-color);
}

.form-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

/* Add status message styling */
.status-message {
  margin-top: 20px;
  padding: 12px;
  border-radius: 4px;
}

.status-message.success {
  background-color: rgba(46, 204, 113, 0.1);
  border-left: 4px solid #2ecc71;
  color: #27ae60;
}

.status-message.error {
  background-color: rgba(231, 76, 60, 0.1);
  border-left: 4px solid #e74c3c;
  color: #c0392b;
}
</style>
