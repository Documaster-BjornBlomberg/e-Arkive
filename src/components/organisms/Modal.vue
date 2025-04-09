<template>
  <div v-if="isOpen" class="modal-overlay" @click.self="closeModal">
    <div class="modal-container">
      <div class="modal-header">
        <h2>{{ title }}</h2>
        <Icon name="close" class="close-icon" @click="closeModal" />
      </div>
      <div class="modal-content">
        <slot></slot>
      </div>
      <div class="modal-footer">
        <Button v-if="showCancel" variant="secondary" @click="closeModal">Cancel</Button>
        <Button v-if="showConfirm" :variant="confirmVariant" @click="confirm">{{ confirmText }}</Button>
      </div>
    </div>
  </div>
</template>

<script>
import Icon from '../atoms/Icon.vue'
import Button from '../atoms/Button.vue'

export default {
  name: 'Modal',
  components: {
    Icon,
    Button
  },
  props: {
    isOpen: {
      type: Boolean,
      default: false
    },
    title: {
      type: String,
      default: 'Modal Title'
    },
    showCancel: {
      type: Boolean,
      default: true
    },
    showConfirm: {
      type: Boolean,
      default: true
    },
    confirmText: {
      type: String,
      default: 'Confirm'
    },
    confirmVariant: {
      type: String,
      default: 'primary'
    }
  },
  methods: {
    closeModal() {
      this.$emit('close')
    },
    confirm() {
      this.$emit('confirm')
    }
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-container {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid #eee;
}

.modal-content {
  padding: 1rem;
  max-height: 60vh;
  overflow-y: auto;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
  padding: 1rem;
  border-top: 1px solid #eee;
}

.close-icon {
  cursor: pointer;
}
</style>