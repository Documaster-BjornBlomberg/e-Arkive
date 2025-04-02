<script setup>
import { ref, onMounted, reactive } from 'vue'

const files = ref([])
const editingFile = ref(null)
const editingMetadata = ref([])
const isLoading = ref(true)
const errorMessage = ref('')

// Hämta alla filer från servern
const fetchFiles = async () => {
  isLoading.value = true
  errorMessage.value = ''
  
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
    errorMessage.value = 'Fel vid hämtning av filer: ' + error.message
    console.error(error)
  } finally {
    isLoading.value = false
  }
}

// Starta redigering av metadata för en fil
const startEditMetadata = (file) => {
  editingFile.value = file.id
  // Skapa en kopia av filens metadata för redigering
  editingMetadata.value = file.metadata 
    ? file.metadata.map(meta => ({ key: meta.key, value: meta.value })) 
    : []
  
  // Lägg till ett tomt metadatafält om det inte finns några
  if (editingMetadata.value.length === 0) {
    addMetadataField()
  }
}

// Avbryt redigering
const cancelEditMetadata = () => {
  editingFile.value = null
  editingMetadata.value = []
}

// Lägg till ett nytt tomt metadatafält
const addMetadataField = () => {
  editingMetadata.value.push({ key: '', value: '' })
}

// Ta bort ett metadatafält
const removeMetadataField = (index) => {
  editingMetadata.value.splice(index, 1)
}

// Spara uppdaterade metadata
const saveMetadata = async () => {
  if (!editingFile.value) return

  // Filtrera bort tomma metadata-poster
  const validMetadata = editingMetadata.value.filter(meta => 
    meta.key.trim() !== '' && meta.value.trim() !== ''
  )

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

    // Uppdatera filen i lokal state
    const updatedFile = files.value.find(f => f.id === editingFile.value)
    if (updatedFile) {
      updatedFile.metadata = result.data.updateMetadata.metadata
    }

    // Återställ redigeringsläge
    editingFile.value = null
    editingMetadata.value = []
    
    alert('Metadata uppdaterad!')
  } catch (error) {
    alert('Fel vid uppdatering av metadata: ' + error.message)
  }
}

// Ta bort en specifik metadata
const deleteMetadataKey = async (fileId, key) => {
  if (!confirm(`Är du säker på att du vill ta bort metadata med nyckeln "${key}"?`)) {
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

    // Uppdatera filen i lokal state
    const updatedFile = files.value.find(f => f.id === fileId)
    if (updatedFile) {
      updatedFile.metadata = result.data.deleteMetadata.metadata
    }
    
    alert('Metadata borttagen!')
  } catch (error) {
    alert('Fel vid borttagning av metadata: ' + error.message)
  }
}

// Ta bort fil
const deleteFile = async (fileId) => {
  if (!confirm('Är du säker på att du vill ta bort denna fil?')) {
    return
  }

  const query = `
    mutation DeleteFile($id: ID!) {
      deleteFile(id: $id)
    }
  `

  const variables = { id: fileId }

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

    if (result.data.deleteFile) {
      // Ta bort filen från lokal state
      files.value = files.value.filter(f => f.id !== fileId)
      alert('Fil borttagen!')
    } else {
      throw new Error('Kunde inte ta bort filen')
    }
  } catch (error) {
    alert('Fel vid borttagning av fil: ' + error.message)
  }
}

// Ladda ner en fil
const downloadFile = async (fileId) => {
  const query = `
    query DownloadFile($id: ID!) {
      downloadFile(id: $id) {
        id
        name
        size
        contentType
        fileData
      }
    }
  `

  const variables = { id: fileId }

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

    const file = result.data.downloadFile
    if (!file || !file.fileData) {
      throw new Error('Fildata saknas')
    }

    // Skapa en nedladdningslänk
    const link = document.createElement('a')
    link.href = `data:${file.contentType};base64,${file.fileData}`
    link.download = file.name
    link.click()
  } catch (error) {
    alert('Fel vid nedladdning av fil: ' + error.message)
  }
}

// Formatera filstorlek
const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
}

// Hämta filer när komponenten laddas
onMounted(fetchFiles)
</script>

<template>
  <div class="document-list">
    <h1>Dokumentlista</h1>
    
    <div v-if="isLoading" class="loading">
      <p>Laddar dokument...</p>
    </div>
    
    <div v-else-if="errorMessage" class="error-message">
      <p>{{ errorMessage }}</p>
      <button @click="fetchFiles" class="retry-btn">Försök igen</button>
    </div>
    
    <div v-else-if="files.length === 0" class="no-files">
      <p>Inga dokument hittades.</p>
    </div>
    
    <div v-else class="file-list">
      <button @click="fetchFiles" class="refresh-btn">Uppdatera listan</button>
      
      <table class="file-table">
        <thead>
          <tr>
            <th>Namn</th>
            <th>Storlek</th>
            <th>Typ</th>
            <th>Skapad</th>
            <th>Metadata</th>
            <th>Åtgärder</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="file in files" :key="file.id">
            <td>{{ file.name }}</td>
            <td>{{ formatFileSize(file.size) }}</td>
            <td>{{ file.contentType }}</td>
            <td>{{ new Date(file.createdAt).toLocaleString() }}</td>
            <td class="metadata-cell">
              <div v-if="editingFile === file.id" class="metadata-editor">
                <div v-for="(meta, index) in editingMetadata" :key="index" class="metadata-field">
                  <div class="metadata-inputs">
                    <div class="input-group">
                      <label>Nyckel:</label>
                      <input type="text" v-model="meta.key" />
                    </div>
                    <div class="input-group">
                      <label>Värde:</label>
                      <input type="text" v-model="meta.value" />
                    </div>
                    <button @click="removeMetadataField(index)" class="delete-btn" title="Ta bort detta fält">🗑️</button>
                  </div>
                </div>
                
                <div class="metadata-actions">
                  <button @click="addMetadataField" class="add-btn" title="Lägg till metadata">
                    <span>+ Lägg till fält</span>
                  </button>
                  <button @click="saveMetadata" class="save-btn" title="Spara ändringar">
                    <span>✓ Spara</span>
                  </button>
                  <button @click="cancelEditMetadata" class="cancel-btn" title="Avbryt">
                    <span>✕ Avbryt</span>
                  </button>
                </div>
              </div>
              <div v-else class="metadata-display">
                <div v-if="file.metadata && file.metadata.length > 0" class="metadata-list">
                  <div v-for="meta in file.metadata" :key="meta.key" class="metadata-item">
                    <span class="metadata-key">{{ meta.key }}:</span>
                    <span class="metadata-value">{{ meta.value }}</span>
                    <button @click="deleteMetadataKey(file.id, meta.key)" class="delete-meta-btn" title="Ta bort denna metadata">🗑️</button>
                  </div>
                </div>
                <p v-else class="no-metadata">Ingen metadata</p>
              </div>
              <button v-if="editingFile !== file.id" @click="startEditMetadata(file)" class="edit-btn" title="Redigera metadata">
                ✏️ Redigera metadata
              </button>
            </td>
            <td class="actions-cell">
              <div class="action-buttons">
                <button @click="downloadFile(file.id)" class="download-btn" title="Ladda ner">
                  📥 Ladda ner
                </button>
                <button @click="deleteFile(file.id)" class="delete-file-btn" title="Ta bort fil">
                  🗑️ Ta bort
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.document-list {
  padding: 20px;
}

.loading, .error-message, .no-files {
  text-align: center;
  margin: 40px 0;
}

.error-message {
  color: #e74c3c;
}

.refresh-btn {
  margin-bottom: 20px;
}

.file-table {
  width: 100%;
  border-collapse: collapse;
}

.file-table th, .file-table td {
  border: 1px solid var(--table-border-color);
  padding: 10px;
}

.file-table th {
  background-color: var(--header-bg);
  font-weight: bold;
  text-align: left;
}

.metadata-cell {
  min-width: 250px;
  max-width: 400px;
}

.metadata-display {
  max-height: 200px;
  overflow-y: auto;
}

.metadata-list {
  margin-bottom: 8px;
}

.metadata-item {
  display: flex;
  align-items: center;
  margin-bottom: 5px;
  padding: 3px;
  border-bottom: 1px solid var(--border-color);
}

.metadata-key {
  font-weight: bold;
  margin-right: 5px;
}

.metadata-value {
  flex-grow: 1;
  margin-right: 10px;
  word-break: break-all;
}

.delete-meta-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1rem;
  padding: 3px;
}

.delete-meta-btn:hover {
  color: #e74c3c;
}

.no-metadata {
  font-style: italic;
  color: #7f8c8d;
  margin-bottom: 8px;
}

.metadata-editor {
  border: 1px solid var(--border-color);
  padding: 10px;
  border-radius: 4px;
  background-color: var(--background-color);
}

.metadata-field {
  margin-bottom: 10px;
}

.metadata-inputs {
  display: flex;
  gap: 10px;
  align-items: center;
}

.input-group {
  flex-grow: 1;
}

.input-group label {
  display: block;
  margin-bottom: 3px;
  font-size: 0.8rem;
}

.input-group input {
  width: 100%;
  padding: 5px;
  border: 1px solid var(--border-color);
  border-radius: 3px;
  background-color: var(--background-color);
  color: var(--text-color);
}

.metadata-actions {
  display: flex;
  gap: 10px;
  margin-top: 15px;
}

.add-btn, .save-btn, .cancel-btn, .delete-btn {
  display: flex;
  align-items: center;
  gap: 5px;
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.download-btn, .delete-file-btn, .edit-btn {
  width: 100%;
  text-align: left;
  display: flex;
  align-items: center;
  gap: 5px;
}

.edit-btn {
  margin-top: 8px;
  background-color: #f39c12;
}

.delete-file-btn:hover {
  background-color: #c0392b;
}

.delete-file-btn {
  background-color: #e74c3c;
}

.edit-btn:hover {
  background-color: #d35400;
}

.download-btn {
  background-color: #3498db;
}

.download-btn:hover {
  background-color: #2980b9;
}

@media (max-width: 768px) {
  .file-table {
    font-size: 0.9rem;
  }
  
  .metadata-inputs {
    flex-direction: column;
    gap: 5px;
  }
}
</style>
