import { apiClient } from './config';
import { Company, CompanyUserWithDetails, MyCompaniesResponse, CompanyUserData } from './types';

export const companyApi = {
  // Get user's companies
  getUserCompanies: async (): Promise<MyCompaniesResponse> => {
    const response = await apiClient.get<MyCompaniesResponse>('/private/companies/my');
    return response.data;
  },

  // Get company users (with JOIN on users table)
  getCompanyUsers: async (companyId: string): Promise<{ users: CompanyUserWithDetails[]; total: number }> => {
    const response = await apiClient.get(`/private/companies/${companyId}/users`);
    return response.data;
  },

  // Update user role
  updateUserRole: async (companyUserId: string, role: number): Promise<void> => {
    await apiClient.put(`/private/company-users/${companyUserId}/role`, { role });
  },

  // Remove user from company
  removeUser: async (companyUserId: string): Promise<void> => {
    await apiClient.delete(`/private/company-users/${companyUserId}`);
  },

  // Invite user to company
  inviteUser: async (companyId: string, email: string, role: number): Promise<void> => {
    await apiClient.post(`/private/companies/${companyId}/invite`, { email, role });
  },
};

