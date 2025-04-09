<template>
  <div class="search-input">
    <Icon name="search" class="search-icon" />
    <input
      type="text"
      class="input-field"
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
      :placeholder="placeholder"
      @keyup.enter="$emit('search')"
    />
    <button 
      v-if="modelValue" 
      class="clear-button" 
      @click="clearSearch" 
      title="Rensa sökning"
    >
      ×
    </button>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue';
import Icon from '../atoms/Icon.vue';

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: 'Sök...'
  }
});

const emit = defineEmits(['update:modelValue', 'search']);

const clearSearch = () => {
  emit('update:modelValue', '');
  emit('search');
};
</script>

<style scoped>
.search-input {
  position: relative;
  width: 100%;
}

.search-icon {
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
  color: #666;
}

.input-field {
  width: 100%;
  padding: 10px 30px 10px 35px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
}

.input-field:focus {
  border-color: #4CAF50;
}

.clear-button {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  font-size: 18px;
  color: #999;
  cursor: pointer;
  padding: 0;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.2s;
}

.clear-button:hover {
  color: #ff5252;
}
</style>