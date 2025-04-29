/**
 * GraphQL client composable
 * This provides a simple way to make GraphQL queries to our backend
 */
import { ref } from 'vue';
import { useAuth } from './useAuth';
import { 
  User, 
  File, 
  Node, 
  FileInput, 
  MetadataInput, 
  NodeInput, 
  NodeUpdateInput,
  UserSetting
} from '../types/schema';

interface GraphQLResponse<T> {
  data?: T;
  errors?: Array<{ message: string }>;
}

export function useGraphQL() {
  const loading = ref<boolean>(false);
  const error = ref<string | null>(null);
  const { token } = useAuth();

  // GraphQL endpoint URL - updated to match the running server
  const endpoint = 'http://localhost:8080/query';

  /**
   * Execute a GraphQL query
   * @param query - GraphQL query string
   * @param variables - Variables for the query
   * @returns Query result
   */
  const executeQuery = async <T>(query: string, variables: Record<string, any> = {}): Promise<T | null> => {
    loading.value = true;
    error.value = null;

    try {
      const headers: Record<string, string> = {
        'Content-Type': 'application/json',
      };

      if (token.value) {
        headers['Authorization'] = `Bearer ${token.value}`;
      }

      console.log(`Sending GraphQL request to ${endpoint}`);
      const response = await fetch(endpoint, {
        method: 'POST',
        headers,
        body: JSON.stringify({
          query,
          variables,
        }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result = await response.json() as GraphQLResponse<T>;

      // Check for GraphQL errors
      if (result.errors) {
        throw new Error(result.errors.map(e => e.message).join(', '));
      }

      return result.data || null;
    } catch (err: any) {
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
  const getFiles = async (): Promise<File[]> => {
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

    interface FilesResponse {
      getFiles: File[];
    }

    const data = await executeQuery<FilesResponse>(query);
    return data?.getFiles || [];
  };

  /**
   * Fetch files for a specific node
   * @param nodeId - The ID of the node to fetch files for
   * @returns The files in the node
   */
  const getFilesByNodeId = async (nodeId: string): Promise<File[]> => {
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

    interface FilesByNodeResponse {
      getFilesByNodeId: File[];
    }

    const data = await executeQuery<FilesByNodeResponse>(query, { nodeId });
    return data?.getFilesByNodeId || [];
  };

  /**
   * Get a single file by ID
   */
  const getFileById = async (id: string): Promise<File | null> => {
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

    interface FileResponse {
      getFile: File;
    }

    const data = await executeQuery<FileResponse>(query, { id });
    return data?.getFile || null;
  };

  /**
   * Save a new file to the backend
   * @param input - The file input data
   * @returns The saved file
   */
  const saveFile = async (input: FileInput): Promise<File | null> => {
    const mutation = `
      mutation SaveFile($input: FileInput!) {
        saveFile(input: $input) {
          id
          name
          size
          contentType
          createdAt
          nodeId
        }
      }
    `;

    interface SaveFileResponse {
      saveFile: File;
    }

    const data = await executeQuery<SaveFileResponse>(mutation, { input });
    return data?.saveFile || null;
  };

  /**
   * Update file metadata
   * @param fileId - The ID of the file to update
   * @param metadataInput - The metadata to update
   * @returns The updated file
   */
  const updateFileMetadata = async (fileId: string, metadataInput: MetadataInput[]): Promise<File | null> => {
    const mutation = `
      mutation UpdateMetadata($fileId: ID!, $metadataInput: [MetadataInput]!) {
        updateMetadata(fileId: $fileId, metadataInput: $metadataInput) {
          id
          name
          metadata {
            key
            value
          }
        }
      }
    `;

    interface UpdateMetadataResponse {
      updateMetadata: File;
    }

    const data = await executeQuery<UpdateMetadataResponse>(mutation, { fileId, metadataInput });
    return data?.updateMetadata || null;
  };

  /**
   * Delete file metadata
   * @param fileId - The ID of the file to update
   * @param keys - The keys of the metadata to delete
   * @returns The updated file
   */
  const deleteFileMetadata = async (fileId: string, keys: string[]): Promise<File | null> => {
    const mutation = `
      mutation DeleteMetadata($fileId: ID!, $keys: [String!]!) {
        deleteMetadata(fileId: $fileId, keys: $keys) {
          id
          name
          metadata {
            key
            value
          }
        }
      }
    `;

    interface DeleteMetadataResponse {
      deleteMetadata: File;
    }

    const data = await executeQuery<DeleteMetadataResponse>(mutation, { fileId, keys });
    return data?.deleteMetadata || null;
  };

  /**
   * Delete a file
   * @param id - The ID of the file to delete
   * @returns Success status
   */
  const deleteFile = async (id: string): Promise<boolean> => {
    const mutation = `
      mutation DeleteFile($id: ID!) {
        deleteFile(id: $id)
      }
    `;

    interface DeleteFileResponse {
      deleteFile: boolean;
    }

    const data = await executeQuery<DeleteFileResponse>(mutation, { id });
    return !!data?.deleteFile;
  };

  /**
   * Move a file to a new node
   * @param fileId - The ID of the file to move
   * @param nodeId - The ID of the destination node
   * @returns The moved file
   */
  const moveFile = async (fileId: string, nodeId: string): Promise<File | null> => {
    const mutation = `
      mutation MoveFile($fileId: ID!, $nodeId: ID!) {
        moveFile(fileId: $fileId, nodeId: $nodeId) {
          id
          name
          nodeId
        }
      }
    `;

    interface MoveFileResponse {
      moveFile: File;
    }

    const data = await executeQuery<MoveFileResponse>(mutation, { fileId, nodeId });
    return data?.moveFile || null;
  };

  /**
   * Fetch all root nodes from the backend (nodes with no parent)
   * @returns The root nodes
   */
  const getRootNodes = async (): Promise<Node[]> => {
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

    interface RootNodesResponse {
      getRootNodes: Node[];
    }

    const data = await executeQuery<RootNodesResponse>(query);
    return data?.getRootNodes || [];
  };

  /**
   * Get a specific node by ID
   * @param id - The ID of the node to fetch
   * @returns The node
   */
  const getNodeById = async (id: string): Promise<Node | null> => {
    const query = `
      query GetNodeById($id: ID!) {
        getNodeById(id: $id) {
          id
          name
          parentId
          createdAt
          updatedAt
        }
      }
    `;

    interface NodeResponse {
      getNodeById: Node;
    }

    const data = await executeQuery<NodeResponse>(query, { id });
    return data?.getNodeById || null;
  };

  /**
   * Fetch all children for a specific node
   * @param parentId - The ID of the parent node
   * @returns The child nodes
   */
  const getChildNodes = async (parentId: string): Promise<Node[]> => {
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

    interface ChildNodesResponse {
      getChildNodes: Node[];
    }

    const data = await executeQuery<ChildNodesResponse>(query, { parentId });
    return data?.getChildNodes || [];
  };

  /**
   * Create a new node
   * @param input - The node input data
   * @returns The created node
   */
  const createNode = async (input: NodeInput): Promise<Node | null> => {
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

    interface CreateNodeResponse {
      createNode: Node;
    }

    const data = await executeQuery<CreateNodeResponse>(mutation, { input });
    return data?.createNode || null;
  };

  /**
   * Update a node
   * @param id - The ID of the node to update
   * @param input - The update data
   * @returns The updated node
   */
  const updateNode = async (id: string, input: NodeUpdateInput): Promise<Node | null> => {
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

    interface UpdateNodeResponse {
      updateNode: Node;
    }

    const data = await executeQuery<UpdateNodeResponse>(mutation, { id, input });
    return data?.updateNode || null;
  };

  /**
   * Delete a node
   * @param id - The ID of the node to delete
   * @returns Success status
   */
  const deleteNode = async (id: string): Promise<boolean> => {
    const mutation = `
      mutation DeleteNode($id: ID!) {
        deleteNode(id: $id)
      }
    `;

    interface DeleteNodeResponse {
      deleteNode: boolean;
    }

    const data = await executeQuery<DeleteNodeResponse>(mutation, { id });
    return !!data?.deleteNode;
  };

  /**
   * Get a user setting
   * @param key - The key of the setting to fetch
   * @returns The setting
   */
  const getUserSetting = async (key: string): Promise<UserSetting | null> => {
    const query = `
      query GetUserSetting($key: String!) {
        getUserSetting(key: $key) {
          id
          key
          value
          createdAt
          updatedAt
        }
      }
    `;

    interface UserSettingResponse {
      getUserSetting: UserSetting;
    }

    const data = await executeQuery<UserSettingResponse>(query, { key });
    return data?.getUserSetting || null;
  };

  /**
   * Save a user setting
   * @param key - The key of the setting
   * @param value - The value of the setting
   * @returns The saved setting
   */
  const saveUserSetting = async (key: string, value: string): Promise<UserSetting | null> => {
    const mutation = `
      mutation SaveUserSetting($key: String!, $value: String!) {
        saveUserSetting(key: $key, value: $value) {
          id
          key
          value
          createdAt
          updatedAt
        }
      }
    `;

    interface SaveUserSettingResponse {
      saveUserSetting: UserSetting;
    }

    const data = await executeQuery<SaveUserSettingResponse>(mutation, { key, value });
    return data?.saveUserSetting || null;
  };

  return {
    loading,
    error,
    executeQuery,
    getFiles,
    getFilesByNodeId,
    getFileById,
    saveFile,
    updateFileMetadata,
    deleteFileMetadata,
    deleteFile,
    moveFile,
    getRootNodes,
    getNodeById,
    getChildNodes,
    createNode,
    updateNode,
    deleteNode,
    getUserSetting,
    saveUserSetting
  };
}
