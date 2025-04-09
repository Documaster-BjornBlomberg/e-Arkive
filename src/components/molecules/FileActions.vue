<template>
  <div class="action-buttons">
    <Button 
      @click="$emit('download')" 
      class="download-btn" 
      title="Ladda ner"
    >
      📥 Ladda ner
    </Button>
    <Button 
      v-if="showEdit" 
      @click="$emit('edit')" 
      class="edit-btn" 
      title="Redigera metadata"
    >
      ✏️ Redigera metadata
    </Button>
    <Button 
      @click="$emit('delete')" 
      class="delete-file-btn" 
      title="Ta bort fil"
    >
      🗑️ Ta bort
    </Button>
  </div>
</template>

<script setup>
import Button from '../atoms/Button.vue';

defineProps({
  showEdit: {
    type: Boolean,
    default: true
  }
});

defineEmits(['download', 'edit', 'delete']);
</script>

<style scoped>
.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.download-btn {
  background-color: #3498db;
}

.download-btn:hover {
  background-color: #2980b9;
}

.edit-btn {
  background-color: #f39c12;
}

.edit-btn:hover {
  background-color: #e67e22;
}

.delete-file-btn {
  background-color: #e74c3c;
}

.delete-file-btn:hover {
  background-color: #c0392b;
}
</style>