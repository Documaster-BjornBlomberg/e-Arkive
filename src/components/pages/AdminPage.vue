<template>
  <MainLayout>
    <div class="admin-page">
      <h1>Administration</h1>
      
      <div class="tabs">
        <tab-button 
          v-for="tab in tabs" 
          :key="tab.id"
          :active="activeTab === tab.id"
          @click="activeTab = tab.id"
        >
          {{ tab.name }}
        </tab-button>
      </div>
      
      <div class="tab-content">
        <!-- Users management tab -->
        <div v-if="activeTab === 'users'">
          <users-list />
        </div>
        
        <!-- Groups management tab -->
        <div v-if="activeTab === 'groups'">
          <groups-list />
        </div>
        
        <!-- Node management tab -->
        <div v-if="activeTab === 'nodes'">
          <node-management />
        </div>
        
        <!-- Settings tab -->
        <div v-if="activeTab === 'settings'">
          <div class="settings-section">
            <h3>System Settings</h3>
            <p>This section is under development.</p>
          </div>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuth } from '../../composables/useAuth';
import { useRouter } from 'vue-router';
import MainLayout from '../templates/MainLayout.vue';
import TabButton from '../atoms/TabButton.vue';
import UsersList from '../organisms/UsersList.vue';
import GroupsList from '../organisms/GroupsList.vue';
import NodeManagement from '../organisms/NodeManagement.vue';

// Router and auth setup
const router = useRouter();
const { isAuthenticated, checkAuth } = useAuth();

// Tabs configuration
const tabs = [
  { id: 'users', name: 'Users' },
  { id: 'groups', name: 'Groups' },
  { id: 'nodes', name: 'Node Access' },
  { id: 'settings', name: 'Settings' }
];

// Active tab state
const activeTab = ref('users');

// Check authentication on mount
onMounted(async () => {
  const authenticated = await checkAuth();
  if (!authenticated) {
    router.push('/login');
  }
});
</script>

<style scoped>
.admin-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

h1 {
  margin-bottom: 24px;
  color: var(--text-color);
}

.tabs {
  display: flex;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 24px;
  overflow-x: auto;
  white-space: nowrap;
}

.tab-content {
  padding-top: 8px;
}

.settings-section {
  background-color: var(--surface-color);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 20px;
}

.settings-section h3 {
  margin-top: 0;
  margin-bottom: 16px;
  color: var(--text-color);
}
</style>
