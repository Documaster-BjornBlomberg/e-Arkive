<template>
  <div class="file-input-container">
    <label :for="id" class="file-label">{{ label }}</label>
    <div class="custom-file-input">
      <Button 
        variant="info" 
        :disabled="disabled"
        @click="triggerFileInput"
        class="file-select-button"
      >
        <span>Välj fil</span>
      </Button>
      <span class="selected-file-name">{{ selectedFileName || 'Ingen fil vald' }}</span>
      <input
        :id="id"
        ref="fileInput"
        type="file"
        class="hidden-file-input"
        @change="handleChange"
        :disabled="disabled"
      />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import Button from './Button.vue';

const props = defineProps({
  id: {
    type: String,
    required: true,
  },
  label: {
    type: String,
    required: true,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['change']);
const fileInput = ref(null);
const selectedFileName = ref('');

const triggerFileInput = () => {
  if (!props.disabled && fileInput.value) {
    fileInput.value.click();
  }
};

const handleChange = (event) => {
  if (event.target.files && event.target.files.length > 0) {
    selectedFileName.value = event.target.files[0].name;
  } else {
    selectedFileName.value = '';
  }
  emit('change', event);
};
</script>

<style scoped>
.file-input-container {
  margin-bottom: 20px;
}

.file-label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
  color: var(--text-color);
}

.custom-file-input {
  display: flex;
  align-items: center;
  gap: 12px;
}

.hidden-file-input {
  display: none;
}

.selected-file-name {
  color: var(--text-color);
  font-size: 0.9rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 300px;
}
</style>
