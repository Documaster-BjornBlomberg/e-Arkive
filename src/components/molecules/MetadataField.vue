<script setup>
import MetadataInputField from '../atoms/MetadataInputField.vue'
import Button from '../atoms/Button.vue'

const props = defineProps({
  metadataItem: {
    type: Object,
    required: true
  },
  disabled: {
    type: Boolean,
    default: false
  },
  removable: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:key', 'update:value', 'remove'])

const updateKey = (value) => {
  emit('update:key', value)
}

const updateValue = (value) => {
  emit('update:value', value)
}
</script>

<template>
  <div class="metadata-item">
    <div class="metadata-inputs">
      <div class="input-field">
        <MetadataInputField
          label="Nyckel"
          :model-value="metadataItem.key"
          placeholder="t.ex. Författare, Datum, Kategori"
          :disabled="disabled"
          @update:model-value="updateKey"
        />
      </div>
      
      <div class="input-field">
        <MetadataInputField
          label="Värde"
          :model-value="metadataItem.value"
          placeholder="t.ex. John Doe, 2025-04-02"
          :disabled="disabled"
          @update:model-value="updateValue"
        />
      </div>
    </div>
    
    <div class="button-container">
      <Button 
        v-if="removable"
        @click="$emit('remove')" 
        variant="danger"
        :disabled="disabled"
        title="Ta bort metadata"
      >
        ✕
      </Button>
    </div>
  </div>
</template>

<style scoped>
.metadata-item {
  display: flex;
  align-items: flex-start; /* Changed from center to prevent vertical alignment issues */
  gap: 10px;
  margin-bottom: 15px;
}

.metadata-inputs {
  display: flex;
  flex: 1;
  gap: 15px;
  flex-wrap: wrap; /* Allow metadata fields to wrap on small screens */
}

.input-field {
  flex: 1;
  min-width: 150px; /* Ensure fields have a minimum width before wrapping */
}

.button-container {
  margin-top: 29px; /* Align button with inputs by adding top margin to match input label height */
}

@media (max-width: 600px) {
  .metadata-inputs {
    flex-direction: column;
  }
  
  .button-container {
    margin-top: 0;
    align-self: center;
  }
}
</style>