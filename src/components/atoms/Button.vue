<template>
  <button 
    :class="['button', variant, { 'full-width': fullWidth }]"
    :title="title"
    @click="$emit('click', $event)">
    <slot></slot>
  </button>
</template>

<script setup>
defineProps({
  variant: {
    type: String,
    default: 'default',
    validator: (value) => ['default', 'primary', 'success', 'warning', 'danger', 'info'].includes(value)
  },
  title: {
    type: String,
    default: ''
  },
  fullWidth: {
    type: Boolean,
    default: false
  }
});

defineEmits(['click']);
</script>

<style scoped>
.button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s, transform 0.1s;
  font-size: 0.9rem;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  color: var(--button-text);
}

.button:hover {
  transform: translateY(-1px);
}

.button:active {
  transform: translateY(0);
}

.default {
  background-color: #95a5a6;
}

.default:hover {
  background-color: #7f8c8d;
}

.primary {
  background-color: var(--button-bg);
}

.primary:hover {
  background-color: var(--button-hover-bg);
}

.success {
  background-color: #2ecc71;
}

.success:hover {
  background-color: #27ae60;
}

.warning {
  background-color: #f39c12;
}

.warning:hover {
  background-color: #e67e22;
}

.danger {
  background-color: #e74c3c;
}

.danger:hover {
  background-color: #c0392b;
}

.info {
  background-color: #3498db;
}

.info:hover {
  background-color: #2980b9;
}

.full-width {
  width: 100%;
}

/* Add disabled button styling */
.button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}
</style>
