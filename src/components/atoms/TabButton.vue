<template>
  <button 
    @click="$emit('click')"
    :class="{ 'active': isActive }"
    class="tab-btn"
  >
    {{ label }}
  </button>
</template>

<script setup>
defineProps({
  label: {
    type: String,
    required: true
  },
  isActive: {
    type: Boolean,
    default: false
  }
});

defineEmits(['click']);
</script>

<style scoped>
.tab-btn {
  background: none;
  border: none;
  padding: 8px 15px;
  cursor: pointer;
  border-radius: 4px 4px 0 0;
  transition: all 0.2s;
  color: var(--text-color);
  opacity: 0.7;
}

.tab-btn.active {
  border-bottom: 2px solid var(--button-bg);
  opacity: 1;
  font-weight: bold;
}

.tab-btn:hover {
  background-color: rgba(128, 128, 128, 0.1);
}
</style>