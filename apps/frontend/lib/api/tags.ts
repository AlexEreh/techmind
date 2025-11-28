import { apiClient } from './config';
import { Tag } from './types';

export const tagsApi = {
  // Get all tags for a company
  getByCompany: async (companyId: string): Promise<{ tags: Tag[]; total: number }> => {
    const response = await apiClient.get(`/private/document-tags/company/${companyId}`);
    return response.data;
  },

  // Get tags for a document
  getByDocument: async (documentId: string): Promise<{ tags: Tag[]; total: number }> => {
    const response = await apiClient.get(`/private/document-tags/document/${documentId}`);
    return response.data;
  },

  // Get tag by ID
  getById: async (id: string): Promise<Tag> => {
    const response = await apiClient.get(`/private/document-tags/tags/${id}`);
    return response.data;
  },

  // Create tag
  create: async (data: { company_id: string; name: string }): Promise<Tag> => {
    const response = await apiClient.post('/private/document-tags/tags', data);
    return response.data;
  },

  // Update tag
  update: async (id: string, data: { name: string }): Promise<Tag> => {
    const response = await apiClient.put(`/private/document-tags/tags/${id}`, data);
    return response.data;
  },

  // Delete tag
  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/private/document-tags/tags/${id}`);
  },

  // Add tag to document
  addToDocument: async (documentId: string, tagId: string): Promise<{ message: string }> => {
    const response = await apiClient.post('/private/document-tags/add', {
      document_id: documentId,
      tag_id: tagId
    });
    return response.data;
  },

  // Remove tag from document
  removeFromDocument: async (documentId: string, tagId: string): Promise<{ message: string }> => {
    const response = await apiClient.post('/private/document-tags/remove', {
      document_id: documentId,
      tag_id: tagId
    });
    return response.data;
  },
};

