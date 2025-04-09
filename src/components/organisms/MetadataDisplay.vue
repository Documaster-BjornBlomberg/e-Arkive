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
        :metadata="metadata" 
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

const searchInput = ref(props.metadataSearch);
const activeTabValue = ref(props.activeTab);

// Update local values when props change
watch(() => props.activeTab, (newVal) => {
  activeTabValue.value = newVal;
});

watch(() => props.metadataSearch, (newVal) => {
  searchInput.value = newVal;
});

// Lista av metadata-fliktabbar
const metadataTabs = [
  { id: 'list', label: 'Lista' },
  { id: 'cards', label: 'Kort' },
  { id: 'categories', label: 'Kategorier' }
];

// Filtrerad metadata baserad på sökning
const filteredMetadata = computed(() => {
  if (!props.metadata) return [];
  
  if (!searchInput.value.trim()) return props.metadata;
  
  const search = searchInput.value.toLowerCase();
  return props.metadata.filter(meta => 
    meta.key.toLowerCase().includes(search) || 
    meta.value.toLowerCase().includes(search)
  );
});
</script>

<style scoped>
.metadata-section {
  margin-bottom: 20px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.section-header h3 {
  margin: 0;
}

.metadata-search {
  padding: 5px 10px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  width: 150px;
}

.metadata-tabs {
  display: flex;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 20px;
}
</style>