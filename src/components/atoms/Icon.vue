<template>
  <span 
    class="icon"
    :class="[`icon-${name}`, { clickable }]"
    @click="clickable && $emit('click')"
    :title="title"
  >
    <i :class="iconClass"></i>
  </span>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  name: {
    type: String,
    required: true
  },
  size: {
    type: String,
    default: 'md' // options: sm, md, lg, xl
  },
  title: {
    type: String,
    default: ''
  },
  clickable: {
    type: Boolean,
    default: false
  }
});

const iconMappings = {
  'search': 'fas fa-search',
  'download': 'fas fa-download',
  'upload': 'fas fa-upload',
  'delete': 'fas fa-trash',
  'edit': 'fas fa-edit',
  'view': 'fas fa-eye',
  'close': 'fas fa-times',
  'add': 'fas fa-plus',
  'document': 'fas fa-file',
  'image': 'fas fa-file-image',
  'pdf': 'fas fa-file-pdf',
  'excel': 'fas fa-file-excel',
  'word': 'fas fa-file-word',
  'video': 'fas fa-file-video',
  'audio': 'fas fa-file-audio',
  'archive': 'fas fa-file-archive',
  'code': 'fas fa-file-code',
  'folder': 'fas fa-folder',
  'calendar': 'fas fa-calendar',
  'user': 'fas fa-user',
  'settings': 'fas fa-cog',
  'info': 'fas fa-info-circle',
  'warning': 'fas fa-exclamation-triangle',
  'success': 'fas fa-check-circle',
  'error': 'fas fa-times-circle'
};

const iconClass = computed(() => {
  const baseClass = iconMappings[props.name] || 'fas fa-question-circle';
  return [
    baseClass,
    `icon-size-${props.size}`
  ];
});
</script>

<style scoped>
.icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.icon.clickable {
  cursor: pointer;
  transition: color 0.2s;
}

.icon.clickable:hover {
  color: #4CAF50;
}

.icon-size-sm i {
  font-size: 0.875rem;
}

.icon-size-md i {
  font-size: 1rem;
}

.icon-size-lg i {
  font-size: 1.5rem;
}

.icon-size-xl i {
  font-size: 2rem;
}
</style>