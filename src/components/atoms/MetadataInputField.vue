<script setup>
import { computed } from 'vue';

const props = defineProps({
  label: {
    type: String,
    required: true
  },
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  },
  disabled: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:modelValue']);

const inputValue = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});
</script>

<template>
  <div class="metadata-input-field">
    <label class="label">{{ label }}</label>
    <input
      type="text"
      class="input"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      @input="$event => emit('update:modelValue', $event.target.value)"
    />
  </div>
</template>

<style scoped>
.metadata-input-field {
  margin-bottom: 12px;
}

.label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
  color: var(--text-color);
}

.input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  background-color: var(--background-color);
  color: var(--text-color);
}

.input:focus {
  outline: none;
  border-color: var(--button-bg);
  box-shadow: 0 0 0 2px rgba(76, 175, 80, 0.2);
}

.input:disabled {
  background-color: rgba(128, 128, 128, 0.1);
  cursor: not-allowed;
  opacity: 0.7;
}

.input::placeholder {
  color: var(--text-color);
  opacity: 0.5;
}
</style>