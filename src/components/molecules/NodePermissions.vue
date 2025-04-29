<template>
  <div class="node-permissions">
    <h5>Permissions</h5>
    
    <div class="permissions-table">
      <div class="permission-row header">
        <div class="permission-cell"></div>
        <div class="permission-cell">Read</div>
        <div class="permission-cell">Write</div>
        <div class="permission-cell">Execute</div>
      </div>
      
      <div class="permission-row">
        <div class="permission-cell">Owner</div>
        <div class="permission-cell">
          <input 
            type="checkbox" 
            :checked="hasPermission(4)" 
            @change="togglePermission(4)" 
          />
        </div>
        <div class="permission-cell">
          <input 
            type="checkbox" 
            :checked="hasPermission(2)" 
            @change="togglePermission(2)" 
          />
        </div>
        <div class="permission-cell">
          <input 
            type="checkbox" 
            :checked="hasPermission(1)" 
            @change="togglePermission(1)" 
          />
        </div>
      </div>
      
      <div class="permission-row">
        <div class="permission-cell">Group</div>
        <div class="permission-cell">
          <input 
            type="checkbox" 
            :checked="hasPermission(4 << 3)" 
            @change="togglePermission(4 << 3)" 
          />
        </div>
        <div class="permission-cell">
          <input 
            type="checkbox" 
            :checked="hasPermission(2 << 3)" 
            @change="togglePermission(2 << 3)" 
          />
        </div>
        <div class="permission-cell">
          <input 
            type="checkbox" 
            :checked="hasPermission(1 << 3)" 
            @change="togglePermission(1 << 3)" 
          />
        </div>
      </div>
      
      <div class="permission-row">
        <div class="permission-cell">Others</div>
        <div class="permission-cell">
          <input 
            type="checkbox" 
            :checked="hasPermission(4 << 6)" 
            @change="togglePermission(4 << 6)" 
          />
        </div>
        <div class="permission-cell">
          <input 
            type="checkbox" 
            :checked="hasPermission(2 << 6)" 
            @change="togglePermission(2 << 6)" 
          />
        </div>
        <div class="permission-cell">
          <input 
            type="checkbox" 
            :checked="hasPermission(1 << 6)" 
            @change="togglePermission(1 << 6)" 
          />
        </div>
      </div>
    </div>
    
    <div class="permission-numeric">
      <span class="label">Numeric Value:</span>
      <span class="value">{{ modelValue.toString().padStart(3, '0') }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  modelValue: number;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: number): void;
}>();

// Check if a specific permission bit is set
const hasPermission = (bit: number): boolean => {
  return (props.modelValue & bit) === bit;
};

// Toggle a permission bit
const togglePermission = (bit: number): void => {
  const newValue = hasPermission(bit)
    ? props.modelValue & ~bit  // Turn bit off
    : props.modelValue | bit;  // Turn bit on
  
  emit('update:modelValue', newValue);
};
</script>

<style scoped>
.node-permissions {
  background-color: var(--background-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 16px;
}

h5 {
  margin-top: 0;
  margin-bottom: 16px;
  color: var(--text-color-secondary);
}

.permissions-table {
  border: 1px solid var(--border-color);
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 12px;
}

.permission-row {
  display: flex;
  border-bottom: 1px solid var(--border-color);
}

.permission-row:last-child {
  border-bottom: none;
}

.permission-row.header {
  background-color: var(--surface-color);
  font-weight: bold;
}

.permission-cell {
  flex: 1;
  padding: 8px 12px;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  border-right: 1px solid var(--border-color);
}

.permission-cell:first-child {
  text-align: left;
  justify-content: flex-start;
}

.permission-cell:last-child {
  border-right: none;
}

.permission-numeric {
  display: flex;
  gap: 8px;
  margin-top: 12px;
  font-family: monospace;
}

.label {
  font-weight: bold;
}

.value {
  background-color: var(--surface-color);
  padding: 2px 6px;
  border-radius: 4px;
  border: 1px solid var(--border-color);
}
</style>
