<template>
  <div class="user-card" @click="$emit('click', user)">
    <div class="user-card-header">
      <div class="user-avatar">
        <span class="material-icons">person</span>
      </div>
      <div class="user-info">
        <h4 class="user-name">{{ user.name }}</h4>
        <span class="user-username">@{{ user.username }}</span>
      </div>
    </div>
    
    <div class="user-groups">
      <h5 class="groups-heading">Groups</h5>
      <div v-if="user.groups && user.groups.length > 0" class="groups-list">
        <Badge 
          v-for="group in user.groups" 
          :key="group.id" 
          :text="group.name"
          variant="info"
        />
      </div>
      <div v-else class="no-groups">
        Not a member of any groups
      </div>
    </div>
    
    <div class="user-actions">
      <button @click.stop="$emit('edit', user)" class="edit-btn">
        <span class="material-icons">edit</span>
        Edit
      </button>
      <button @click.stop="$emit('delete', user)" class="delete-btn">
        <span class="material-icons">delete</span>
        Delete
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import Badge from '../atoms/Badge.vue';
import { User } from '../../types/schema';

defineProps<{
  user: User;
}>();

defineEmits(['click', 'edit', 'delete']);
</script>

<style scoped>
.user-card {
  background-color: var(--surface-color);
  border-radius: 8px;
  border: 1px solid var(--border-color);
  padding: 16px;
  transition: all 0.2s;
  cursor: pointer;
}

.user-card:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.user-card-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.user-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background-color: var(--primary-color-light);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
}

.user-avatar .material-icons {
  font-size: 24px;
  color: var(--primary-color);
}

.user-name {
  margin: 0 0 4px 0;
  font-size: 16px;
  font-weight: bold;
  color: var(--text-color);
}

.user-username {
  font-size: 14px;
  color: var(--text-color-secondary);
}

.user-groups {
  margin-bottom: 16px;
}

.groups-heading {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: var(--text-color-secondary);
  font-weight: normal;
}

.groups-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.no-groups {
  font-size: 14px;
  color: var(--text-color-secondary);
  font-style: italic;
}

.user-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 12px;
}

.edit-btn, .delete-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.edit-btn {
  background-color: var(--surface-color);
  color: var(--text-color);
  border: 1px solid var(--border-color);
}

.edit-btn:hover {
  background-color: var(--hover-color);
}

.delete-btn {
  background-color: var(--surface-color);
  color: var(--error-color);
  border: 1px solid var(--error-color);
}

.delete-btn:hover {
  background-color: rgba(231, 76, 60, 0.1);
}

.material-icons {
  font-size: 18px;
}
</style>
