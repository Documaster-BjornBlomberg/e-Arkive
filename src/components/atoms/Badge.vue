<template>
  <span class="badge" :class="variant">
    <slot></slot>
  </span>
</template>

<script setup>
defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'secondary', 'success', 'danger', 'warning', 'info'].includes(value)
  }
})
</script>

<style scoped>
.badge {
  display: inline-block;
  padding: 2px 6px;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: bold;
  color: white;
}

.primary {
  background-color: var(--button-bg);
}

.secondary {
  background-color: #95a5a6;
}

.success {
  background-color: #2ecc71;
}

.danger {
  background-color: #e74c3c;
}

.warning {
  background-color: #f39c12;
}

.info {
  background-color: #3498db;
}
</style>