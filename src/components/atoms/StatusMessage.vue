<template>
  <div class="upload-status" v-if="message">
    <p :class="{ 'error': isError }">{{ message }}</p>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  message: {
    type: String,
    default: '',
  },
});

const isError = computed(() => {
  return props.message.toLowerCase().startsWith('error');
});
</script>

<style scoped>
.upload-status {
  padding: 10px;
  border-radius: 4px;
  background-color: rgba(46, 204, 113, 0.1); /* Success background */
  border-left: 4px solid #2ecc71; /* Success border */
  margin-top: 15px; /* Added margin */
}

.upload-status p {
  margin: 0;
}

.upload-status .error {
  background-color: rgba(231, 76, 60, 0.1); /* Error background */
  border-left-color: #e74c3c; /* Error border */
  color: #e74c3c; /* Error text color */
}
</style>
