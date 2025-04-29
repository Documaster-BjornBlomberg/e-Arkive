import { ref, watch, onMounted } from 'vue';
import { useGraphQL } from './useGraphQL';

export function useTheme() {
  const { executeQuery, getUserSetting, saveUserSetting } = useGraphQL();
  const currentTheme = ref<string>('light');
  const isLoading = ref<boolean>(true);
  const error = ref<Error | null>(null);

  // Load theme from user settings
  const loadTheme = async (): Promise<void> => {
    isLoading.value = true;
    error.value = null;
    
    try {
      // Try to get theme from user settings
      const themeSetting = await getUserSetting('theme');
      
      if (themeSetting && themeSetting.value) {
        currentTheme.value = themeSetting.value;
      } else {
        // Check if we have a preferred theme stored in localStorage as fallback
        const savedTheme = localStorage.getItem('preferred-theme');
        if (savedTheme) {
          currentTheme.value = savedTheme;
        } else {
          // Check for system preference
          if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
            currentTheme.value = 'dark';
          }
        }
      }
      
      applyTheme(currentTheme.value);
    } catch (err: any) {
      console.error('Error loading theme:', err);
      error.value = err;
      
      // Fallback to localStorage
      const savedTheme = localStorage.getItem('preferred-theme');
      if (savedTheme) {
        currentTheme.value = savedTheme;
        applyTheme(currentTheme.value);
      }
    } finally {
      isLoading.value = false;
    }
  };
  
  // Save theme to user settings
  const saveTheme = async (theme: string): Promise<boolean> => {
    error.value = null;
    
    // First update local state and apply immediately for responsive feel
    currentTheme.value = theme;
    applyTheme(theme);
    
    // Save to localStorage as fallback
    localStorage.setItem('preferred-theme', theme);
    
    try {
      const setting = await saveUserSetting('theme', theme);
      return !!setting;
    } catch (err: any) {
      console.error('Error saving theme:', err);
      error.value = err;
      return false;
    }
  };
  
  // Toggle between light and dark theme
  const toggleTheme = (): void => {
    const newTheme = currentTheme.value === 'light' ? 'dark' : 'light';
    saveTheme(newTheme);
  };
  
  // Apply theme to document
  const applyTheme = (theme: string): void => {
    document.documentElement.setAttribute('data-theme', theme);
  };
  
  // Watch for theme changes to apply
  watch(currentTheme, (newTheme) => {
    applyTheme(newTheme);
  });
  
  onMounted(() => {
    // Apply default theme before loading
    applyTheme(currentTheme.value);
    loadTheme();
  });
  
  return {
    currentTheme,
    isLoading,
    error,
    loadTheme,
    saveTheme,
    toggleTheme
  };
}
