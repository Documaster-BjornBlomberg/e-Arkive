<template>
  <div class="metadata-section">
    <h2>Document Metadata</h2>
    <p class="metadata-help">Add key-value pairs to describe your document</p>
    
    <div class="metadata-list">
      <MetadataItem
        v-for="(item, index) in modelValue"
        :key="index"
        :item="item"
        :disabled="disabled"
        :disable-remove="modelValue.length === 1"
        @update:item="updateItem(index, $event)"
        @remove="removeItem(index)"
      />
    </div>
    
    <div class="metadata-actions">
      <Button
        variant="add"
        @click="addItem"
        :disabled="disabled"
      >
        Add More Metadata
      </Button>
    </div>
  </div>
</template>

<script setup>
import MetadataItem from './MetadataItem.vue';
import Button from '../atoms/Button.vue';

const props = defineProps({
  modelValue: {
    type: Array,
    required: true,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['update:modelValue']);

const updateItem = (index, newItem) => {
  const newList = [...props.modelValue];
  newList[index] = newItem;
  emit('update:modelValue', newList);
};

const removeItem = (index) => {
  const newList = [...props.modelValue];
  newList.splice(index, 1);
  emit('update:modelValue', newList);
};

const addItem = () => {
  emit('update:modelValue', [...props.modelValue, { key: '', value: '' }]);
};
</script>

<style scoped>
.metadata-section {
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 15px;
  background-color: rgba(128, 128, 128, 0.05);
}

.metadata-section h2 {
  margin-top: 0;
  margin-bottom: 5px;
}

.metadata-help {
  color: #666;
  margin-bottom: 15px;
  font-size: 0.9em;
}

.metadata-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.metadata-actions {
  margin-top: 10px;
}
</style>
