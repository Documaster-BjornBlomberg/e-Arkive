<template>
  <MainLayout>
    <div class="document-list">
      <div class="list-header">
        <div class="view-controls">
          <Button
            class="view-toggle-btn"
            :title="viewMode === 'sidepanel' ? 'Byt till expanderbar vy' : 'Byt till sidopanel vy'"
            @click="toggleViewMode"
          >
            {{ viewMode === 'sidepanel' ? '⬍ Expanderbar vy' : '⬌ Sidopanel vy' }}
          </Button>
          <Button
            class="refresh-btn"
            title="Uppdatera listan"
            @click="fetchFiles"
            :disabled="loading"
          >
            {{ loading ? 'Laddar...' : '🔄 Uppdatera' }}
          </Button>
        </div>
      </div>

      <!-- Using the FileList organism to display files -->
      <div class="main-content" :class="{ 'with-sidepanel': viewMode === 'sidepanel' && selectedFile }">
        <FileList 
          :files="files"
          :isLoading="loading"
          :viewMode="viewMode"
          :selectedFileId="selectedFileId"
          :expandedFileId="expandedFileId"
          @select-file="selectFile"
          @toggle-expand="toggleExpand"
        />
      </div>

      <!-- Show error message if there's an error fetching files -->
      <div v-if="error" class="error-message">
        <p>Ett fel uppstod: {{ error }}</p>
        <Button @click="fetchFiles">Försök igen</Button>
      </div>

      <!-- File detail view using the side panel -->
      <FileDetailSidepanel
        v-if="viewMode === 'sidepanel'"
        :file="selectedFile"
        :isOpen="!!selectedFile"
        @close="closeDetail"
      />
    </div>
  </MainLayout>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import MainLayout from '../templates/MainLayout.vue';
import Button from '../atoms/Button.vue';
import FileList from '../organisms/FileList.vue';
import FileDetailSidepanel from '../organisms/FileDetailSidepanel.vue';
import { useGraphQL } from '../../composables/useGraphQL';

// Set up GraphQL client
const { loading, error, getFiles, getFileById } = useGraphQL();

const viewMode = ref('sidepanel');
const files = ref([]);
const selectedFileId = ref(null);
const expandedFileId = ref(null);
const selectedFile = ref(null);

const toggleViewMode = () => {
  viewMode.value = viewMode.value === 'sidepanel' ? 'expandable' : 'sidepanel';
};

const fetchFiles = async () => {
  try {
    const fileData = await getFiles();
    files.value = fileData;
    
    // Log successful file fetch
    console.log(`Successfully loaded ${fileData.length} files`);
  } catch (err) {
    console.error('Error fetching files:', err);
  }
};

const selectFile = async (file) => {
  selectedFileId.value = file.id;
  
  try {
    // Get detailed file info if needed
    const detailedFile = await getFileById(file.id);
    selectedFile.value = detailedFile || file;
    console.log('File selected:', selectedFile.value);
  } catch (err) {
    console.error('Error fetching file details:', err);
    selectedFile.value = file; // Fall back to the basic file info
  }
};

const closeDetail = () => {
  selectedFileId.value = null;
  selectedFile.value = null;
};

const toggleExpand = (fileId) => {
  if (expandedFileId.value === fileId) {
    expandedFileId.value = null;
  } else {
    expandedFileId.value = fileId;
    // Load file details for expanded view if needed
    const file = files.value.find(f => f.id === fileId);
    if (file) selectFile(file);
  }
};

// Watch for changes in the file list
watch(files, (newFiles) => {
  // If we have files and one was selected before, try to reselect it
  if (newFiles.length > 0 && selectedFileId.value) {
    const stillExists = newFiles.some(file => file.id === selectedFileId.value);
    if (!stillExists) {
      selectedFileId.value = null;
      selectedFile.value = null;
    }
  }
});

onMounted(() => {
  console.log('ListPage mounted, fetching files...');
  fetchFiles();
});
</script>

<style scoped>
/* Styles from DocumentList.vue */
.document-list {
  padding: 20px;
  position: relative;
  display: flex;
  flex-direction: column;
  height: calc(100vh - 100px);
}

.list-header {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--border-color);
}

.main-content {
  flex: 1;
  transition: width 0.3s ease;
  width: 100%;
}

.main-content.with-sidepanel {
  width: calc(100% - 450px);
}

.view-controls {
  display: flex;
  gap: 10px;
}

.error-message {
  margin-top: 20px;
  padding: 15px;
  background-color: rgba(231, 76, 60, 0.1);
  border-left: 4px solid #e74c3c;
  border-radius: 4px;
  color: #c0392b;
}
</style>
