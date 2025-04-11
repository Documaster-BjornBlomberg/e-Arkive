<template>
  <div class="folder-node">
    <div 
      class="folder-item" 
      @click="toggleExpand"
      :class="{ 'selected': isSelected }"
    >
      <span class="folder-toggle">
        <span v-if="hasChildren || hasLoadedChildren" class="toggle-icon">
          {{ isExpanded ? '▼' : '▶' }}
        </span>
        <span v-else class="toggle-spacer"></span>
      </span>
      <span class="folder-icon">📁</span>
      <span class="folder-name">{{ node.name }}</span>
    </div>
    
    <div v-if="isExpanded" class="folder-children">
      <div v-if="isLoading" class="loading">Laddar...</div>
      <template v-else>
        <FolderNode 
          v-for="child in children" 
          :key="child.id" 
          :node="child"
          :selected-node-id="selectedNodeId"
          @select-node="$emit('select-node', $event)"
        />
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, inject } from 'vue';
import { useNodeHandling } from '../../composables/useNodeHandling';

const props = defineProps({
  node: {
    type: Object,
    required: true
  },
  selectedNodeId: {
    type: String,
    default: null
  }
});

const emit = defineEmits(['select-node']);

// Get node handling functions from the parent or inject them
const {
  toggleNodeExpanded, 
  isNodeExpanded, 
  loadedChildrenMap,
  fetchNodeChildren
} = useNodeHandling();

const isLoading = ref(false);
const children = ref([]);

// Computed properties
const isExpanded = computed(() => isNodeExpanded(props.node.id));
const isSelected = computed(() => props.selectedNodeId === props.node.id);
const hasChildren = computed(() => props.node.children && props.node.children.length > 0);
const hasLoadedChildren = computed(() => {
  return loadedChildrenMap.has(props.node.id) && loadedChildrenMap.get(props.node.id).length > 0;
});

// Methods
const toggleExpand = async () => {
  await toggleNodeExpanded(props.node.id);
  
  if (isExpanded.value && !children.value.length) {
    await loadChildren();
  }
  
  emit('select-node', props.node.id);
};

const loadChildren = async () => {
  if (isLoading.value) return;
  
  isLoading.value = true;
  try {
    const nodeChildren = await fetchNodeChildren(props.node.id);
    children.value = nodeChildren || [];
  } catch (error) {
    console.error(`Failed to load children for node ${props.node.id}:`, error);
  } finally {
    isLoading.value = false;
  }
};

// Initialize component
onMounted(() => {
  if (loadedChildrenMap.has(props.node.id)) {
    children.value = loadedChildrenMap.get(props.node.id);
  }
});
</script>

<style scoped>
.folder-node {
  font-size: 0.9rem;
}

.folder-item {
  display: flex;
  align-items: center;
  padding: 4px 0;
  cursor: pointer;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  border-radius: 4px;
  transition: background-color 0.2s ease;
}

.folder-item:hover {
  background-color: rgba(128, 128, 128, 0.1);
}

.folder-item.selected {
  background-color: rgba(51, 122, 183, 0.2);
}

.folder-toggle {
  width: 20px;
  display: inline-block;
  text-align: center;
}

.toggle-icon {
  font-size: 0.7rem;
}

.toggle-spacer {
  display: inline-block;
  width: 10px;
}

.folder-icon {
  margin-right: 5px;
}

.folder-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
}

.folder-children {
  padding-left: 20px;
}

.loading {
  padding: 4px 0 4px 25px;
  font-style: italic;
  font-size: 0.8rem;
  color: #777;
}
</style>