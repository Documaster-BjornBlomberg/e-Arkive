<script setup>
import { ref, onMounted, reactive } from 'vue'

const files = ref([])
const editingFile = ref(null)
const editingMetadata = ref([])
const newMetadata = reactive({ key: '', value: '' })

const fetchFiles = async () => {
  const query = `
    query {
      getFiles {
        id
        name
        size
        contentType
        createdAt
        metadata {
          key
          value
        }
      }
    }
  `

  try {
    const response = await fetch('http://localhost:8080/query', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query }),
    })

    const result = await response.json()
    if (result.errors) {
      throw new Error(result.errors[0].message)
    }

    files.value = result.data.getFiles
  } catch (error) {
    alert('Error fetching files: ' + error.message)
  }
}

const startEditMetadata = (file) => {
  editingFile.value = file.id
  // Create a deep copy of the metadata to edit
  editingMetadata.value = file.metadata ? 
    file.metadata.map(meta => ({ key: meta.key, value: meta.value })) : 
    []
}

const cancelEditMetadata = () => {
  editingFile.value = null
  editingMetadata.value = []
}

const addMetadataField = () => {
  editingMetadata.value.push({ key: '', value: '' })
}

const removeMetadataField = (index) => {
  editingMetadata.value.splice(index, 1)
}

const saveMetadata = async () => {
  if (!editingFile.value) return

  // Filter out empty metadata entries
  const validMetadata = editingMetadata.value.filter(meta => meta.key.trim() !== '' && meta.value.trim() !== '')

  const query = `
    mutation($fileId: ID!, $metadataInput: [MetadataInput!]!) {
      updateMetadata(fileId: $fileId, metadataInput: $metadataInput) {
        id
        metadata {
          key
          value
        }
      }
    }
  `

  const variables = {
    fileId: editingFile.value,
    metadataInput: validMetadata
  }

  try {
    const response = await fetch('http://localhost:8080/query', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query, variables }),
    })

    const result = await response.json()
    if (result.errors) {
      throw new Error(result.errors[0].message)
    }

    // Update the file in the local state
    const updatedFile = files.value.find(f => f.id === editingFile.value)
    if (updatedFile) {
      updatedFile.metadata = result.data.updateMetadata.metadata
    }

    // Reset editing state
    editingFile.value = null
    editingMetadata.value = []
    
    alert('Metadata updated successfully!')
  } catch (error) {
    alert('Error updating metadata: ' + error.message)
  }
}

const deleteMetadataKey = async (fileId, key) => {
  if (!confirm(`Are you sure you want to delete the metadata with key "${key}"?`)) {
    return
  }

  const query = `
    mutation($fileId: ID!, $keys: [String!]!) {
      deleteMetadata(fileId: $fileId, keys: $keys) {
        id
        metadata {
          key
          value
        }
      }
    }
  `

  const variables = {
    fileId,
    keys: [key]
  }

  try {
    const response = await fetch('http://localhost:8080/query', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query, variables }),
    })

    const result = await response.json()
    if (result.errors) {
      throw new Error(result.errors[0].message)
    }

    // Update the file in the local state
    const updatedFile = files.value.find(f => f.id === fileId)
    if (updatedFile) {
      updatedFile.metadata = result.data.deleteMetadata.metadata
    }
    
    alert('Metadata deleted successfully!')
  } catch (error) {
    alert('Error deleting metadata: ' + error.message)
  }
}

onMounted(fetchFiles)
</script>

<template>
  <div>
    <h1>Document List</h1>
    <table class="file-table">
      <thead>
        <tr>
          <th>Name</th>
          <th>Size</th>
          <th>Content Type</th>
          <th>Created At</th>
          <th>Metadata</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="file in files" :key="file.id">
          <td>{{ file.name }}</td>
          <td>{{ file.size }}</td>
          <td>{{ file.contentType }}</td>
          <td>{{ file.createdAt }}</td>
          <td>
            <div v-if="editingFile === file.id" class="metadata-editor">
              <div v-for="(meta, index) in editingMetadata" :key="index" class="metadata-edit-row">
                <input type="text" v-model="meta.key" placeholder="Key" />
                <input type="text" v-model="meta.value" placeholder="Value" />
                <button @click="removeMetadataField(index)" class="delete-btn">Remove</button>
              </div>
              <div class="metadata-actions">
                <button @click="addMetadataField" class="add-btn">Add Field</button>
                <button @click="saveMetadata" class="save-btn">Save</button>
                <button @click="cancelEditMetadata" class="cancel-btn">Cancel</button>
              </div>
            </div>
            <div v-else>
              <ul v-if="file.metadata && file.metadata.length > 0">
                <li v-for="meta in file.metadata" :key="meta.key" class="metadata-item">
                  <span>{{ meta.key }}: {{ meta.value }}</span>
                  <button @click="deleteMetadataKey(file.id, meta.key)" class="delete-meta-btn">Delete</button>
                </li>
              </ul>
              <p v-else>No metadata</p>
            </div>
          </td>
          <td>
            <button v-if="editingFile !== file.id" @click="startEditMetadata(file)" class="edit-btn">Edit Metadata</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.file-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

.file-table th, .file-table td {
  border: 1px solid var(--table-border-color);
  padding: 8px;
  text-align: left;
  vertical-align: top;
}

.file-table th {
  background-color: var(--header-bg);
  font-weight: bold;
}

.file-table tr:nth-child(even) {
  background-color: rgba(128, 128, 128, 0.05);
}

.file-table tr:hover {
  background-color: rgba(128, 128, 128, 0.1);
}

.metadata-editor {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.metadata-edit-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.metadata-edit-row input {
  padding: 4px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background-color: var(--background-color);
  color: var(--text-color);
}

.metadata-actions {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

button {
  padding: 4px 8px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  background-color: var(--button-bg);
  color: var(--button-text);
}

button:hover {
  background-color: var(--button-hover-bg);
}

.delete-btn, .delete-meta-btn {
  background-color: #e74c3c;
}

.delete-btn:hover, .delete-meta-btn:hover {
  background-color: #c0392b;
}

.save-btn {
  background-color: #2ecc71;
}

.save-btn:hover {
  background-color: #27ae60;
}

.cancel-btn {
  background-color: #7f8c8d;
}

.cancel-btn:hover {
  background-color: #95a5a6;
}

.add-btn {
  background-color: #3498db;
}

.add-btn:hover {
  background-color: #2980b9;
}

.edit-btn {
  background-color: #f39c12;
}

.edit-btn:hover {
  background-color: #d35400;
}

.metadata-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.delete-meta-btn {
  font-size: 0.8em;
  padding: 2px 4px;
}
</style>
