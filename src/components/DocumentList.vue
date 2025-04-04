<script setup>
import { ref, computed, onMounted, reactive, watch } from 'vue'

// Tillstånd
const files = ref([])
const editingFile = ref(null)
const editingMetadata = ref([])
const isLoading = ref(true)
const errorMessage = ref('')
const selectedFile = ref(null)
const sidepanelOpen = ref(false)
const activeMetadataTab = ref('list')
const viewMode = ref('sidepanel') // 'sidepanel' eller 'expandable'
const expandedFileId = ref(null)
const metadataSearch = ref('')

// Filtrerad metadata baserad på sökning
const filteredMetadata = computed(() => {
  if (!selectedFile.value || !selectedFile.value.metadata) return [];
  
  if (!metadataSearch.value.trim()) return selectedFile.value.metadata;
  
  const search = metadataSearch.value.toLowerCase();
  return selectedFile.value.metadata.filter(meta => 
    meta.key.toLowerCase().includes(search) || 
    meta.value.toLowerCase().includes(search)
  );
})

// Grupperad metadata för kategorivyn
const groupedMetadata = computed(() => {
  if (!selectedFile.value || !selectedFile.value.metadata) return {};

  const groups = {};
  selectedFile.value.metadata.forEach(meta => {
    const parts = meta.key.split('.');
    const category = parts.length > 1 ? parts[0] : 'Övrigt';
    
    if (!groups[category]) {
      groups[category] = [];
    }
    
    groups[category].push({
      key: parts.length > 1 ? parts.slice(1).join('.') : meta.key,
      originalKey: meta.key,
      value: meta.value
    });
  });
  
  return groups;
})

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

// Öppna sidopanel för en fil
const openSidepanel = (file) => {
  selectedFile.value = file
  sidepanelOpen.value = true
}

// Stäng sidopanel
const closeSidepanel = () => {
  sidepanelOpen.value = false
  setTimeout(() => {
    selectedFile.value = null
  }, 300) // Vänta på animation innan filen nollställs
}

// Växla expandering av en fil i expanderbart läge
const toggleExpanded = (fileId) => {
  expandedFileId.value = expandedFileId.value === fileId ? null : fileId
}

// Växla visningsläge mellan sidopanel och expanderbart
const toggleViewMode = () => {
  viewMode.value = viewMode.value === 'sidepanel' ? 'expandable' : 'sidepanel'
  // Återställ tillstånd när läge ändras
  expandedFileId.value = null
  closeSidepanel()
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
      
      // Om den uppdaterade filen är den för tillfället valda filen i sidopanelen
      if (selectedFile.value && selectedFile.value.id === updatedFile.id) {
        selectedFile.value = { ...updatedFile }
      }
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
      
      // Om den uppdaterade filen är den för tillfället valda filen i sidopanelen
      if (selectedFile.value && selectedFile.value.id === updatedFile.id) {
        selectedFile.value = { ...updatedFile }
      }
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
      // Om den borttagna filen är den valda filen i sidopanelen
      if (selectedFile.value && selectedFile.value.id === fileId) {
        closeSidepanel()
      }
      
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

// Lista av metadata-fliktabbar
const metadataTabs = [
  { id: 'list', label: 'Lista' },
  { id: 'cards', label: 'Kort' },
  { id: 'categories', label: 'Kategorier' }
]

// Hämta filer när komponenten laddas
onMounted(() => {
  fetchFiles()
  // Sätt standardläget för metadata-flikarna
  activeMetadataTab.value = 'list'
})
</script>

<template>
  <div class="document-list">
    <div class="list-header">
      <h1>Dokumentlista</h1>
      <div class="view-controls">
        <button @click="toggleViewMode" class="view-toggle-btn" :title="viewMode === 'sidepanel' ? 'Byt till expanderbar vy' : 'Byt till sidopanel vy'">
          {{ viewMode === 'sidepanel' ? '⬍ Expanderbar vy' : '⬌ Sidopanel vy' }}
        </button>
        <button @click="fetchFiles" class="refresh-btn" title="Uppdatera listan">
          🔄 Uppdatera
        </button>
      </div>
    </div>
    
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
    
    <!-- Sidopanel-vy (standardvy) -->
    <div v-else-if="viewMode === 'sidepanel'" class="file-list-sidepanel">
      <!-- Huvudinnehåll -->
      <div class="main-content" :class="{ 'sidebar-open': sidepanelOpen }">
        <div class="files-container">
          <div 
            v-for="file in files" 
            :key="file.id" 
            class="file-card"
            :class="{ 'selected': selectedFile && selectedFile.id === file.id }"
            @click="openSidepanel(file)"
          >
            <h3 class="file-name">{{ file.name }}</h3>
            <div class="file-info-summary">
              <div class="file-size">{{ formatFileSize(file.size) }}</div>
              <div class="file-type">{{ file.contentType }}</div>
            </div>
            <div class="file-metadata-summary">
              <span class="metadata-count">{{ file.metadata ? file.metadata.length : 0 }} metadatafält</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Sidopanelen -->
      <div class="detail-sidepanel" :class="{ 'open': sidepanelOpen }">
        <div class="sidepanel-minimize-strip" @click="closeSidepanel" title="Minimera sidopanel">
          <div class="minimize-icon">▶</div>
        </div>
        <div v-if="selectedFile" class="sidepanel-content">
          <div class="sidepanel-header">
            <h2>{{ selectedFile.name }}</h2>
          </div>
          
          <div class="file-details">
            <div class="detail-row">
              <span class="detail-label">Storlek:</span>
              <span class="detail-value">{{ formatFileSize(selectedFile.size) }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Typ:</span>
              <span class="detail-value">{{ selectedFile.contentType }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Skapad:</span>
              <span class="detail-value">{{ new Date(selectedFile.createdAt).toLocaleString() }}</span>
            </div>
          </div>
          
          <!-- Metadata-sektion -->
          <div class="metadata-section">
            <div class="section-header">
              <h3>Metadata ({{ selectedFile.metadata ? selectedFile.metadata.length : 0 }})</h3>
              <input 
                type="text" 
                v-model="metadataSearch" 
                placeholder="Sök metadata..." 
                class="metadata-search"
              />
            </div>
            
            <div class="metadata-tabs">
              <button 
                v-for="tab in metadataTabs" 
                :key="tab.id" 
                @click="activeMetadataTab = tab.id"
                :class="{ 'active': activeMetadataTab === tab.id }"
                class="tab-btn"
              >
                {{ tab.label }}
              </button>
            </div>
            
            <div class="metadata-content">
              <!-- Lista -->
              <div v-if="activeMetadataTab === 'list'" class="metadata-list-view">
                <div v-if="filteredMetadata.length > 0">
                  <div v-for="meta in filteredMetadata" :key="meta.key" class="metadata-list-item">
                    <div class="metadata-key">{{ meta.key }}:</div>
                    <div class="metadata-value">{{ meta.value }}</div>
                    <button 
                      @click.stop="deleteMetadataKey(selectedFile.id, meta.key)" 
                      class="delete-meta-btn" 
                      title="Ta bort metadata"
                    >
                      🗑️
                    </button>
                  </div>
                </div>
                <p v-else class="no-metadata">{{ metadataSearch ? 'Inga matchande metadata' : 'Ingen metadata' }}</p>
              </div>
              
              <!-- Kort -->
              <div v-else-if="activeMetadataTab === 'cards'" class="metadata-cards-view">
                <div v-if="filteredMetadata.length > 0" class="metadata-cards">
                  <div v-for="meta in filteredMetadata" :key="meta.key" class="metadata-card">
                    <div class="card-key">{{ meta.key }}</div>
                    <div class="card-value">{{ meta.value }}</div>
                    <button 
                      @click.stop="deleteMetadataKey(selectedFile.id, meta.key)" 
                      class="delete-meta-btn" 
                      title="Ta bort metadata"
                    >
                      🗑️
                    </button>
                  </div>
                </div>
                <p v-else class="no-metadata">{{ metadataSearch ? 'Inga matchande metadata' : 'Ingen metadata' }}</p>
              </div>
              
              <!-- Kategorier -->
              <div v-else-if="activeMetadataTab === 'categories'" class="metadata-categories-view">
                <div v-if="Object.keys(groupedMetadata).length > 0">
                  <div v-for="(group, category) in groupedMetadata" :key="category" class="metadata-category">
                    <h4 class="category-title">{{ category }}</h4>
                    <div v-for="meta in group" :key="meta.originalKey" class="category-item">
                      <div class="category-item-key">{{ meta.key }}:</div>
                      <div class="category-item-value">{{ meta.value }}</div>
                      <button 
                        @click.stop="deleteMetadataKey(selectedFile.id, meta.originalKey)" 
                        class="delete-meta-btn" 
                        title="Ta bort metadata"
                      >
                        🗑️
                      </button>
                    </div>
                  </div>
                </div>
                <p v-else class="no-metadata">{{ metadataSearch ? 'Inga matchande metadata' : 'Ingen metadata' }}</p>
              </div>
            </div>
          </div>
          
          <div class="sidepanel-actions">
            <button @click="downloadFile(selectedFile.id)" class="download-btn" title="Ladda ner">
              📥 Ladda ner
            </button>
            <button @click="startEditMetadata(selectedFile)" class="edit-btn" title="Redigera metadata">
              ✏️ Redigera metadata
            </button>
            <button @click="deleteFile(selectedFile.id)" class="delete-file-btn" title="Ta bort fil">
              🗑️ Ta bort
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Expanderbar vy -->
    <div v-else class="file-list-expandable">
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
              <div class="metadata-summary" @click="toggleExpanded(file.id)" title="Visa/dölj metadata">
                <span>{{ file.metadata ? file.metadata.length : 0 }} metadatafält</span>
                <span class="toggle-icon">{{ expandedFileId === file.id ? '▼' : '►' }}</span>
              </div>
              
              <div v-if="expandedFileId === file.id" class="metadata-expanded">
                <div v-if="file.metadata && file.metadata.length > 0" class="metadata-grid">
                  <div v-for="meta in file.metadata" :key="meta.key" class="metadata-grid-item">
                    <span class="metadata-key">{{ meta.key }}:</span>
                    <span class="metadata-value">{{ meta.value }}</span>
                    <button @click="deleteMetadataKey(file.id, meta.key)" class="delete-meta-btn" title="Ta bort denna metadata">🗑️</button>
                  </div>
                </div>
                <p v-else class="no-metadata">Ingen metadata</p>
                <button @click="startEditMetadata(file)" class="edit-btn" title="Redigera metadata">
                  ✏️ Redigera metadata
                </button>
              </div>
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
  
  <!-- Redigeringsmodal för metadata -->
  <div v-if="editingFile" class="modal-overlay" @click="cancelEditMetadata">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h2>Redigera metadata</h2>
        <button @click="cancelEditMetadata" class="close-btn" title="Stäng">×</button>
      </div>
      
      <div class="metadata-editor">
        <div v-for="(meta, index) in editingMetadata" :key="index" class="metadata-field">
          <div class="metadata-inputs">
            <div class="input-group">
              <label>Nyckel:</label>
              <input type="text" v-model="meta.key" placeholder="t.ex. författare, version" />
            </div>
            <div class="input-group">
              <label>Värde:</label>
              <input type="text" v-model="meta.value" placeholder="t.ex. Björn Blomberg, 1.0" />
            </div>
            <button @click="removeMetadataField(index)" class="delete-btn" title="Ta bort detta fält">🗑️</button>
          </div>
        </div>
        
        <div class="metadata-actions">
          <button @click="addMetadataField" class="add-btn" title="Lägg till metadata">
            <span>+ Lägg till fält</span>
          </button>
        </div>
      </div>
      
      <div class="modal-footer">
        <button @click="saveMetadata" class="save-btn" title="Spara ändringar">
          <span>✓ Spara</span>
        </button>
        <button @click="cancelEditMetadata" class="cancel-btn" title="Avbryt">
          <span>✕ Avbryt</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.document-list {
  padding: 0;
  position: relative;
  height: calc(100vh - 200px);
  min-height: 500px;
  overflow: hidden;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.view-controls {
  display: flex;
  gap: 10px;
}

.loading, .error-message, .no-files {
  text-align: center;
  margin: 40px 0;
}

.error-message {
  color: #e74c3c;
}

/* Sidopanel-vy */
.file-list-sidepanel {
  display: flex;
  height: 100%;
  position: relative;
}

.main-content {
  width: 100%;
  transition: width 0.3s;
  overflow-y: auto;
  height: 100%;
}

.main-content.sidebar-open {
  width: 50%;
}

.files-container {
  padding-right: 20px;
}

.file-card {
  padding: 15px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  margin-bottom: 10px;
  cursor: pointer;
  transition: all 0.2s;
  background-color: var(--background-color);
}

.file-card:hover {
  border-color: var(--button-bg);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.file-card.selected {
  border-color: var(--button-bg);
  background-color: rgba(76, 175, 80, 0.05);
}

.file-name {
  margin: 0 0 10px 0;
  font-size: 1.1rem;
  word-break: break-word;
}

.file-info-summary {
  display: flex;
  gap: 15px;
  margin-bottom: 10px;
  font-size: 0.9rem;
}

.file-metadata-summary {
  font-size: 0.85rem;
  color: #666;
}

.detail-sidepanel {
  position: absolute;
  top: 0;
  right: 0;
  height: 100%;
  width: 0;
  background-color: var(--background-color);
  border-left: 1px solid var(--border-color);
  overflow: hidden;
  transition: width 0.3s;
  box-shadow: -5px 0 15px rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.detail-sidepanel.open {
  width: 50%;
}

.sidepanel-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background-color: var(--background-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--border-color);
}

.sidepanel-header h2 {
  margin: 0;
  word-break: break-word;
}

.sidepanel-minimize-strip {
  position: absolute;
  left: 0;
  top: 0;
  width: 15px;
  height: 100%;
  background-color: var(--border-color);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 100;
  transition: background-color 0.2s;
}

.sidepanel-minimize-strip:hover {
  background-color: rgba(0, 0, 0, 0.2);
}

.minimize-icon {
  color: var(--text-color);
  font-size: 0.8rem;
  opacity: 0.6;
}

.sidepanel-minimize-strip:hover .minimize-icon {
  opacity: 1;
}

.sidepanel-content {
  width: calc(100% - 15px); /* Justera för minimera-knappen */
  margin-left: 15px; /* Justera för minimera-knappen */
  height: 100%;
  padding: 20px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.file-details {
  margin-bottom: 20px;
  padding: 10px;
  background-color: rgba(0, 0, 0, 0.02);
  border-radius: 4px;
}

.detail-row {
  display: flex;
  margin-bottom: 8px;
}

.detail-label {
  font-weight: bold;
  width: 100px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.section-header h3 {
  margin: 0;
}

.metadata-search {
  width: 180px;
  padding: 5px 10px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 0.9rem;
}

.metadata-tabs {
  display: flex;
  gap: 5px;
  margin-bottom: 15px;
  border-bottom: 1px solid var(--border-color);
}

.tab-btn {
  background: none;
  border: none;
  padding: 8px 15px;
  cursor: pointer;
  border-radius: 4px 4px 0 0;
  transition: all 0.2s;
  color: var(--text-color);
  opacity: 0.7;
}

.tab-btn.active {
  border-bottom: 2px solid var(--button-bg);
  opacity: 1;
  font-weight: bold;
}

.tab-btn:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.metadata-content {
  max-height: 300px;
  overflow-y: auto;
  padding: 10px;
  background-color: rgba(0, 0, 0, 0.02);
  border-radius: 4px;
}

.metadata-list-view {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.metadata-list-item {
  display: flex;
  align-items: center;
  padding: 5px;
  border-bottom: 1px solid var(--border-color);
}

.metadata-key {
  font-weight: bold;
  margin-right: 5px;
  min-width: 100px;
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
  opacity: 0.6;
}

.delete-meta-btn:hover {
  opacity: 1;
  color: #e74c3c;
}

.no-metadata {
  font-style: italic;
  color: #7f8c8d;
  text-align: center;
  padding: 20px;
}

.metadata-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 10px;
}

.metadata-card {
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 10px;
  position: relative;
  background-color: var(--background-color);
}

.card-key {
  font-weight: bold;
  font-size: 0.9em;
  margin-bottom: 5px;
  color: var(--button-bg);
}

.card-value {
  word-break: break-all;
}

.metadata-card .delete-meta-btn {
  position: absolute;
  top: 5px;
  right: 5px;
}

.metadata-category {
  margin-bottom: 15px;
}

.category-title {
  margin: 0 0 8px 0;
  font-size: 1em;
  color: var(--button-bg);
  padding-bottom: 3px;
  border-bottom: 1px solid var(--border-color);
}

.category-item {
  display: flex;
  align-items: center;
  padding: 3px;
  margin-left: 15px;
}

.category-item-key {
  font-weight: bold;
  margin-right: 5px;
}

.category-item-value {
  flex-grow: 1;
  margin-right: 10px;
  word-break: break-all;
}

.sidepanel-actions {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.sidepanel-actions button {
  width: 100%;
}

/* Knappstilar */
.download-btn, .edit-btn, .delete-file-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
}

.view-toggle-btn, .refresh-btn {
  padding: 8px 16px;
}

.download-btn {
  background-color: #3498db;
}

.download-btn:hover {
  background-color: #2980b9;
}

.edit-btn {
  background-color: #f39c12;
}

.edit-btn:hover {
  background-color: #e67e22;
}

.delete-file-btn {
  background-color: #e74c3c;
}

.delete-file-btn:hover {
  background-color: #c0392b;
}

.add-btn {
  background-color: #2ecc71;
}

.add-btn:hover {
  background-color: #27ae60;
}

.save-btn {
  background-color: #2ecc71;
}

.save-btn:hover {
  background-color: #27ae60;
}

.cancel-btn {
  background-color: #95a5a6;
}

.cancel-btn:hover {
  background-color: #7f8c8d;
}

.delete-btn {
  background-color: #e74c3c;
  padding: 5px 10px;
}

.delete-btn:hover {
  background-color: #c0392b;
}
</style>
