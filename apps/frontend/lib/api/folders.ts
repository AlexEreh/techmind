import { apiClient } from './config';
import { Folder } from './types';

export const foldersApi = {
  // Get folders by parent (lazy loading support)
  getByParent: async (companyId: string, parentId?: string): Promise<{ folders: Folder[]; total: number }> => {
    const response = await apiClient.post('/private/folders/by-parent', {
      company_id: companyId,
      parent_id: parentId || null,
    });
    return response.data;
  },

  // Get all company folders (for search/overview)
  getAllCompanyFolders: async (companyId: string): Promise<{ folders: Folder[]; total: number }> => {
    const response = await apiClient.get(`/private/folders/company/${companyId}`);
    return response.data;
  },

  // Get specific folder
  getById: async (id: string): Promise<Folder> => {
    const response = await apiClient.get(`/private/folders/${id}`);
    return response.data;
  },

  // Create folder
  create: async (data: { company_id: string; name: string; parent_id?: string }): Promise<Folder> => {
    const response = await apiClient.post('/private/folders', data);
    return response.data;
  },

  // Rename folder
  rename: async (id: string, name: string): Promise<Folder> => {
    const response = await apiClient.put(`/private/folders/${id}/rename`, { name });
    return response.data;
  },

  // Delete folder
  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/private/folders/${id}`);
  },
};

