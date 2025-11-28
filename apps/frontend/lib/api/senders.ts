import { apiClient } from './config';
import { Sender } from './types';

export const sendersApi = {
  // Get all senders for a company
  getByCompany: async (companyId: string): Promise<{ senders: Sender[]; total: number }> => {
    const response = await apiClient.get(`/private/senders/company/${companyId}`);
    return response.data;
  },

  // Get sender by ID
  getById: async (id: string): Promise<Sender> => {
    const response = await apiClient.get(`/private/senders/${id}`);
    return response.data;
  },

  // Create sender
  create: async (data: { company_id: string; name: string; email?: string }): Promise<Sender> => {
    const response = await apiClient.post('/private/senders', data);
    return response.data;
  },

  // Update sender
  update: async (id: string, data: { name: string; email?: string }): Promise<Sender> => {
    const response = await apiClient.put(`/private/senders/${id}`, data);
    return response.data;
  },

  // Delete sender
  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/private/senders/${id}`);
  },
};

