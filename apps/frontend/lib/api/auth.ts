import { apiClient } from './config';
import { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse } from './types';
import axios from "axios";

export const authApi = {
    login: async (data: LoginRequest): Promise<LoginResponse> => {
        // Создаем отдельный клиент без Authorization для публичных запросов
        const response = await axios.post<LoginResponse>(
            `${process.env.NEXT_PUBLIC_API_URL || '/api'}/public/auth/login`,
            data
        );
        return response.data;
    },

    register: async (data: RegisterRequest): Promise<RegisterResponse> => {
        // Используем чистый axios без перехватчиков
        const response = await axios.post<RegisterResponse>(
            `${process.env.NEXT_PUBLIC_API_URL || '/api'}/public/auth/register`,
            data
        );
        return response.data;
    },

  logout: () => {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
  },
};

