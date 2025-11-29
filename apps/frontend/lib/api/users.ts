import { apiClient } from './config';
import { User } from './types';

export const usersApi = {
  // Получить пользователя по ID
  getById: async (userId: string): Promise<User> => {
    const response = await apiClient.get<User>(`/private/users/${userId}`);
    return response.data;
  },
};

