<template>
  <span class="badge" :class="badgeClass">
    {{ text }}
    <span v-if="closable" class="badge-close" @click.stop="$emit('close')">×</span>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  text: string;
  variant?: 'primary' | 'success' | 'warning' | 'error' | 'info' | 'default';
  closable?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  closable: false
});

defineEmits(['close']);

const badgeClass = computed(() => {
  return {
    'badge-primary': props.variant === 'primary',
    'badge-success': props.variant === 'success',
    'badge-warning': props.variant === 'warning',
    'badge-error': props.variant === 'error',
    'badge-info': props.variant === 'info',
    'badge-default': props.variant === 'default'
  };
});
</script>

<style scoped>
.badge {
  display: inline-flex;
  align-items: center;
  padding: 0.25em 0.75em;
  font-size: 0.75rem;
  font-weight: 500;
  border-radius: 9999px;
  white-space: nowrap;
}

.badge-primary {
  background-color: var(--primary-color);
  color: white;
}

.badge-success {
  background-color: var(--success-color);
  color: white;
}

.badge-warning {
  background-color: #f39c12;
  color: white;
}

.badge-error {
  background-color: var(--error-color);
  color: white;
}

.badge-info {
  background-color: #3498db;
  color: white;
}

.badge-default {
  background-color: var(--border-color);
  color: var(--text-color);
}

.badge-close {
  margin-left: 0.25em;
  font-size: 1.1em;
  line-height: 0.8;
  cursor: pointer;
}

.badge-close:hover {
  opacity: 0.7;
}
</style>
