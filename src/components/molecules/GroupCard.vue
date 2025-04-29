<template>
  <div class="group-card" @click="$emit('click', group)">
    <div class="group-card-header">
      <div class="group-icon">
        <span class="material-icons">group</span>
      </div>
      <div class="group-info">
        <h4 class="group-name">{{ group.name }}</h4>
        <span class="group-members-count">
          {{ group.members ? group.members.length : 0 }} members
        </span>
      </div>
    </div>
    
    <div class="group-members">
      <h5 class="members-heading">Members</h5>
      <div v-if="group.members && group.members.length > 0" class="members-list">
        <div v-for="member in displayedMembers" :key="member.id" class="member-item">
          <span class="material-icons">person</span>
          <span class="member-name">{{ member.name }}</span>
        </div>
        <div v-if="hasMoreMembers" class="more-members">
          And {{ group.members.length - maxDisplayedMembers }} more...
        </div>
      </div>
      <div v-else class="no-members">
        No members in this group
      </div>
    </div>
    
    <div class="group-actions">
      <button @click.stop="$emit('edit', group)" class="edit-btn">
        <span class="material-icons">edit</span>
        Edit
      </button>
      <button @click.stop="$emit('delete', group)" class="delete-btn">
        <span class="material-icons">delete</span>
        Delete
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { Group } from '../../types/schema';

const props = defineProps<{
  group: Group;
}>();

defineEmits(['click', 'edit', 'delete']);

const maxDisplayedMembers = 3;

const displayedMembers = computed(() => {
  if (!props.group.members) return [];
  return props.group.members.slice(0, maxDisplayedMembers);
});

const hasMoreMembers = computed(() => {
  return props.group.members && props.group.members.length > maxDisplayedMembers;
});
</script>

<style scoped>
.group-card {
  background-color: var(--surface-color);
  border-radius: 8px;
  border: 1px solid var(--border-color);
  padding: 16px;
  transition: all 0.2s;
  cursor: pointer;
}

.group-card:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.group-card-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.group-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background-color: var(--primary-color-light);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
}

.group-icon .material-icons {
  font-size: 24px;
  color: var(--primary-color);
}

.group-name {
  margin: 0 0 4px 0;
  font-size: 16px;
  font-weight: bold;
  color: var(--text-color);
}

.group-members-count {
  font-size: 14px;
  color: var(--text-color-secondary);
}

.group-members {
  margin-bottom: 16px;
}

.members-heading {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: var(--text-color-secondary);
  font-weight: normal;
}

.members-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.member-item {
  display: flex;
  align-items: center;
  padding: 4px 0;
}

.member-item .material-icons {
  font-size: 16px;
  margin-right: 8px;
  color: var(--text-color-secondary);
}

.member-name {
  font-size: 14px;
}

.more-members {
  margin-top: 4px;
  font-size: 14px;
  color: var(--text-color-secondary);
  font-style: italic;
}

.no-members {
  font-size: 14px;
  color: var(--text-color-secondary);
  font-style: italic;
}

.group-actions {
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
