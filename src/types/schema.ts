// GraphQL Schema Types

export interface User {
  id: string;
  name: string;
  username: string;
  settings?: UserSetting[];
  groups?: Group[];
}

export interface UserSetting {
  id: string;
  key: string;
  value: string;
  createdAt: string;
  updatedAt: string;
}

export interface Group {
  id: string;
  name: string;
  members?: User[];
}

export interface Node {
  id: string;
  name: string;
  parentId?: string | null;
  createdAt: string;
  updatedAt: string;
  children?: Node[];
  parent?: Node;
  files?: File[];
  ownerUserId?: string;
  ownerGroupId?: string;
  ownerUser?: User;
  ownerGroup?: Group;
  permissions?: number;
}

export interface File {
  id: string;
  name: string;
  size: number;
  contentType: string;
  createdAt: string;
  fileData?: string;
  metadata?: Metadata[];
  nodeId?: string | null;
  node?: Node;
}

export interface Metadata {
  key: string;
  value: string;
}

export interface FileInput {
  name: string;
  size: number;
  contentType: string;
  fileData: string;
  metadata?: MetadataInput[];
  nodeId?: string;
}

export interface MetadataInput {
  key: string;
  value: string;
}

export interface NodeInput {
  name: string;
  parentId?: string;
  ownerUserId?: string;
  ownerGroupId?: string;
  permissions?: number;
}

export interface NodeUpdateInput {
  name?: string;
  parentId?: string;
  ownerUserId?: string;
  ownerGroupId?: string;
  permissions?: number;
}

export interface AuthPayload {
  token: string;
  user: User;
}

export interface UserInput {
  name: string;
  username: string;
  password: string;
}

export interface UserUpdateInput {
  id: string;
  name?: string;
  password?: string;
}

export interface GroupInput {
  name: string;
}

export interface GroupUpdateInput {
  id: string;
  name: string;
}
