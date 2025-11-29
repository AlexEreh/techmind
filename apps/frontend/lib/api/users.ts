import { apiClient } from './config';
import { User } from './types';
export const usersApi = {
  getById: async (id: string): Promise<User> => {
    const response = await apiClient.get(`/private/users/${id}`);
    return response.data;
  },
};
