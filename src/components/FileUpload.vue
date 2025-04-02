<script setup>
import { ref } from 'vue'

const file = ref(null)
const metadataList = ref([{ key: '', value: '' }])
const uploadStatus = ref('')
const isUploading = ref(false)

const addMetadata = () => {
  metadataList.value.push({ key: '', value: '' })
}

const removeMetadata = (index) => {
  metadataList.value.splice(index, 1)
}

const resetForm = () => {
  file.value = null
  metadataList.value = [{ key: '', value: '' }]
  uploadStatus.value = ''
  // Reset the file input
  const fileInput = document.querySelector('input[type="file"]')
  if (fileInput) fileInput.value = ''
}

const uploadFile = async () => {
  if (!file.value) {
    alert('Please select a file to upload')
    return
  }

  isUploading.value = true
  uploadStatus.value = 'Reading file...'

  const reader = new FileReader()
  reader.onload = async (e) => {
    try {
      const base64Data = e.target.result.split(',')[1]

      // Filter out empty metadata entries
      const metadata = metadataList.value.filter(meta => meta.key.trim() && meta.value.trim())

      uploadStatus.value = 'Preparing upload...'

      const fileInput = {
        name: file.value.name,
        size: file.value.size,
        contentType: file.value.type,
        fileData: base64Data,
        metadata,
      }

      const query = `
        mutation($input: FileInput!) {
          saveFile(input: $input) {
            id
            name
          }
        }
      `

      const variables = { input: fileInput }

      uploadStatus.value = 'Uploading to server...'

      const response = await fetch('http://localhost:8080/query', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ query, variables }),
      })

      const result = await response.json()
      if (result.errors) {
        throw new Error(result.errors[0].message)
      }

      uploadStatus.value = 'Upload successful!'
      setTimeout(() => {
        resetForm()
      }, 3000)
    } catch (error) {
      uploadStatus.value = `Error: ${error.message}`
    } finally {
      isUploading.value = false
    }
  }

  reader.onerror = () => {
    uploadStatus.value = 'Error reading file'
    isUploading.value = false
  }

  reader.readAsDataURL(file.value)
}

const onFileChange = (e) => {
  file.value = e.target.files[0]
}
</script>

<template>
  <div class="upload-container">
    <h1>Upload a Document</h1>
    
    <div class="upload-form">
      <div class="file-input-container">
        <label for="file-upload" class="file-label">Select File:</label>
        <input 
          id="file-upload" 
          type="file" 
          @change="onFileChange" 
          :disabled="isUploading"
        />
      </div>

      <div class="file-info" v-if="file">
        <p><strong>Selected File:</strong> {{ file.name }}</p>
        <p><strong>Size:</strong> {{ Math.round(file.size / 1024) }} KB</p>
        <p><strong>Type:</strong> {{ file.type }}</p>
      </div>

      <div class="metadata-section">
        <h2>Document Metadata</h2>
        <p class="metadata-help">Add key-value pairs to describe your document</p>
        
        <div class="metadata-list">
          <div 
            v-for="(meta, index) in metadataList" 
            :key="index"
            class="metadata-item"
          >
            <div class="metadata-inputs">
              <div class="input-group">
                <label>Key:</label>
                <input 
                  type="text" 
                  v-model="meta.key" 
                  placeholder="e.g. Author, Date, Category" 
                  :disabled="isUploading"
                />
              </div>
              
              <div class="input-group">
                <label>Value:</label>
                <input 
                  type="text" 
                  v-model="meta.value" 
                  placeholder="e.g. John Doe, 2025-04-02" 
                  :disabled="isUploading"
                />
              </div>
            </div>
            
            <button 
              @click="removeMetadata(index)" 
              class="remove-btn"
              :disabled="metadataList.length === 1 || isUploading"
              title="Remove this metadata"
            >
              ✕
            </button>
          </div>
        </div>
        
        <div class="metadata-actions">
          <button 
            @click="addMetadata" 
            class="add-btn"
            :disabled="isUploading"
          >
            Add More Metadata
          </button>
        </div>
      </div>

      <div class="upload-status" v-if="uploadStatus">
        <p :class="{ 'error': uploadStatus.startsWith('Error') }">{{ uploadStatus }}</p>
      </div>

      <div class="form-actions">
        <button 
          @click="uploadFile" 
          class="upload-btn"
          :disabled="!file || isUploading"
        >
          {{ isUploading ? 'Uploading...' : 'Upload Document' }}
        </button>
        
        <button 
          @click="resetForm" 
          class="reset-btn"
          :disabled="isUploading"
        >
          Reset Form
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.upload-container {
  max-width: 800px;
  margin: 0 auto;
}

.upload-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.file-input-container {
  margin-bottom: 10px;
}

.file-label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.file-info {
  background-color: rgba(128, 128, 128, 0.1);
  padding: 10px;
  border-radius: 4px;
  margin-bottom: 10px;
}

.metadata-section {
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 15px;
  background-color: rgba(128, 128, 128, 0.05);
}

.metadata-section h2 {
  margin-top: 0;
  margin-bottom: 5px;
}

.metadata-help {
  color: #666;
  margin-bottom: 15px;
  font-size: 0.9em;
}

.metadata-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.metadata-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.metadata-inputs {
  display: flex;
  flex: 1;
  gap: 10px;
}

.input-group {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.input-group label {
  margin-bottom: 3px;
  font-size: 0.9em;
}

input[type="text"], input[type="file"] {
  padding: 8px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background-color: var(--background-color);
  color: var(--text-color);
}

.metadata-actions {
  margin-top: 10px;
}

.form-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  transition: background-color 0.2s;
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.add-btn {
  background-color: #3498db;
  color: white;
}

.add-btn:hover:not(:disabled) {
  background-color: #2980b9;
}

.remove-btn {
  background-color: #e74c3c;
  color: white;
  padding: 4px 8px;
  font-size: 0.8em;
}

.remove-btn:hover:not(:disabled) {
  background-color: #c0392b;
}

.upload-btn {
  background-color: var(--button-bg);
  color: var(--button-text);
  padding: 10px 20px;
}

.upload-btn:hover:not(:disabled) {
  background-color: var(--button-hover-bg);
}

.reset-btn {
  background-color: #7f8c8d;
  color: white;
}

.reset-btn:hover:not(:disabled) {
  background-color: #95a5a6;
}

.upload-status {
  padding: 10px;
  border-radius: 4px;
  background-color: rgba(46, 204, 113, 0.1);
  border-left: 4px solid #2ecc71;
}

.upload-status .error {
  background-color: rgba(231, 76, 60, 0.1);
  border-left-color: #e74c3c;
}
</style>
