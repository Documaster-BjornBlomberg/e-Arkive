<template>
  <div class="file-action-buttons">
    <Button
      variant="info"
      title="Ladda ner fil"
      @click="$emit('download')"
    >
      Ladda ner
    </Button>
    
    <Button
      v-if="!isEditing"
      variant="primary"
      title="Redigera metadata"
      @click="$emit('edit')"
    >
      Redigera metadata
    </Button>
    
    <Button
      v-else
      variant="success"
      title="Spara metadata"
      @click="$emit('save')"
    >
      Spara
    </Button>
    
    <Button
      v-if="isEditing"
      variant="default"
      title="Avbryt redigering"
      @click="$emit('cancel')"
    >
      Avbryt
    </Button>
    
    <Button
      variant="danger"
      title="Ta bort fil"
      @click="$emit('delete')"
    >
      Ta bort
    </Button>
  </div>
</template>

<script setup>
import Button from '../atoms/Button.vue';

defineProps({
  isEditing: {
    type: Boolean,
    default: false
  }
});

defineEmits(['download', 'edit', 'save', 'cancel', 'delete']);
</script>

<style scoped>
.file-action-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 20px;
}
</style>