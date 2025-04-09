<template>
  <div class="file-info" v-if="file">
    <p><strong>Selected File:</strong> {{ file.name }}</p>
    <p><strong>Size:</strong> {{ fileSize }} KB</p>
    <p><strong>Type:</strong> {{ file.type }}</p>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  file: {
    type: Object,
    default: null,
  },
});

const fileSize = computed(() => {
  return props.file ? Math.round(props.file.size / 1024) : 0;
});
</script>

<style scoped>
.file-info {
  background-color: rgba(128, 128, 128, 0.1);
  padding: 10px;
  border-radius: 4px;
  margin-bottom: 10px;
  border: 1px solid var(--border-color); /* Added border for consistency */
}

.file-info p {
  margin: 5px 0;
  font-size: 0.9em;
}
</style>
