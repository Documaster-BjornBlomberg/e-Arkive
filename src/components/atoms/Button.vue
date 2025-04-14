<template>
  <button 
    :class="['button', variant, { 'full-width': fullWidth }]"
    :title="title"
    :type="type"
    :disabled="disabled || loading"
    @click="$emit('click', $event)">
    <span v-if="loading" class="loading-icon">
      <span class="material-icons">sync</span>
    </span>
    <span v-else-if="icon" class="material-icons button-icon">{{ icon }}</span>
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
  icon: {
    type: String,
    default: ''
  },
  title: {
    type: String,
    default: ''
  },
  fullWidth: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  type: {
    type: String,
    default: 'button'
  },
  primary: {
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

.button:hover:not(:disabled) {
  transform: translateY(-1px);
}

.button:active:not(:disabled) {
  transform: translateY(0);
}

.default {
  background-color: #95a5a6;
}

.default:hover:not(:disabled) {
  background-color: #7f8c8d;
}

.primary {
  background-color: var(--button-bg);
}

.primary:hover:not(:disabled) {
  background-color: var(--button-hover-bg);
}

.success {
  background-color: #2ecc71;
}

.success:hover:not(:disabled) {
  background-color: #27ae60;
}

.warning {
  background-color: #f39c12;
}

.warning:hover:not(:disabled) {
  background-color: #e67e22;
}

.danger {
  background-color: #e74c3c;
}

.danger:hover:not(:disabled) {
  background-color: #c0392b;
}

.info {
  background-color: #3498db;
}

.info:hover:not(:disabled) {
  background-color: #2980b9;
}

.full-width {
  width: 100%;
}

.button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.button-icon {
  font-size: 18px;
}

.loading-icon {
  display: inline-flex;
  align-items: center;
  animation: spin 1s infinite linear;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
