<template>
  <div class="file-table-view">
    <div class="table-config">
      <Button 
        variant="default"
        @click="showColumnConfig = !showColumnConfig"
        class="config-button">
        <span class="material-icons">view_column</span>
        Visa kolumner
      </Button>
      
      <!-- Column configuration modal -->
      <div v-if="showColumnConfig" class="column-config-modal">
        <div class="config-header">
          <h3>Konfigurera kolumner</h3>
          <button class="close-btn" @click="showColumnConfig = false">×</button>
        </div>
        <div class="config-content">
          <div 
            v-for="col in availableColumns" 
            :key="col.id"
            class="column-option">
            <label>
              <input 
                type="checkbox"
                :checked="activeColumns.includes(col.id)"
                @change="toggleColumn(col.id)">
              {{ col.label }}
            </label>
          </div>
        </div>
        <div class="config-actions">
          <Button variant="default" @click="resetColumns">Återställ</Button>
          <Button variant="primary" @click="showColumnConfig = false">Stäng</Button>
        </div>
      </div>
    </div>
    
    <div class="table-container">
      <table class="file-table">
        <thead>
          <tr>
            <th v-if="showPreview"></th>
            <th 
              v-for="col in visibleColumns" 
              :key="col.id"
              :class="{ sortable: col.sortable }"
              @click="col.sortable && handleSort(col.id)">
              {{ col.label }}
              <span 
                v-if="col.sortable" 
                class="sort-icon material-icons">
                {{ getSortIcon(col.id) }}
              </span>
            </th>
            <th class="actions-column">Åtgärder</th>
          </tr>
        </thead>
        <tbody>
          <tr 
            v-for="file in sortedFiles" 
            :key="file.id"
            :class="{ selected: selectedFileId === file.id, expanded: expandedFileId === file.id }"
            @click="$emit('select', file.id)">
            
            <!-- Preview toggle column -->
            <td v-if="showPreview" class="preview-toggle">
              <button 
                class="expand-button"
                @click.stop="$emit('toggle-expand', file.id)">
                <span class="material-icons">
                  {{ expandedFileId === file.id ? 'expand_less' : 'expand_more' }}
                </span>
              </button>
            </td>
            
            <!-- Dynamic columns -->
            <td v-for="col in visibleColumns" :key="col.id">
              <template v-if="col.id === 'icon'">
                <FileIcon :fileName="file.name" :fileType="file.contentType" />
              </template>
              <template v-else-if="col.id === 'name'">
                <div class="name-cell">
                  {{ file.name }}
                  <div class="hover-preview">
                    <div class="preview-header">
                      <h4>{{ file.name }}</h4>
                      <div class="preview-actions">
                        <button 
                          class="preview-action"
                          title="Ladda ner"
                          @click.stop="$emit('download', file)">
                          <span class="material-icons">download</span>
                        </button>
                        <button 
                          class="preview-action"
                          title="Redigera metadata"
                          @click.stop="$emit('edit', file)">
                          <span class="material-icons">edit</span>
                        </button>
                      </div>
                    </div>
                    <div class="preview-metadata">
                      <div v-if="file.metadata?.length" class="metadata-preview">
                        <div 
                          v-for="(meta, index) in file.metadata.slice(0, 3)" 
                          :key="index"
                          class="metadata-item">
                          <span class="metadata-key">{{ meta.key }}:</span>
                          <span class="metadata-value">{{ meta.value }}</span>
                        </div>
                        <div v-if="file.metadata.length > 3" class="metadata-more">
                          +{{ file.metadata.length - 3 }} fler fält
                        </div>
                      </div>
                      <div v-else class="no-metadata">
                        Ingen metadata tillgänglig
                      </div>
                    </div>
                  </div>
                </div>
              </template>
              <template v-else-if="col.id === 'size'">
                {{ formatFileSize(file.size) }}
              </template>
              <template v-else-if="col.id === 'type'">
                {{ file.contentType }}
              </template>
              <template v-else-if="col.id === 'created'">
                {{ new Date(file.createdAt).toLocaleString() }}
              </template>
              <template v-else-if="col.id === 'metadataCount'">
                {{ file.metadata?.length || 0 }}
              </template>
              <template v-else-if="col.isMetadataField">
                {{ getMetadataValue(file, col.metadataKey) }}
              </template>
            </td>
            
            <!-- Actions column -->
            <td class="actions-column">
              <div class="action-buttons">
                <button 
                  class="action-btn"
                  title="Ladda ner"
                  @click.stop="$emit('download', file)">
                  <span class="material-icons">download</span>
                </button>
                <button 
                  class="action-btn"
                  title="Redigera metadata"
                  @click.stop="$emit('edit', file)">
                  <span class="material-icons">edit</span>
                </button>
              </div>
            </td>
          </tr>
          
          <!-- Expanded preview row -->
          <tr v-if="expandedFileId" class="preview-row">
            <td :colspan="totalColumns" class="preview-cell">
              <FilePreview :file="getExpandedFile()" />
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import Button from '../atoms/Button.vue';
import FileIcon from '../atoms/FileIcon.vue';
import FilePreview from '../organisms/FilePreview.vue';

const props = defineProps({
  files: {
    type: Array,
    required: true
  },
  selectedFileId: {
    type: String,
    default: null
  },
  expandedFileId: {
    type: String,
    default: null
  }
});

defineEmits(['select', 'toggle-expand', 'download', 'edit']);

// Column configuration state
const showColumnConfig = ref(false);
const availableColumns = [
  { id: 'icon', label: 'Ikon', sortable: false },
  { id: 'name', label: 'Namn', sortable: true },
  { id: 'size', label: 'Storlek', sortable: true },
  { id: 'type', label: 'Typ', sortable: true },
  { id: 'created', label: 'Skapad', sortable: true },
  { id: 'metadataCount', label: 'Metadata', sortable: true }
];

// Load saved column configuration or use defaults
const savedColumns = localStorage.getItem('tableColumns');
const activeColumns = ref(
  savedColumns ? JSON.parse(savedColumns) : ['icon', 'name', 'size', 'type', 'created']
);

// Track sort state
const sortField = ref('name');
const sortDirection = ref('asc');

// Show preview column if preview feature is enabled
const showPreview = ref(true);

// Computed columns that are currently visible
const visibleColumns = computed(() => {
  return availableColumns.filter(col => activeColumns.value.includes(col.id));
});

// Total number of columns (including actions and preview)
const totalColumns = computed(() => {
  return visibleColumns.value.length + (showPreview.value ? 2 : 1);
});

// Column management functions
const toggleColumn = (columnId) => {
  const index = activeColumns.value.indexOf(columnId);
  if (index === -1) {
    activeColumns.value.push(columnId);
  } else {
    activeColumns.value.splice(index, 1);
  }
  saveColumnConfig();
};

const resetColumns = () => {
  activeColumns.value = ['icon', 'name', 'size', 'type', 'created'];
  saveColumnConfig();
};

const saveColumnConfig = () => {
  localStorage.setItem('tableColumns', JSON.stringify(activeColumns.value));
};

// Sorting functions
const handleSort = (field) => {
  if (sortField.value === field) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc';
  } else {
    sortField.value = field;
    sortDirection.value = 'asc';
  }
};

const getSortIcon = (field) => {
  if (sortField.value !== field) return 'sort';
  return sortDirection.value === 'asc' ? 'arrow_upward' : 'arrow_downward';
};

// Sort files based on current sort field and direction
const sortedFiles = computed(() => {
  const sorted = [...props.files].sort((a, b) => {
    let valueA, valueB;
    
    switch (sortField.value) {
      case 'name':
        valueA = a.name.toLowerCase();
        valueB = b.name.toLowerCase();
        break;
      case 'size':
        valueA = a.size;
        valueB = b.size;
        break;
      case 'type':
        valueA = a.contentType.toLowerCase();
        valueB = b.contentType.toLowerCase();
        break;
      case 'created':
        valueA = new Date(a.createdAt).getTime();
        valueB = new Date(b.createdAt).getTime();
        break;
      case 'metadataCount':
        valueA = a.metadata?.length || 0;
        valueB = b.metadata?.length || 0;
        break;
      default:
        return 0;
    }
    
    if (valueA < valueB) return sortDirection.value === 'asc' ? -1 : 1;
    if (valueA > valueB) return sortDirection.value === 'asc' ? 1 : -1;
    return 0;
  });
  
  return sorted;
});

// Helper functions
const formatFileSize = (bytes) => {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
};

const getMetadataValue = (file, key) => {
  return file.metadata?.find(m => m.key === key)?.value || '';
};

const getExpandedFile = () => {
  return props.files.find(f => f.id === props.expandedFileId);
};
</script>

<style scoped>
.file-table-view {
  width: 100%;
  background-color: var(--surface-color);
}

.table-config {
  padding: 12px;
  border-bottom: 1px solid var(--border-color);
  position: relative;
}

.config-button {
  display: flex;
  align-items: center;
  gap: 6px;
}

.column-config-modal {
  position: absolute;
  top: 100%;
  left: 12px;
  background-color: var(--surface-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 100;
  min-width: 200px;
}

.config-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border-bottom: 1px solid var(--border-color);
}

.config-header h3 {
  margin: 0;
  font-size: 1rem;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 4px;
  color: var(--text-color);
}

.config-content {
  padding: 12px;
  max-height: 300px;
  overflow-y: auto;
}

.column-option {
  padding: 6px 0;
}

.column-option label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.config-actions {
  padding: 12px;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.table-container {
  overflow-x: auto;
  overflow-y: hidden;
}

.file-table {
  width: 100%;
  border-collapse: collapse;
}

.file-table th {
  background-color: var(--header-bg);
  font-weight: 500;
  text-align: left;
  padding: 12px;
  border-bottom: 1px solid var(--border-color);
  white-space: nowrap;
}

.file-table th.sortable {
  cursor: pointer;
  user-select: none;
}

.file-table th.sortable:hover {
  background-color: var(--hover-color);
}

.sort-icon {
  font-size: 18px;
  vertical-align: middle;
  margin-left: 4px;
}

.file-table td {
  padding: 12px;
  border-bottom: 1px solid var(--border-color);
  vertical-align: middle;
}

.file-table tr:hover {
  background-color: var(--hover-color);
}

.file-table tr.selected {
  background-color: var(--selection-color);
}

.preview-toggle {
  width: 40px;
  text-align: center;
}

.expand-button {
  background: none;
  border: none;
  padding: 4px;
  cursor: pointer;
  color: var(--text-color-secondary);
  border-radius: 4px;
}

.expand-button:hover {
  background-color: var(--hover-color);
  color: var(--text-color);
}

.actions-column {
  width: 100px;
  text-align: right;
  white-space: nowrap;
}

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 4px;
}

.action-btn {
  background: none;
  border: none;
  padding: 4px;
  cursor: pointer;
  color: var(--text-color-secondary);
  border-radius: 4px;
}

.action-btn:hover {
  background-color: var(--hover-color);
  color: var(--text-color);
}

.preview-row {
  background-color: var(--surface-color) !important;
}

.preview-cell {
  padding: 0 !important;
}

.name-cell {
  position: relative;
}

.hover-preview {
  display: none;
  position: absolute;
  left: 0;
  top: 100%;
  background-color: var(--surface-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 15px;
  min-width: 300px;
  max-width: 400px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 10;
}

.name-cell:hover .hover-preview {
  display: block;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.preview-header h4 {
  margin: 0;
  font-size: 1rem;
  color: var(--text-color);
  word-break: break-word;
  flex: 1;
  padding-right: 10px;
}

.preview-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.preview-action {
  background: none;
  border: none;
  padding: 4px;
  cursor: pointer;
  color: var(--text-color-secondary);
  border-radius: 4px;
}

.preview-action:hover {
  background-color: var(--hover-color);
  color: var(--text-color);
}

.metadata-preview {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.metadata-item {
  font-size: 0.8rem;
  display: flex;
  gap: 6px;
}

.metadata-key {
  color: var(--text-color-secondary);
  flex-shrink: 0;
}

.metadata-value {
  color: var(--text-color);
  word-break: break-word;
}

.metadata-more {
  font-size: 0.8rem;
  color: var(--text-color-secondary);
  font-style: italic;
  text-align: center;
  padding-top: 4px;
}

.no-metadata {
  font-size: 0.8rem;
  color: var(--text-color-secondary);
  font-style: italic;
  text-align: center;
}
</style>