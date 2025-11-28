import { apiClient } from './config';
import { Document, SearchRequest } from './types';

export const documentsApi = {
  // Search documents
  search: async (params: SearchRequest): Promise<{ documents: Document[]; total: number }> => {
    const response = await apiClient.post('/private/documents/search', params);
    return response.data;
  },

  // Get document by ID
  getById: async (id: string): Promise<Document> => {
    const response = await apiClient.get(`/private/documents/${id}`);
    return response.data;
  },

  // Get documents by folder
  getByFolder: async (folderId: string): Promise<{ documents: Document[]; total: number }> => {
    const response = await apiClient.get(`/private/documents/folder/${folderId}`);
    return response.data;
  },

  // Get all company documents
  getByCompany: async (companyId: string): Promise<{ documents: Document[]; total: number }> => {
    const response = await apiClient.get(`/private/documents/company/${companyId}`);
    return response.data;
  },

  // Update document
  update: async (id: string, data: { name?: string; folder_id?: string; sender_id?: string }): Promise<Document> => {
    const response = await apiClient.put(`/private/documents/${id}`, data);
    return response.data;
  },

  // Delete document
  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/private/documents/${id}`);
  },

  // Upload document
  upload: async (data: {
    company_id: string;
    name: string;
    folder_id?: string;
    sender_id?: string;
    file: File;
  }): Promise<Document> => {
    const formData = new FormData();
    formData.append('file', data.file);
    formData.append('company_id', data.company_id);
    formData.append('name', data.name);
    if (data.folder_id) {
      formData.append('folder_id', data.folder_id);
    }
    if (data.sender_id) {
      formData.append('sender_id', data.sender_id);
    }
    const response = await apiClient.post('/private/documents', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data;
  },

  // Get download URL
  getDownloadUrl: async (id: string): Promise<{ url: string; expires_at: string }> => {
    const response = await apiClient.get(`/private/documents/${id}/download`);
    return response.data;
  },

  // Get preview URL
  getPreviewUrl: async (id: string): Promise<{ url: string; expires_at: string }> => {
    const response = await apiClient.get(`/private/documents/${id}/preview`);
    return response.data;
  },
};

