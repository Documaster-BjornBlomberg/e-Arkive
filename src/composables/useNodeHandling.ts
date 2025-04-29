import { ref, reactive, computed, nextTick } from 'vue';
import { useGraphQL } from './useGraphQL';
import { Node, File } from '../types/schema';

interface FlatNode {
  id: string;
  name: string;
}

export function useNodeHandling() {
  const { 
    loading: nodeLoading, 
    error: nodeError, 
    getRootNodes, 
    getChildNodes, 
    getNodeById, 
    getFilesByNodeId, 
    createNode, 
    updateNode, 
    deleteNode 
  } = useGraphQL();
  
  // State refs
  const rootNodes = ref<Node[]>([]);
  const expandedNodeIds = reactive(new Set<string>());
  const loadedChildrenMap = reactive(new Map<string, Node[]>());
  const isNodePanelOpen = ref<boolean>(true);
  const selectedNodeId = ref<string>('1'); // Default to root node (ID: 1)
  const nodeFiles = ref<File[]>([]);
  const currentNode = ref<Node | null>(null);
  const allNodesFlat = ref<FlatNode[]>([]); // For upload dropdown

  // Recursive helper to find a node by ID in the tree
  const findNodeById = (nodes: Node[], nodeId: string): Node | null => {
    for (const node of nodes) {
      if (node.id === nodeId) {
        return node;
      }
      if (node.children && node.children.length > 0) {
        const found = findNodeById(node.children, nodeId);
        if (found) {
          return found;
        }
      }
    }
    return null;
  };

  // Fetch root nodes (top level folders)
  const fetchRootNodes = async (): Promise<Node[]> => {
    try {
      const nodes = await getRootNodes();
      rootNodes.value = nodes || [];
      return nodes;
    } catch (error) {
      console.error('Error fetching root nodes:', error);
      return [];
    }
  };

  // Fetch children for a specific node
  const fetchNodeChildren = async (nodeId: string): Promise<Node[]> => {
    if (loadedChildrenMap.has(nodeId)) {
      return loadedChildrenMap.get(nodeId) || [];
    }

    try {
      const children = await getChildNodes(nodeId);
      const fetchedChildren = children || [];
      loadedChildrenMap.set(nodeId, fetchedChildren);

      // Find the parent node in the tree and attach the children
      await nextTick(); // Ensure DOM updates if needed before searching
      const parentNode = findNodeById(rootNodes.value, nodeId);
      if (parentNode) {
        parentNode.children = fetchedChildren;
      } else {
        // This might happen if the node structure is complex or fetched out of order
        console.warn(`Parent node with ID ${nodeId} not found in rootNodes tree to attach children.`);
      }

      return fetchedChildren;
    } catch (error) {
      console.error(`Error fetching children for node ${nodeId}:`, error);
      return [];
    }
  };

  // Fetch files for the selected node
  const fetchNodeFiles = async (nodeId: string): Promise<File[]> => {
    try {
      const files = await getFilesByNodeId(nodeId);
      nodeFiles.value = files;
      return files;
    } catch (error) {
      console.error(`Error fetching files for node ${nodeId}:`, error);
      return [];
    }
  };

  // Fetch the complete node with all properties
  const fetchNodeById = async (nodeId: string): Promise<Node | null> => {
    try {
      const node = await getNodeById(nodeId);
      currentNode.value = node;
      return node;
    } catch (error) {
      console.error(`Error fetching node ${nodeId}:`, error);
      return null;
    }
  };
  
  // Toggle node expansion state
  const toggleNodeExpand = async (nodeId: string): Promise<Node[]> => {
    if (expandedNodeIds.has(nodeId)) {
      expandedNodeIds.delete(nodeId);
      return [];
    } else {
      expandedNodeIds.add(nodeId);
      // Fetch children if not already loaded and attach them
      const children = await fetchNodeChildren(nodeId);
      return children;
    }
  };
  
  // Check if a node is expanded
  const isNodeExpanded = (nodeId: string): boolean => {
    return expandedNodeIds.has(nodeId);
  };
  
  // Get children for a specific node
  const getNodeChildren = (nodeId: string): Node[] => {
    return loadedChildrenMap.get(nodeId) || [];
  };
  
  // Select a node and fetch its files and details
  const selectNode = async (nodeId: string): Promise<void> => {
    selectedNodeId.value = nodeId;
    await Promise.all([
      fetchNodeById(nodeId),
      fetchNodeFiles(nodeId),
      // Also expand the node to show its children
      !expandedNodeIds.has(nodeId) ? toggleNodeExpand(nodeId) : Promise.resolve()
    ]);
  };
  
  // Add a new node (folder)
  const addNewNode = async (nodeName: string, parentId?: string): Promise<Node | null> => {
    try {
      const newNode = await createNode({
        name: nodeName,
        parentId: parentId || undefined
      });
      
      if (!newNode) return null;
      
      // Refresh the parent's children list if it exists and attach
      if (parentId) {
        // Re-fetch and attach children for the parent
        loadedChildrenMap.delete(parentId); // Force re-fetch
        await fetchNodeChildren(parentId);
        // Ensure parent is expanded to show the new node
        expandedNodeIds.add(parentId);
      } else {
        // If it's a root node, refresh the root nodes list
        await fetchRootNodes();
      }
      
      return newNode;
    } catch (error) {
      console.error('Error creating node:', error);
      throw error;
    }
  };
  
  // Update a node's name
  const updateNodeName = async (nodeId: string, newName: string): Promise<Node | null> => {
    try {
      const updatedNode = await updateNode(nodeId, {
        name: newName
      });
      
      if (!updatedNode) return null;
      
      // Update local data
      if (currentNode.value && currentNode.value.id === nodeId) {
        currentNode.value.name = newName;
      }
      
      // Refresh parent node's children if needed and re-attach
      const parentId = updatedNode.parentId;
      if (parentId) {
        loadedChildrenMap.delete(parentId); // Force re-fetch
        await fetchNodeChildren(parentId);
      } else {
        // If it's a root node, refresh the root nodes list
        await fetchRootNodes();
        // Need to re-find the updated node in the new root list if it was a root node
        const nodeInRoot = findNodeById(rootNodes.value, nodeId);
        if (nodeInRoot) nodeInRoot.name = newName; // Update name in potentially new object instance
      }
      
      return updatedNode;
    } catch (error) {
      console.error('Error updating node:', error);
      throw error;
    }
  };
  
  // Delete a node
  const removeNode = async (nodeId: string): Promise<boolean> => {
    try {
      const success = await deleteNode(nodeId);
      
      if (!success) return false;
      
      // Remove from expanded nodes set
      expandedNodeIds.delete(nodeId);
      
      // Remove from loaded children map
      loadedChildrenMap.delete(nodeId);
      
      // If this was the selected node, select the parent or root
      if (selectedNodeId.value === nodeId) {
        const parentId = currentNode.value?.parentId || '1';
        selectedNodeId.value = parentId;
        await fetchNodeById(parentId);
        await fetchNodeFiles(parentId);
      }
      
      // Refresh the parent's children after deletion
      const parentIdToDeleteFrom = currentNode.value?.parentId;
      if (parentIdToDeleteFrom) {
        loadedChildrenMap.delete(parentIdToDeleteFrom); // Force re-fetch
        await fetchNodeChildren(parentIdToDeleteFrom);
      } else {
        // If a root node was deleted, refresh root nodes
        await fetchRootNodes();
      }
      
      return true;
    } catch (error) {
      console.error('Error deleting node:', error);
      throw error;
    }
  };
  
  // Initialize node handling
  const initialize = async (): Promise<void> => {
    // First, fetch the root nodes
    await fetchRootNodes();
    
    // Select the root node (ID: 1) - this will also fetch its children via selectNode
    await selectNode('1');
    
    // Ensure the root node is expanded to show its children
    expandedNodeIds.add('1');
  };

  // Helper function to recursively fetch all nodes and flatten them
  const fetchAllNodesRecursiveHelper = async (
    nodeId: string, 
    depth: number = 0, 
    flatList: FlatNode[] = []
  ): Promise<void> => {
    try {
      // Fetch the node itself to get its name
      const node = await getNodeById(nodeId);
      if (!node) return; // Skip if node fetch fails

      // Add the current node to the list with indentation
      flatList.push({
        id: node.id,
        name: `${'--'.repeat(depth)} ${node.name}`, // Indent based on depth
      });

      // Fetch children for the current node
      const children = await fetchNodeChildren(nodeId);
      if (children && children.length > 0) {
        // Recursively fetch for each child
        for (const child of children) {
          await fetchAllNodesRecursiveHelper(child.id, depth + 1, flatList);
        }
      }
    } catch (error) {
      console.error(`Error fetching node or children during recursive fetch for ${nodeId}:`, error);
      // Continue with other nodes even if one fails
    }
  };

  // Function to get all nodes formatted for the upload dropdown
  const getAllNodesForUpload = async (): Promise<FlatNode[]> => {
    allNodesFlat.value = []; // Clear previous list
    try {
      // Start recursion from the root node(s)
      const roots = await fetchRootNodes(); // Ensure roots are fetched first
      for (const rootNode of roots) {
        await fetchAllNodesRecursiveHelper(rootNode.id, 0, allNodesFlat.value);
      }
      return allNodesFlat.value;
    } catch (error) {
      console.error('Error getting all nodes for upload:', error);
      return [];
    }
  };
  
  // Return all the necessary functions and reactive data
  return {
    nodeLoading,
    nodeError,
    rootNodes,
    nodeFiles,
    currentNode,
    isNodePanelOpen,
    selectedNodeId,
    expandedNodeIds,
    fetchRootNodes,
    fetchNodeChildren,
    fetchNodeFiles,
    fetchNodeById,
    toggleNodeExpand,
    isNodeExpanded,
    getNodeChildren,
    selectNode,
    addNewNode,
    updateNodeName,
    removeNode,
    initialize,
    getAllNodesForUpload
  };
}
