import { ref, reactive } from 'vue';
import { useGraphQL } from './useGraphQL';

export function useNodeHandling() {
  const { loading: nodeLoading, error: nodeError } = useGraphQL();
  const rootNodes = ref([]);
  const expandedNodeIds = reactive(new Set());
  const loadedChildrenMap = reactive(new Map());
  const isNodePanelOpen = ref(true);

  // Fetch root nodes (top level folders)
  const fetchRootNodes = async () => {
    try {
      const { getRootNodes } = useGraphQL();
      const nodes = await getRootNodes();
      rootNodes.value = nodes || [];
      return nodes;
    } catch (error) {
      console.error('Error fetching root nodes:', error);
      return [];
    }
  };

  // Fetch children for a specific node
  const fetchNodeChildren = async (nodeId) => {
    if (loadedChildrenMap.has(nodeId)) {
      return loadedChildrenMap.get(nodeId);
    }

    try {
      const { getChildNodes } = useGraphQL();
      const children = await getChildNodes(nodeId);
      loadedChildrenMap.set(nodeId, children || []);
      return children;
    } catch (error) {
      console.error(`Error fetching children for node ${nodeId}:`, error);
      return [];
    }
  };

  // Toggle node expanded state
  const toggleNodeExpanded = async (nodeId) => {
    if (expandedNodeIds.has(nodeId)) {
      expandedNodeIds.delete(nodeId);
    } else {
      expandedNodeIds.add(nodeId);
      // Load children if not already loaded
      if (!loadedChildrenMap.has(nodeId)) {
        await fetchNodeChildren(nodeId);
      }
    }
  };

  // Check if a node is expanded
  const isNodeExpanded = (nodeId) => {
    return expandedNodeIds.has(nodeId);
  };

  // Toggle node panel visibility
  const toggleNodePanel = () => {
    isNodePanelOpen.value = !isNodePanelOpen.value;
  };

  return {
    rootNodes,
    nodeLoading,
    nodeError,
    expandedNodeIds,
    loadedChildrenMap,
    isNodePanelOpen,
    fetchRootNodes,
    fetchNodeChildren,
    toggleNodeExpanded,
    isNodeExpanded,
    toggleNodePanel
  };
}