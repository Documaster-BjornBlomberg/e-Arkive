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
          nodeId
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
   * Fetch files for a specific node
   * @param {string} nodeId - The ID of the node to fetch files for
   * @returns {Promise<Array>} The files in the node
   */
  const getFilesByNodeId = async (nodeId) => {
    const query = `
      query GetFilesByNodeId($nodeId: ID!) {
        getFilesByNodeId(nodeId: $nodeId) {
          id
          name
          size
          contentType
          createdAt
          nodeId
          metadata {
            key
            value
          }
        }
      }
    `;

    const data = await executeQuery(query, { nodeId });
    return data?.getFilesByNodeId || [];
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
          nodeId
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
          nodeId
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

  /**
   * Move a file to a different node
   * @param {string} fileId - The ID of the file to move
   * @param {string} nodeId - The ID of the destination node
   * @returns {Promise<Object>} The updated file
   */
  const moveFile = async (fileId, nodeId) => {
    const mutation = `
      mutation MoveFile($fileId: ID!, $nodeId: ID!) {
        moveFile(fileId: $fileId, nodeId: $nodeId) {
          id
          name
          nodeId
        }
      }
    `;

    const data = await executeQuery(mutation, { fileId, nodeId });
    return data?.moveFile;
  };

  /**
   * Fetch all root nodes from the backend (nodes with no parent)
   * @returns {Promise<Array>} The root nodes
   */
  const getRootNodes = async () => {
    const query = `
      query {
        getRootNodes {
          id
          name
          parentId
          createdAt
          updatedAt
        }
      }
    `;

    const data = await executeQuery(query);
    return data?.getRootNodes || [];
  };

  /**
   * Fetch all children for a specific node
   * @param {string} parentId - The ID of the parent node
   * @returns {Promise<Array>} The child nodes
   */
  const getChildNodes = async (parentId) => {
    const query = `
      query GetChildNodes($parentId: ID!) {
        getChildNodes(parentId: $parentId) {
          id
          name
          parentId
          createdAt
          updatedAt
        }
      }
    `;

    const data = await executeQuery(query, { parentId });
    return data?.getChildNodes || [];
  };

  /**
   * Create a new node
   * @param {Object} input - The node input data
   * @returns {Promise<Object>} The created node
   */
  const createNode = async (input) => {
    const mutation = `
      mutation CreateNode($input: NodeInput!) {
        createNode(input: $input) {
          id
          name
          parentId
          createdAt
          updatedAt
        }
      }
    `;

    const data = await executeQuery(mutation, { input });
    return data?.createNode;
  };

  /**
   * Get a node by ID
   * @param {string} id - The ID of the node
   * @returns {Promise<Object>} The node
   */
  const getNodeById = async (id) => {
    const query = `
      query GetNodeById($id: ID!) {
        getNodeById(id: $id) {
          id
          name
          parentId
          createdAt
          updatedAt
          children {
            id
            name
          }
          files {
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
      }
    `;

    const data = await executeQuery(query, { id });
    return data?.getNodeById;
  };

  /**
   * Update an existing node
   * @param {string} id - The ID of the node to update
   * @param {Object} input - The update input data
   * @returns {Promise<Object>} The updated node
   */
  const updateNode = async (id, input) => {
    const mutation = `
      mutation UpdateNode($id: ID!, $input: NodeUpdateInput!) {
        updateNode(id: $id, input: $input) {
          id
          name
          parentId
          createdAt
          updatedAt
        }
      }
    `;

    const data = await executeQuery(mutation, { id, input });
    return data?.updateNode;
  };

  /**
   * Delete a node
   * @param {string} id - The ID of the node to delete
   * @returns {Promise<boolean>} Success status
   */
  const deleteNode = async (id) => {
    const mutation = `
      mutation DeleteNode($id: ID!) {
        deleteNode(id: $id)
      }
    `;

    const data = await executeQuery(mutation, { id });
    return data?.deleteNode;
  };

  return {
    loading,
    error,
    executeQuery,
    getFiles,
    getFilesByNodeId,
    getFileById,
    saveFile,
    moveFile,
    getRootNodes,
    getChildNodes,
    getNodeById,
    createNode,
    updateNode,
    deleteNode
  };
}