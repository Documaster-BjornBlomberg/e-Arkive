/**
 * User and group management composable
 */
import { ref } from 'vue';
import { useGraphQL } from './useGraphQL';
import { User, Group } from '../types/schema';

export function useUserManagement() {
  const { executeQuery, loading, error } = useGraphQL();
  
  // State refs
  const users = ref<User[]>([]);
  const groups = ref<Group[]>([]);
  const selectedUser = ref<User | null>(null);
  const selectedGroup = ref<Group | null>(null);

  /**
   * Fetch all users
   */
  const getUsers = async (): Promise<User[]> => {
    const query = `
      query {
        getUsers {
          id
          name
          username
          groups {
            id
            name
          }
        }
      }
    `;

    interface UsersResponse {
      getUsers: User[];
    }

    const data = await executeQuery<UsersResponse>(query);
    if (data?.getUsers) {
      users.value = data.getUsers;
      return data.getUsers;
    }
    return [];
  };

  /**
   * Get a single user by ID
   */
  const getUserById = async (id: string): Promise<User | null> => {
    const query = `
      query GetUserById($id: ID!) {
        getUserById(id: $id) {
          id
          name
          username
          groups {
            id
            name
          }
        }
      }
    `;

    interface UserResponse {
      getUserById: User;
    }

    const data = await executeQuery<UserResponse>(query, { id });
    if (data?.getUserById) {
      selectedUser.value = data.getUserById;
      return data.getUserById;
    }
    return null;
  };

  /**
   * Fetch all groups
   */
  const getGroups = async (): Promise<Group[]> => {
    const query = `
      query {
        getGroups {
          id
          name
          members {
            id
            name
            username
          }
        }
      }
    `;

    interface GroupsResponse {
      getGroups: Group[];
    }

    const data = await executeQuery<GroupsResponse>(query);
    if (data?.getGroups) {
      groups.value = data.getGroups;
      return data.getGroups;
    }
    return [];
  };

  /**
   * Get a single group by ID
   */
  const getGroupById = async (id: string): Promise<Group | null> => {
    const query = `
      query GetGroup($id: ID!) {
        getGroup(id: $id) {
          id
          name
          members {
            id
            name
            username
          }
        }
      }
    `;

    interface GroupResponse {
      getGroup: Group;
    }

    const data = await executeQuery<GroupResponse>(query, { id });
    if (data?.getGroup) {
      selectedGroup.value = data.getGroup;
      return data.getGroup;
    }
    return null;
  };

  /**
   * Create a new group
   */
  const createGroup = async (name: string): Promise<Group | null> => {
    const mutation = `
      mutation CreateGroup($name: String!) {
        createGroup(name: $name) {
          id
          name
        }
      }
    `;

    interface CreateGroupResponse {
      createGroup: Group;
    }

    const data = await executeQuery<CreateGroupResponse>(mutation, { name });
    if (data?.createGroup) {
      // Refresh groups list
      await getGroups();
      return data.createGroup;
    }
    return null;
  };

  /**
   * Update a group
   */
  const updateGroup = async (id: string, name: string): Promise<Group | null> => {
    const mutation = `
      mutation UpdateGroup($id: ID!, $name: String!) {
        updateGroup(id: $id, name: $name) {
          id
          name
        }
      }
    `;

    interface UpdateGroupResponse {
      updateGroup: Group;
    }

    const data = await executeQuery<UpdateGroupResponse>(mutation, { id, name });
    if (data?.updateGroup) {
      // Refresh groups list
      await getGroups();
      return data.updateGroup;
    }
    return null;
  };

  /**
   * Delete a group
   */
  const deleteGroup = async (id: string): Promise<boolean> => {
    const mutation = `
      mutation DeleteGroup($id: ID!) {
        deleteGroup(id: $id)
      }
    `;

    interface DeleteGroupResponse {
      deleteGroup: boolean;
    }

    const data = await executeQuery<DeleteGroupResponse>(mutation, { id });
    if (data?.deleteGroup) {
      // Refresh groups list
      await getGroups();
      return true;
    }
    return false;
  };

  /**
   * Add a user to a group
   */
  const addUserToGroup = async (userId: string, groupId: string): Promise<boolean> => {
    const mutation = `
      mutation AddUserToGroup($userId: ID!, $groupId: ID!) {
        addUserToGroup(userId: $userId, groupId: $groupId)
      }
    `;

    interface AddUserToGroupResponse {
      addUserToGroup: boolean;
    }

    const data = await executeQuery<AddUserToGroupResponse>(mutation, { userId, groupId });
    return !!data?.addUserToGroup;
  };

  /**
   * Remove a user from a group
   */
  const removeUserFromGroup = async (userId: string, groupId: string): Promise<boolean> => {
    const mutation = `
      mutation RemoveUserFromGroup($userId: ID!, $groupId: ID!) {
        removeUserFromGroup(userId: $userId, groupId: $groupId)
      }
    `;

    interface RemoveUserFromGroupResponse {
      removeUserFromGroup: boolean;
    }

    const data = await executeQuery<RemoveUserFromGroupResponse>(mutation, { userId, groupId });
    return !!data?.removeUserFromGroup;
  };

  /**
   * Set node permissions
   */
  const setNodePermissions = async (nodeId: string, permissions: number): Promise<boolean> => {
    const mutation = `
      mutation SetNodePermissions($nodeId: ID!, $permissions: Int!) {
        setNodePermissions(nodeId: $nodeId, permissions: $permissions) {
          id
          permissions
        }
      }
    `;

    interface SetNodePermissionsResponse {
      setNodePermissions: {
        id: string;
        permissions: number;
      };
    }

    const data = await executeQuery<SetNodePermissionsResponse>(mutation, { nodeId, permissions });
    return !!data?.setNodePermissions;
  };

  /**
   * Set node ownership
   */
  const setNodeOwnership = async (nodeId: string, ownerUserId?: string, ownerGroupId?: string): Promise<boolean> => {
    const mutation = `
      mutation SetNodeOwnership($nodeId: ID!, $ownerUserId: ID, $ownerGroupId: ID) {
        setNodeOwnership(nodeId: $nodeId, ownerUserId: $ownerUserId, ownerGroupId: $ownerGroupId) {
          id
          ownerUserId
          ownerGroupId
        }
      }
    `;

    interface SetNodeOwnershipResponse {
      setNodeOwnership: {
        id: string;
        ownerUserId?: string;
        ownerGroupId?: string;
      };
    }

    const data = await executeQuery<SetNodeOwnershipResponse>(mutation, { nodeId, ownerUserId, ownerGroupId });
    return !!data?.setNodeOwnership;
  };
  /**
   * Create a new user
   */
  const createUser = async (userInput: { name: string; username: string; password: string }): Promise<User | null> => {
    const mutation = `
      mutation Register($username: String!, $password: String!) {
        register(username: $username, password: $password) {
          token
          user {
            id
            name
            username
          }
        }
      }
    `;

    interface RegisterResponse {
      register: {
        token: string;
        user: User;
      }
    }

    try {
      const data = await executeQuery<RegisterResponse>(mutation, { 
        username: userInput.username, 
        password: userInput.password 
      });
      
      if (data?.register?.user) {
        // Update the user's name in a separate mutation since register doesn't support name
        await updateUser({ 
          id: data.register.user.id, 
          name: userInput.name 
        });
        
        // Refresh users list
        await getUsers();
        return data.register.user;
      }
      return null;
    } catch (err) {
      console.error('Error creating user:', err);
      throw err;
    }
  };
  /**
   * Update user details
   */
  const updateUser = async (userInput: { id: string; name?: string; password?: string }): Promise<User | null> => {
    // For the name update
    const updateNameMutation = userInput.name ? `
      mutation UpdateUser($id: ID!, $name: String!) {
        updateUser(id: $id, name: $name) {
          id
          name
          username
        }
      }
    ` : null;

    // For the password update
    const updatePasswordMutation = userInput.password ? `
      mutation UpdatePassword($currentPassword: String!, $newPassword: String!) {
        updatePassword(currentPassword: $currentPassword, newPassword: $newPassword)
      }
    ` : null;

    interface UpdateUserResponse {
      updateUser: User;
    }

    try {
      let updatedUser: User | null = null;
          // Update name if provided
      if (updateNameMutation && userInput.name) {
        const data = await executeQuery<UpdateUserResponse>(updateNameMutation, { 
          id: userInput.id, 
          name: userInput.name 
        });
        
        if (data?.updateUser) {
          updatedUser = data.updateUser;
        }
      }
      
      // Update password if provided
      if (updatePasswordMutation && userInput.password) {
        // Note: This is a simplified implementation. In reality, we would need
        // the current password from the user to update the password
        await executeQuery(updatePasswordMutation, {
          currentPassword: 'currentPassword', // This would come from the user
          newPassword: userInput.password
        });
      }
      
      // Refresh users list
      await getUsers();
      return updatedUser;
    } catch (err) {
      console.error('Error updating user:', err);
      throw err;
    }
  };

  /**
   * Delete a user
   */
  const deleteUser = async (id: string): Promise<boolean> => {
    const mutation = `
      mutation DeleteUser($id: ID!) {
        deleteUser(id: $id)
      }
    `;

    interface DeleteUserResponse {
      deleteUser: boolean;
    }

    try {
      const data = await executeQuery<DeleteUserResponse>(mutation, { id });
      
      if (data?.deleteUser) {
        // Refresh users list
        await getUsers();
        return true;
      }
      return false;
    } catch (err) {
      console.error('Error deleting user:', err);
      throw err;
    }
  };

  return {
    loading,
    error,
    users,
    groups,
    selectedUser,
    selectedGroup,
    getUsers,
    getUserById,
    getGroups,
    getGroupById,
    createGroup,
    updateGroup,
    deleteGroup,
    addUserToGroup,
    removeUserFromGroup,
    setNodePermissions,
    setNodeOwnership,
    createUser,
    updateUser,
    deleteUser
  };
}
