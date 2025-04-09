<template>
  <div class="metadata-section">
    <h2>Dokumentets metadata</h2>
    <p class="metadata-help">Lägg till nyckel-värde par för att beskriva ditt dokument</p>
    <div class="metadata-list">
      <MetadataField
        v-for="(meta, index) in metadata"
        :key="index"
        :metadata-item="meta"
        :disabled="disabled"
        @update:key="updateMetadataKey(index, $event)"
        @update:value="updateMetadataValue(index, $event)"
        @remove="removeMetadataItem(index)"
      />
    </div>
    <div class="metadata-actions">
      <Button 
        @click="addMetadata" 
        variant="info"
        :disabled="disabled"
      >
        Lägg till metadata
      </Button>
    </div>
  </div>
</template>

<script setup>
import MetadataField from '../molecules/MetadataField.vue';
import Button from '../atoms/Button.vue';
import { defineProps, defineEmits } from 'vue';

const props = defineProps({
  metadata: {
    type: Array,
    default: () => []
  },
  disabled: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:metadata', 'add-metadata']);

const updateMetadataKey = (index, newKey) => {
  const updatedMetadata = [...props.metadata];
  updatedMetadata[index] = { ...updatedMetadata[index], key: newKey };
  emit('update:metadata', updatedMetadata);
};

const updateMetadataValue = (index, newValue) => {
  const updatedMetadata = [...props.metadata];
  updatedMetadata[index] = { ...updatedMetadata[index], value: newValue };
  emit('update:metadata', updatedMetadata);
};

const removeMetadataItem = (index) => {
  const updatedMetadata = [...props.metadata];
  updatedMetadata.splice(index, 1);
  emit('update:metadata', updatedMetadata);
};

const addMetadata = () => {
  emit('add-metadata');
};
</script>

<style scoped>
.metadata-section {
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 15px;
  background-color: rgba(128, 128, 128, 0.05);
  color: var(--text-color);
}

.metadata-section h2 {
  margin-top: 0;
  margin-bottom: 5px;
  color: var(--text-color);
}

.metadata-help {
  color: var(--text-color);
  opacity: 0.7;
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