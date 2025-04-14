<template>
  <div 
    class="folder-node"
    :class="{ 'expanded': isExpanded, 'selected': isSelected }"
    @click.stop="$emit('select')"
    @contextmenu.prevent="$emit('contextmenu', $event)">
    <div class="folder-node-content">
      <span 
        class="expand-icon" 
        @click.stop="$emit('toggle-expand', node)"
        v-if="hasChildren || isExpanded">
        {{ isExpanded ? '▼' : '▶' }}
      </span>
      <span class="expand-icon placeholder" v-else></span>
      
      <span class="material-icons folder-icon">
        {{ isExpanded ? 'folder_open' : 'folder' }}
      </span>
      
      <span class="folder-name">{{ node.name }}</span>
      
      <span class="file-count" v-if="node.files && node.files.length > 0">
        {{ node.files.length }}
      </span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  node: {
    type: Object,
    required: true
  },
  isExpanded: {
    type: Boolean,
    default: false
  },
  isSelected: {
    type: Boolean,
    default: false
  }
});

defineEmits(['toggle-expand', 'select', 'contextmenu']);

// Check if the node has children
const hasChildren = computed(() => {
  return props.node.children && props.node.children.length > 0;
});
</script>

<style scoped>
.folder-node {
  cursor: pointer;
  user-select: none;
  border-radius: 4px;
  margin: 2px 0;
  transition: background-color 0.2s;
  color: var(--text-color);
}

.folder-node:hover {
  background-color: rgba(128, 128, 128, 0.1);
}

.folder-node.selected {
  background-color: var(--primary-color-light-transparent);
  font-weight: 500;
}

.folder-node-content {
  display: flex;
  align-items: center;
  padding: 6px 8px;
  gap: 5px;
}

.folder-icon {
  color: var(--folder-color, #ffd54f);
  font-size: 18px;
}

.folder-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--text-color);
}

.expand-icon {
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  color: var(--text-color-secondary);
}

.expand-icon.placeholder {
  visibility: hidden;
}

.file-count {
  background-color: rgba(128, 128, 128, 0.2);
  color: var(--text-color-secondary);
  border-radius: 10px;
  padding: 1px 6px;
  font-size: 11px;
  margin-left: 5px;
}
</style>