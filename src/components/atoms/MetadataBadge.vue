<template>
  <span class="metadata-badge">
    {{ count }} metadatafält
  </span>
</template>

<script setup>
defineProps({
  count: {
    type: Number,
    default: 0
  }
});
</script>

<style scoped>
.metadata-badge {
  font-size: 0.85rem;
  color: #666;
  background-color: rgba(0, 0, 0, 0.05);
  border-radius: 10px;
  padding: 2px 8px;
}
</style>