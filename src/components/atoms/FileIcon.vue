<template>
  <div class="file-icon" :class="iconClass">
    <span>{{ fileExtension }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  fileName: {
    type: String,
    required: true
  },
  fileType: {
    type: String,
    default: ''
  }
});

const fileExtension = computed(() => {
  const parts = props.fileName.split('.');
  if (parts.length > 1) {
    return parts[parts.length - 1].toUpperCase();
  }
  return '';
});

const iconClass = computed(() => {
  const type = props.fileType.split('/')[0];
  
  if (type === 'image') {
    return 'image-file';
  } else if (type === 'application') {
    if (props.fileType.includes('pdf')) {
      return 'pdf-file';
    } else if (props.fileType.includes('word') || props.fileType.includes('document')) {
      return 'doc-file';
    } else if (props.fileType.includes('excel') || props.fileType.includes('spreadsheet')) {
      return 'sheet-file';
    }
    return 'application-file';
  } else if (type === 'text') {
    return 'text-file';
  } else if (type === 'video') {
    return 'video-file';
  } else if (type === 'audio') {
    return 'audio-file';
  }
  
  return 'generic-file';
});
</script>

<style scoped>
.file-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  color: white;
  font-size: 0.7rem;
  font-weight: bold;
}

.generic-file {
  background-color: #95a5a6;
}

.image-file {
  background-color: #3498db;
}

.pdf-file {
  background-color: #e74c3c;
}

.doc-file {
  background-color: #2980b9;
}

.sheet-file {
  background-color: #27ae60;
}

.text-file {
  background-color: #f39c12;
}

.video-file {
  background-color: #9b59b6;
}

.audio-file {
  background-color: #e67e22;
}

.application-file {
  background-color: #34495e;
}
</style>