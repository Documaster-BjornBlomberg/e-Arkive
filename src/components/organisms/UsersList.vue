<template>
  <div class="users-list">
    <div class="section-header">
      <h3>Users Management</h3>
      <Button @click="showCreateUserModal = true" icon="person_add">Add User</Button>
    </div>
    
    <div class="search-container">
      <input 
        type="text" 
        v-model="searchQuery"
        class="search-input" 
        placeholder="Search users..." 
      />
    </div>
    
    <div v-if="loading" class="loading">
      Loading users...
    </div>
    
    <div v-else-if="filteredUsers.length === 0" class="empty-state">
      No users found.
    </div>
    
    <div v-else class="users-grid">
      <UserCard 
        v-for="user in filteredUsers" 
        :key="user.id" 
        :user="user"
        @edit="openEditUserModal"
        @delete="confirmDeleteUser"
      />
    </div>
    
    <!-- Create User Modal -->
    <Modal v-if="showCreateUserModal" @close="showCreateUserModal = false">
      <template #header>
        <h3>Create New User</h3>
      </template>
      
      <template #body>
        <div class="form-group">
          <label for="newUserName">Name</label>
          <input type="text" id="newUserName" v-model="newUser.name" />
        </div>
        
        <div class="form-group">
          <label for="newUserUsername">Username</label>
          <input type="text" id="newUserUsername" v-model="newUser.username" />
        </div>
        
        <div class="form-group">
          <label for="newUserPassword">Password</label>
          <input type="password" id="newUserPassword" v-model="newUser.password" />
        </div>
        
        <div class="form-group">
          <label>Groups</label>
          <div class="groups-select">
            <div v-for="group in groups" :key="group.id" class="group-checkbox">
              <input 
                type="checkbox" 
                :id="'group-' + group.id" 
                :value="group.id"
                v-model="newUser.groupIds"
              />
              <label :for="'group-' + group.id">{{ group.name }}</label>
            </div>
          </div>
        </div>
      </template>
      
      <template #footer>
        <button 
          @click="handleCreateUser" 
          class="save-btn" 
          :disabled="!isNewUserValid()"
        >
          Create
        </button>
        <button @click="showCreateUserModal = false" class="cancel-btn">Cancel</button>
      </template>
    </Modal>
    
    <!-- Edit User Modal -->
    <Modal v-if="showEditUserModal" @close="showEditUserModal = false">
      <template #header>
        <h3>Edit User: {{ editingUser.name }}</h3>
      </template>
      
      <template #body>
        <div class="form-group">
          <label for="editUserName">Name</label>
          <input type="text" id="editUserName" v-model="editingUser.name" />
        </div>
        
        <div class="form-group">
          <label for="editUserUsername">Username</label>
          <input 
            type="text" 
            id="editUserUsername" 
            v-model="editingUser.username" 
            disabled 
            title="Username cannot be changed"
          />
        </div>
        
        <div class="form-group">
          <label for="editUserPassword">New Password (leave empty to keep current)</label>
          <input type="password" id="editUserPassword" v-model="editingUser.password" />
        </div>
        
        <div class="form-group">
          <label>Groups</label>
          <div class="groups-select">
            <div v-for="group in groups" :key="group.id" class="group-checkbox">
              <input 
                type="checkbox" 
                :id="'edit-group-' + group.id" 
                :value="group.id"
                v-model="editingUser.groupIds"
              />
              <label :for="'edit-group-' + group.id">{{ group.name }}</label>
            </div>
          </div>
        </div>
      </template>
      
      <template #footer>
        <button 
          @click="handleUpdateUser" 
          class="save-btn" 
          :disabled="!isEditedUserValid()"
        >
          Save Changes
        </button>
        <button @click="showEditUserModal = false" class="cancel-btn">Cancel</button>
      </template>
    </Modal>
    
    <!-- Delete Confirmation Modal -->
    <Modal v-if="showDeleteUserModal" @close="showDeleteUserModal = false">
      <template #header>
        <h3>Confirm Delete</h3>
      </template>
      
      <template #body>
        <p>Are you sure you want to delete the user <strong>{{ userToDelete?.name }}</strong>?</p>
        <p class="warning">This action cannot be undone.</p>
      </template>
      
      <template #footer>
        <button @click="handleDeleteUser" class="delete-btn">Delete</button>
        <button @click="showDeleteUserModal = false" class="cancel-btn">Cancel</button>
      </template>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useUserManagement } from '../../composables/useUserManagement';
import UserCard from '../molecules/UserCard.vue';
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
  createUser: apiCreateUser,
  updateUser: apiUpdateUser,
  deleteUser: apiDeleteUser,
  addUserToGroup,
  removeUserFromGroup
} = useUserManagement();

// Local state
const searchQuery = ref('');
const showCreateUserModal = ref(false);
const showEditUserModal = ref(false);
const showDeleteUserModal = ref(false);
const userToDelete = ref<User | null>(null);

interface NewUserForm {
  name: string;
  username: string;
  password: string;
  groupIds: string[];
}

interface EditUserForm extends NewUserForm {
  id: string;
}

const newUser = ref<NewUserForm>({
  name: '',
  username: '',
  password: '',
  groupIds: []
});

const editingUser = ref<EditUserForm>({
  id: '',
  name: '',
  username: '',
  password: '',
  groupIds: []
});

// Computed
const filteredUsers = computed(() => {
  if (!searchQuery.value) return users.value;
  
  const query = searchQuery.value.toLowerCase();
  return users.value.filter(user => 
    user.name.toLowerCase().includes(query) || 
    user.username.toLowerCase().includes(query)
  );
});

// Initialize
onMounted(async () => {
  await Promise.all([getUsers(), getGroups()]);
});

// Methods
const isNewUserValid = () => {
  return newUser.value.name && 
         newUser.value.username && 
         newUser.value.password;
};

const isEditedUserValid = () => {
  return editingUser.value.name.trim() !== '';
};

const resetNewUserForm = () => {
  newUser.value = {
    name: '',
    username: '',
    password: '',
    groupIds: []
  };
};

const openEditUserModal = (user: User) => {
  editingUser.value = {
    id: user.id,
    name: user.name,
    username: user.username,
    password: '',
    groupIds: user.groups ? user.groups.map(group => group.id) : []
  };
  showEditUserModal.value = true;
};

const confirmDeleteUser = (user: User) => {
  userToDelete.value = user;
  showDeleteUserModal.value = true;
};

const handleCreateUser = async () => {
  try {
    // Call API to create user
    const createdUser = await apiCreateUser({
      name: newUser.value.name,
      username: newUser.value.username,
      password: newUser.value.password
    });
    
    if (createdUser && newUser.value.groupIds.length > 0) {
      // Add user to selected groups
      await Promise.all(
        newUser.value.groupIds.map(groupId => 
          addUserToGroup(createdUser.id, groupId)
        )
      );
    }
    
    // Refresh users list
    await getUsers();
    
    // Reset form and close modal
    resetNewUserForm();
    showCreateUserModal.value = false;
  } catch (err) {
    console.error('Error creating user:', err);
    alert('Failed to create user. Please try again.');
  }
};

const handleUpdateUser = async () => {
  try {
    // Update basic user information
    await apiUpdateUser({
      id: editingUser.value.id,
      name: editingUser.value.name,
      password: editingUser.value.password || undefined
    });
    
    // Find the original user to compare groups
    const originalUser = users.value.find(u => u.id === editingUser.value.id);
    if (originalUser) {
      const originalGroupIds = originalUser.groups ? originalUser.groups.map(g => g.id) : [];
      const newGroupIds = editingUser.value.groupIds;
      
      // Remove user from groups they're no longer part of
      const groupsToRemove = originalGroupIds.filter(id => !newGroupIds.includes(id));
      await Promise.all(
        groupsToRemove.map(groupId => removeUserFromGroup(editingUser.value.id, groupId))
      );
      
      // Add user to new groups
      const groupsToAdd = newGroupIds.filter(id => !originalGroupIds.includes(id));
      await Promise.all(
        groupsToAdd.map(groupId => addUserToGroup(editingUser.value.id, groupId))
      );
    }
    
    // Refresh users list
    await getUsers();
    
    // Close modal
    showEditUserModal.value = false;
  } catch (err) {
    console.error('Error updating user:', err);
    alert('Failed to update user. Please try again.');
  }
};

const handleDeleteUser = async () => {
  if (!userToDelete.value) return;
  
  try {
    await apiDeleteUser(userToDelete.value.id);
    
    // Refresh users list
    await getUsers();
    
    // Close modal
    showDeleteUserModal.value = false;
    userToDelete.value = null;
  } catch (err) {
    console.error('Error deleting user:', err);
    alert('Failed to delete user. Please try again.');
  }
};
</script>

<style scoped>
.users-list {
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

.users-grid {
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
  color: var(--text-color);
}

.form-group input[type="text"],
.form-group input[type="password"] {
  width: 100%;
  padding: 10px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background-color: var(--input-background);
  color: var(--text-color);
}

.form-group input[disabled] {
  background-color: var(--hover-color);
  cursor: not-allowed;
}

.groups-select {
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 8px;
  background-color: var(--input-background);
}

.group-checkbox {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.group-checkbox:last-child {
  margin-bottom: 0;
}

.group-checkbox label {
  margin-left: 8px;
  margin-bottom: 0;
  cursor: pointer;
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
</style>
