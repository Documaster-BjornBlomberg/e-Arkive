<template>
  <div class="metadata-item">
    <div class="metadata-inputs">
      <InputField
        label="Key:"
        :modelValue="item.key"
        @update:modelValue="$emit('update:item', { ...item, key: $event })"
        placeholder="e.g. Author, Date, Category"
        :disabled="disabled"
      />
      <InputField
        label="Value:"
        :modelValue="item.value"
        @update:modelValue="$emit('update:item', { ...item, value: $event })"
        placeholder="e.g. John Doe, 2025-04-02"
        :disabled="disabled"
      />
    </div>
    <Button
      variant="remove"
      @click="$emit('remove')"
      :disabled="disableRemove || disabled"
      title="Remove this metadata"
    >
      ✕
    </Button>
  </div>
</template>

<script setup>
import InputField from '../atoms/InputField.vue';
import Button from '../atoms/Button.vue';

defineProps({
  item: {
    type: Object,
    required: true,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  disableRemove: {
    type: Boolean,
    default: false,
  },
});

defineEmits(['update:item', 'remove']);
</script>

<style scoped>
.metadata-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.metadata-inputs {
  display: flex;
  flex: 1;
  gap: 10px;
}
</style>
