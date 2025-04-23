<template>
  <div class="login-page">
    <div class="login-container">
      <h1>e-Arkive Login</h1>
      
      <form @submit.prevent="handleLogin" class="login-form">
        <InputField
          label="Användarnamn"
          v-model="username"
          :disabled="loading"
          placeholder="Ange användarnamn"
        />
        
        <InputField
          label="Lösenord"
          type="password"
          v-model="password"
          :disabled="loading"
          placeholder="Ange lösenord"
        />
        
        <div v-if="statusMessage" class="status-message" :class="{ 'error': isError }">
          {{ statusMessage }}
        </div>
        
        <Button
          type="submit"
          variant="primary"
          :loading="loading"
          :disabled="loading"
          fullWidth
        >
          Logga in
        </Button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuth } from '../../composables/useAuth';
import Button from '../atoms/Button.vue';
import InputField from '../atoms/InputField.vue';

const router = useRouter();
const { login } = useAuth();

const username = ref('');
const password = ref('');
const loading = ref(false);
const statusMessage = ref('');
const isError = ref(false);

const handleLogin = async () => {
  if (!username.value || !password.value) {
    statusMessage.value = 'Vänligen fyll i alla fält';
    isError.value = true;
    return;
  }

  loading.value = true;
  statusMessage.value = '';
  isError.value = false;

  try {
    console.log('Attempting login with:', username.value);
    const result = await login(username.value, password.value);
    console.log('Login result:', result);
    
    if (result.success) {
      console.log('Login successful, redirecting to list page');
      router.push('/list');
    } else {
      statusMessage.value = result.error || 'Fel användarnamn eller lösenord';
      isError.value = true;
    }
  } catch (error) {
    console.error('Login error:', error);
    statusMessage.value = 'Ett fel uppstod vid inloggning';
    isError.value = true;
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 20px;
  background-color: var(--background-color);
}

.login-container {
  width: 100%;
  max-width: 400px;
  padding: 30px;
  border-radius: 8px;
  background-color: var(--surface-color);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

h1 {
  text-align: center;
  margin-bottom: 30px;
  color: var(--text-color);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.status-message {
  padding: 10px;
  border-radius: 4px;
  background-color: var(--success-background);
  color: var(--success-color);
  text-align: center;
}

.status-message.error {
  background-color: var(--error-background);
  color: var(--error-color);
}
</style>
