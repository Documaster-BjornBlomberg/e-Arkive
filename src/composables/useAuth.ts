import { ref, watch } from 'vue';
import { router } from '../main';
import { User } from '../types/schema';

interface LoginResponse {
  success: boolean;
  error?: string;
}

// Create shared auth state outside the composable function
const token = ref<string | null>(localStorage.getItem('auth_token'));
const user = ref<User | null>(null);
const isAuthenticated = ref<boolean>(!!localStorage.getItem('auth_token'));

// GraphQL endpoint URL - updated to match the running server
const endpoint = 'http://localhost:8080/query';

export function useAuth() {
  const login = async (username: string, password: string): Promise<LoginResponse> => {
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
    } catch (error: any) {
      return { success: false, error: error.message };
    }
  };

  const logout = async (): Promise<void> => {
    try {
      // Call the backend to invalidate the token if we have one
      if (token.value) {
        await fetch(endpoint, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            query: `
              mutation Logout($token: String!) {
                logout(token: $token)
              }
            `,
            variables: {
              token: token.value
            }
          })
        });
      }
    } catch (error) {
      console.error('Error during logout:', error);
    } finally {
      // Clear local state even if the server call fails
      token.value = null;
      user.value = null;
      isAuthenticated.value = false;
      localStorage.removeItem('auth_token');
      router.push('/login');
    }
  };

  const fetchCurrentUser = async (): Promise<User | null> => {
    if (!token.value) return null;
    
    try {
      const result = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token.value}`
        },
        body: JSON.stringify({
          query: `
            query {
              me {
                id
                name
                username
              }
            }
          `
        })
      });

      const data = await result.json();
      
      if (data.errors) {
        console.error('Error fetching current user:', data.errors);
        return null;
      }

      return data.data.me;
    } catch (error) {
      console.error('Error fetching current user:', error);
      return null;
    }
  };

  const checkAuth = async (): Promise<boolean> => {
    const storedToken = localStorage.getItem('auth_token');
    if (storedToken) {
      token.value = storedToken;
      isAuthenticated.value = true;
      
      // Fetch user information if we have a token
      const userData = await fetchCurrentUser();
      if (userData) {
        user.value = userData;
      } else {
        // If we couldn't get user data, the token might be invalid
        logout();
        return false;
      }
      
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
