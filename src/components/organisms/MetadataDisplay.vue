<template>
  <div class="metadata-section">
    <div class="section-header">
      <h3>Metadata ({{ filteredMetadata.length }})</h3>
      <input 
        type="text" 
        v-model="searchInput" 
        placeholder="Sök metadata..." 
        class="metadata-search"
        @input="$emit('update:metadataSearch', searchInput)"
      />
    </div>
    
    <div class="metadata-tabs">
      <TabButton 
        v-for="tab in metadataTabs" 
        :key="tab.id" 
        :label="tab.label"
        :isActive="activeTabValue === tab.id"
        @click="$emit('update:activeTab', tab.id)"
      />
    </div>
    
    <!-- List View -->
    <div v-if="activeTabValue === 'list'" class="metadata-list-view">
      <MetadataListView 
        :metadata="filteredMetadata" 
        @delete-metadata="$emit('delete-metadata', $event)"
      />
    </div>
    
    <!-- Card View -->
    <div v-else-if="activeTabValue === 'cards'" class="metadata-card-view">
      <MetadataCardView 
        :metadata="filteredMetadata" 
        @delete-metadata="$emit('delete-metadata', $event)"
      />
    </div>
    
    <!-- Category View -->
    <div v-else-if="activeTabValue === 'categories'" class="metadata-category-view">
      <MetadataCategoryView 
        :metadata="filteredMetadata" 
        @delete-metadata="$emit('delete-metadata', $event)"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import TabButton from '../atoms/TabButton.vue';
import MetadataListView from '../molecules/MetadataListView.vue';
import MetadataCardView from '../molecules/MetadataCardView.vue';
import MetadataCategoryView from '../molecules/MetadataCategoryView.vue';

const props = defineProps({
  metadata: {
    type: Array,
    default: () => []
  },
  activeTab: {
    type: String,
    default: 'list'
  },
  metadataSearch: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['update:activeTab', 'update:metadataSearch', 'delete-metadata']);

// Local state
const activeTabValue = ref(props.activeTab);
const searchInput = ref(props.metadataSearch);

// Available tabs
const metadataTabs = [
  { id: 'list', label: 'Lista' },
  { id: 'cards', label: 'Kort' },
  { id: 'categories', label: 'Kategorier' }
];

// Filter metadata based on search query
const filteredMetadata = computed(() => {
  if (!searchInput.value.trim()) return props.metadata || [];
  
  const query = searchInput.value.toLowerCase();
  return (props.metadata || []).filter(meta => 
    meta.key.toLowerCase().includes(query) || 
    meta.value.toLowerCase().includes(query)
  );
});

// Watch for prop changes
watch(() => props.activeTab, (newValue) => {
  activeTabValue.value = newValue;
});

watch(() => props.metadataSearch, (newValue) => {
  searchInput.value = newValue;
});
</script>

<style scoped>
.metadata-section {
  padding: 15px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background-color: var(--surface-color);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.section-header h3 {
  margin: 0;
}

.metadata-search {
  padding: 6px 10px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  width: 150px;
  font-size: 0.9rem;
  background-color: var(--input-background);
  color: var(--text-color);
}

.metadata-tabs {
  display: flex;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 15px;
}

.metadata-list-view, 
.metadata-card-view, 
.metadata-category-view {
  max-height: 300px;
  overflow-y: auto;
  padding-right: 5px;
}
</style>