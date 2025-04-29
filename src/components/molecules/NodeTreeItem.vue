<template>
  <div class="node-tree-item">
    <div 
      class="node-row" 
      :class="{ 'selected': selectedNodeId === node.id }" 
      @click="$emit('select', node.id)"
    >
      <div class="expander" @click.stop="$emit('toggle-expanded', node.id)" v-if="hasChildren">
        <span class="material-icons">
          {{ isExpanded ? 'keyboard_arrow_down' : 'keyboard_arrow_right' }}
        </span>
      </div>
      <div class="expander spacer" v-else></div>
      
      <div class="node-icon">
        <span class="material-icons">folder</span>
      </div>
      
      <div class="node-name">{{ node.name }}</div>
    </div>
    
    <!-- Render children if expanded -->
    <div class="children" v-if="isExpanded && node.children && node.children.length > 0">
      <node-tree-item
        v-for="child in node.children"
        :key="child.id"
        :node="child"
        :expanded-nodes="expandedNodes"
        :selected-node-id="selectedNodeId"
        @select="$emit('select', $event)"
        @toggle-expanded="$emit('toggle-expanded', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { Node } from '../../types/schema';

const props = defineProps<{
  node: Node;
  expandedNodes: Set<string>;
  selectedNodeId?: string | null;
}>();

defineEmits<{
  (e: 'select', nodeId: string): void;
  (e: 'toggle-expanded', nodeId: string): void;
}>();

// Computed properties
const isExpanded = computed(() => {
  return props.expandedNodes.has(props.node.id);
});

const hasChildren = computed(() => {
  // Consider a node to have children if it has children array with items,
  // or if we haven't loaded its children yet (assuming folders might have children)
  return !!(props.node.children && props.node.children.length > 0);
});
</script>

<style scoped>
.node-tree-item {
  user-select: none;
}

.node-row {
  display: flex;
  align-items: center;
  padding: 6px 8px;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.node-row:hover {
  background-color: var(--hover-color);
}

.node-row.selected {
  background-color: var(--primary-color-light);
}

.expander {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 4px;
}

.expander.spacer {
  visibility: hidden;
}

.node-icon {
  color: var(--primary-color);
  margin-right: 8px;
}

.children {
  padding-left: 24px;
}

.material-icons {
  font-size: 20px;
}
</style>
