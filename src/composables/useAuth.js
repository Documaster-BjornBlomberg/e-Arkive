import { ref, watch } from 'vue';
import { router } from '../main.js';

// Create shared auth state outside the composable function
const token = ref(localStorage.getItem('auth_token'));
const user = ref(null);
const isAuthenticated = ref(!!localStorage.getItem('auth_token'));

// GraphQL endpoint URL - updated to match the running server
const endpoint = 'http://localhost:8080/query';

export function useAuth() {
  const login = async (username, password) => {
    try {
      const result = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          query: `
            mutation Login($username: String!, $password: String!) {
              login(username: $username, password: $password) {
                token
                user {
                  id
                  name
                  username
                }
              }
            }
          `,
          variables: {
            username,
            password
          }
        })
      });

      const data = await result.json();
      
      if (data.errors) {
        throw new Error(data.errors[0].message);
      }

      const { token: authToken, user: userData } = data.data.login;
      token.value = authToken;
      user.value = userData;
      isAuthenticated.value = true;
      localStorage.setItem('auth_token', authToken);
      
      return { success: true };
    } catch (error) {
      return { success: false, error: error.message };
    }
  };

  const logout = () => {
    token.value = null;
    user.value = null;
    isAuthenticated.value = false;
    localStorage.removeItem('auth_token');
    router.push('/login');
  };

  const checkAuth = () => {
    const storedToken = localStorage.getItem('auth_token');
    if (storedToken) {
      token.value = storedToken;
      isAuthenticated.value = true;
      return true;
    } else {
      isAuthenticated.value = false;
      return false;
    }
  };

  // Watch for token changes and update authentication state
  watch(token, (newToken) => {
    isAuthenticated.value = !!newToken;
  });

  return {
    token,
    user,
    isAuthenticated,
    login,
    logout,
    checkAuth
  };
}
