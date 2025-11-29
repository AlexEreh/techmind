export interface User {
  id: string;
  email: string;
  name: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  expires_at: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface RegisterResponse {
  token: string;
  user: User;
}

export interface Company {
  id: string;
  name: string;
}

export interface CompanyUser {
  id: string;
  user_id: string;
  company_id: string;
  role: number;
  user?: User;
}

export interface CompanyUserWithDetails {
  id: string; // user id
  name: string;
  email: string;
  role: number; // role in the company
  company_user_id: string; // id of the company_users record
}

export interface CompanyUserData {
  id: string;
  user_id: string;
  company_id: string;
  role: number;
  company?: Company;
}

export interface MyCompaniesResponse {
  companies: CompanyUserData[];
}

export const UserRole = {
  OWNER: 0,
  ADMIN: 1,
  MEMBER: 2,
  VIEWER: 3,
} as const;

export type UserRoleType = typeof UserRole[keyof typeof UserRole];

export interface Folder {
  id: string;
  company_id: string;
  parent_folder_id?: string;
  name: string;
  size: number;
  count: number;
  children?: Folder[];
}

export interface Tag {
  id: string;
  company_id: string;
  name: string;
}

export interface Sender {
  id: string;
  name: string;
  email?: string;
}

export interface Document {
  id: string;
  company_id: string;
  folder_id?: string;
  name: string;
  file_path: string;
  preview_file_path?: string;
  preview_url?: string;
  download_url?: string;
  file_size: number;
  mime_type: string;
  checksum: string;
  sender_id?: string;
  sender?: Sender;
  created_by?: string;
  updated_by?: string;
  created_at: string;
  updated_at: string;
  tags?: Tag[];
}

export interface FoldersTree {
  folders: Folder[];
  documents: Document[];
}

export interface SearchRequest {
  company_id: string;
  query?: string;
  folder_id?: string;
  tag_ids?: string[];
  sender_id?: string;
  page?: number;
  page_size?: number;
}

