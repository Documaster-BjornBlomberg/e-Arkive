/**
 * GraphQL client composable
 * This provides a simple way to make GraphQL queries to our backend
 */
import { ref } from 'vue';

export function useGraphQL() {
  const loading = ref(false);
  const error = ref(null);

  // GraphQL endpoint URL - updated to match the running server
  const endpoint = 'http://localhost:8080/query';

  /**
   * Execute a GraphQL query
   * @param {string} query - GraphQL query string
   * @param {Object} variables - Variables for the query
   * @returns {Promise<Object>} Query result
   */
  const executeQuery = async (query, variables = {}) => {
    loading.value = true;
    error.value = null;
    
    try {
      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          query,
          variables
        }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result = await response.json();
      
      // Check for GraphQL errors
      if (result.errors) {
        throw new Error(result.errors.map(e => e.message).join(', '));
      }

      return result.data;
    } catch (err) {
      console.error('GraphQL query error:', err);
      error.value = err.message;
      return null;
    } finally {
      loading.value = false;
    }
  };

  /**
   * Fetch all files from the backend
   */
  const getFiles = async () => {
    const query = `
      query {
        getFiles {
          id
          name
          size
          contentType
          createdAt
          metadata {
            key
            value
          }
        }
      }
    `;

    const data = await executeQuery(query);
    return data?.getFiles || [];
  };

  /**
   * Get a single file by ID
   */
  const getFileById = async (id) => {
    const query = `
      query GetFile($id: ID!) {
        getFile(id: $id) {
          id
          name
          size
          contentType
          createdAt
          metadata {
            key
            value
          }
        }
      }
    `;

    const data = await executeQuery(query, { id });
    return data?.getFile;
  };

  /**
   * Save a file to the backend
   * @param {Object} fileInput - The file input data
   * @returns {Promise<Object>} The saved file
   */
  const saveFile = async (fileInput) => {
    const mutation = `
      mutation SaveFile($input: FileInput!) {
        saveFile(input: $input) {
          id
          name
          size
          contentType
          createdAt
          metadata {
            key
            value
          }
        }
      }
    `;

    const data = await executeQuery(mutation, { input: fileInput });
    return data?.saveFile;
  };

  return {
    loading,
    error,
    executeQuery,
    getFiles,
    getFileById,
    saveFile
  };
}