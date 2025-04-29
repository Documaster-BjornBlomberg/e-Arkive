<template>
  <div class="groups-list">
    <div class="section-header">
      <h3>Groups Management</h3>
      <Button @click="showCreateGroupModal = true" icon="group_add">Add Group</Button>
    </div>
    
    <div class="search-container">
      <input 
        type="text" 
        v-model="searchQuery"
        class="search-input" 
        placeholder="Search groups..." 
      />
    </div>
    
    <div v-if="loading" class="loading">
      Loading groups...
    </div>
    
    <div v-else-if="filteredGroups.length === 0" class="empty-state">
      No groups found.
    </div>
    
    <div v-else class="groups-grid">
      <GroupCard 
        v-for="group in filteredGroups" 
        :key="group.id" 
        :group="group"
        @edit="openEditGroupModal"
        @delete="confirmDeleteGroup"
      />
    </div>
    
    <!-- Create Group Modal -->
    <Modal v-if="showCreateGroupModal" @close="showCreateGroupModal = false">
      <template #header>
        <h3>Create New Group</h3>
      </template>
      
      <template #body>
        <div class="form-group">
          <label for="newGroupName">Group Name</label>
          <input type="text" id="newGroupName" v-model="newGroupName" />
        </div>
      </template>
      
      <template #footer>
        <button 
          @click="createGroup" 
          class="save-btn" 
          :disabled="!newGroupName.trim()"
        >
          Create
        </button>
        <button @click="showCreateGroupModal = false" class="cancel-btn">Cancel</button>
      </template>
    </Modal>
    
    <!-- Edit Group Modal -->
    <Modal v-if="showEditGroupModal" @close="showEditGroupModal = false">
      <template #header>
        <h3>Edit Group</h3>
      </template>
      
      <template #body>
        <div class="form-group">
          <label for="editGroupName">Group Name</label>
          <input type="text" id="editGroupName" v-model="editingGroup.name" />
        </div>
        
        <div class="form-group">
          <label>Members</label>
          <div class="group-members">
            <div v-if="!users.length" class="no-users">
              No users available to add to group.
            </div>
            <div v-else-if="!editingGroup.memberIds.length && !nonMemberUsers.length" class="no-users">
              All users are already members of this group.
            </div>
            <div v-else>
              <div class="member-section">
                <h4>Current Members</h4>
                <div v-if="!editingGroup.memberIds.length" class="empty-message">
                  No members in this group yet.
                </div>
                <div v-else class="member-list">
                  <div v-for="userId in editingGroup.memberIds" :key="userId" class="member-item">
                    <span class="member-name">{{ getUserName(userId) }}</span>
                    <button @click="removeMember(userId)" class="remove-btn">
                      <span class="material-icons">remove_circle</span>
                    </button>
                  </div>
                </div>
              </div>
              
              <div v-if="nonMemberUsers.length" class="member-section">
                <h4>Add Members</h4>
                <div class="member-list">
                  <div v-for="user in nonMemberUsers" :key="user.id" class="member-item">
                    <span class="member-name">{{ user.name }} (@{{ user.username }})</span>
                    <button @click="addMember(user.id)" class="add-btn">
                      <span class="material-icons">add_circle</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
      
      <template #footer>
        <button 
          @click="updateGroup" 
          class="save-btn" 
          :disabled="!editingGroup.name.trim()"
        >
          Save Changes
        </button>
        <button @click="showEditGroupModal = false" class="cancel-btn">Cancel</button>
      </template>
    </Modal>
    
    <!-- Delete Confirmation Modal -->
    <Modal v-if="showDeleteGroupModal" @close="showDeleteGroupModal = false">
      <template #header>
        <h3>Confirm Delete</h3>
      </template>
      
      <template #body>
        <p>Are you sure you want to delete the group <strong>{{ groupToDelete?.name }}</strong>?</p>
        <p class="warning">This action cannot be undone.</p>
      </template>
      
      <template #footer>
        <button @click="deleteGroup" class="delete-btn">Delete</button>
        <button @click="showDeleteGroupModal = false" class="cancel-btn">Cancel</button>
      </template>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useUserManagement } from '../../composables/useUserManagement';
import GroupCard from '../molecules/GroupCard.vue';
import Modal from './Modal.vue';
import Button from '../atoms/Button.vue';
import { User, Group } from '../../types/schema';

// Composables
const { 
  loading, 
  error, 
  users,
  groups,
  getUsers, 
  getGroups,
  createGroup: apiCreateGroup,
  updateGroup: apiUpdateGroup,
  deleteGroup: apiDeleteGroup
} = useUserManagement();

// Local state
const searchQuery = ref('');
const showCreateGroupModal = ref(false);
const showEditGroupModal = ref(false);
const showDeleteGroupModal = ref(false);
const newGroupName = ref('');
const groupToDelete = ref<Group | null>(null);

interface EditingGroup {
  id: string;
  name: string;
  memberIds: string[];
  originalMemberIds: string[];
}

const editingGroup = ref<EditingGroup>({
  id: '',
  name: '',
  memberIds: [],
  originalMemberIds: []
});

// Computed
const filteredGroups = computed(() => {
  if (!searchQuery.value) return groups.value;
  
  const query = searchQuery.value.toLowerCase();
  return groups.value.filter(group => 
    group.name.toLowerCase().includes(query)
  );
});

const nonMemberUsers = computed(() => {
  return users.value.filter(user => 
    !editingGroup.value.memberIds.includes(user.id)
  );
});

// Initialize
onMounted(async () => {
  await Promise.all([getUsers(), getGroups()]);
});

// Methods
const getUserName = (userId: string): string => {
  const user = users.value.find(u => u.id === userId);
  return user ? `${user.name} (@${user.username})` : 'Unknown User';
};

const openEditGroupModal = (group: Group) => {
  const memberIds = group.members ? group.members.map(member => member.id) : [];
  
  editingGroup.value = {
    id: group.id,
    name: group.name,
    memberIds: [...memberIds],
    originalMemberIds: [...memberIds]
  };
  
  showEditGroupModal.value = true;
};

const confirmDeleteGroup = (group: Group) => {
  groupToDelete.value = group;
  showDeleteGroupModal.value = true;
};

const createGroup = async () => {
  try {
    await apiCreateGroup(newGroupName.value);
    
    // Refresh groups list
    await getGroups();
    
    // Reset form and close modal
    newGroupName.value = '';
    showCreateGroupModal.value = false;
  } catch (err) {
    console.error('Error creating group:', err);
    alert('Failed to create group. Please try again.');
  }
};

const updateGroup = async () => {
  try {
    // Update group name
    await apiUpdateGroup(editingGroup.value.id, editingGroup.value.name);
    
    // Handle membership changes
    const { addUserToGroup, removeUserFromGroup } = useUserManagement();
    
    // Add new members
    const membersToAdd = editingGroup.value.memberIds.filter(
      id => !editingGroup.value.originalMemberIds.includes(id)
    );
    
    for (const userId of membersToAdd) {
      await addUserToGroup(userId, editingGroup.value.id);
    }
    
    // Remove former members
    const membersToRemove = editingGroup.value.originalMemberIds.filter(
      id => !editingGroup.value.memberIds.includes(id)
    );
    
    for (const userId of membersToRemove) {
      await removeUserFromGroup(userId, editingGroup.value.id);
    }
    
    // Refresh groups list
    await getGroups();
    
    // Close modal
    showEditGroupModal.value = false;
  } catch (err) {
    console.error('Error updating group:', err);
    alert('Failed to update group. Please try again.');
  }
};

const deleteGroup = async () => {
  if (!groupToDelete.value) return;
  
  try {
    await apiDeleteGroup(groupToDelete.value.id);
    
    // Refresh groups list
    await getGroups();
    
    // Close modal
    showDeleteGroupModal.value = false;
    groupToDelete.value = null;
  } catch (err) {
    console.error('Error deleting group:', err);
    alert('Failed to delete group. Please try again.');
  }
};

const addMember = (userId: string) => {
  if (!editingGroup.value.memberIds.includes(userId)) {
    editingGroup.value.memberIds.push(userId);
  }
};

const removeMember = (userId: string) => {
  editingGroup.value.memberIds = editingGroup.value.memberIds.filter(id => id !== userId);
};
</script>

<style scoped>
.groups-list {
  max-width: 1000px;
  margin: 0 auto;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.section-header h3 {
  margin: 0;
}

.search-container {
  margin-bottom: 16px;
}

.search-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 16px;
  background-color: var(--input-background);
  color: var(--text-color);
}

.groups-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.loading, .empty-state {
  padding: 32px;
  text-align: center;
  color: var(--text-color-secondary);
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
}

.form-group input[type="text"] {
  width: 100%;
  padding: 10px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background-color: var(--input-background);
  color: var(--text-color);
}

.group-members {
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background-color: var(--input-background);
  max-height: 300px;
  overflow-y: auto;
}

.member-section {
  padding: 12px;
  border-bottom: 1px solid var(--border-color);
}

.member-section:last-child {
  border-bottom: none;
}

.member-section h4 {
  margin-top: 0;
  margin-bottom: 12px;
  font-size: 14px;
  color: var(--text-color-secondary);
}

.member-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.member-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px;
  border-radius: 4px;
  background-color: var(--surface-color);
}

.member-name {
  font-size: 14px;
}

.add-btn, .remove-btn {
  background: none;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-color-secondary);
  border-radius: 50%;
  width: 24px;
  height: 24px;
}

.add-btn:hover {
  color: var(--success-color);
}

.remove-btn:hover {
  color: var(--error-color);
}

.save-btn, .delete-btn, .cancel-btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  transition: opacity 0.2s;
}

.save-btn {
  background-color: var(--primary-color);
  color: white;
}

.delete-btn {
  background-color: var(--error-color);
  color: white;
}

.cancel-btn {
  background-color: var(--border-color);
  color: var(--text-color);
  margin-left: 8px;
}

.save-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.warning {
  color: var(--error-color);
}

.empty-message, .no-users {
  padding: 12px;
  text-align: center;
  font-style: italic;
  color: var(--text-color-secondary);
}
</style>
