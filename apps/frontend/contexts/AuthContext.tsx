'use client';

import React, { createContext, useContext, useState, useEffect } from 'react';
import { User, Company } from '@/lib/api/types';
import { authApi } from '@/lib/api/auth';
import { companyApi } from '@/lib/api/company';

interface AuthContextType {
  user: User | null;
  companies: Company[];
  currentCompany: Company | null;
  setCurrentCompany: (company: Company) => void;
  reloadCompanies: () => Promise<void>;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, name: string) => Promise<void>;
  logout: () => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [companies, setCompanies] = useState<Company[]>([]);
  const [currentCompany, setCurrentCompany] = useState<Company | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('token');
    const savedUser = localStorage.getItem('user');
    const savedCompany = localStorage.getItem('currentCompany');

    if (token && savedUser) {
      setUser(JSON.parse(savedUser));
      if (savedCompany) {
        setCurrentCompany(JSON.parse(savedCompany));
      }
      loadCompanies();
    }
    setIsLoading(false);
  }, []);

  const loadCompanies = async () => {
    try {
      const response = await companyApi.getUserCompanies();
      // Преобразуем CompanyUserData в Company, извлекая вложенную информацию о компании
      const userCompanies: Company[] = response.companies
        .filter(cu => cu.company) // Фильтруем только те, у которых есть данные компании
        .map(cu => cu.company!); // Извлекаем объект Company

      setCompanies(userCompanies);
      // Не устанавливаем компанию автоматически - пользователь должен выбрать сам
    } catch (error) {
      console.error('Failed to load companies:', error);
    }
  };

  const login = async (email: string, password: string) => {
    const response = await authApi.login({ email, password });
    localStorage.setItem('token', response.token);

    // Assuming the token contains user info or we need to fetch it
    const userData: User = { id: '', email, name: email }; // This should come from API
    setUser(userData);
    localStorage.setItem('user', JSON.stringify(userData));
    await loadCompanies();
  };

  const register = async (email: string, password: string, name: string) => {
    const response = await authApi.register({ email, password, name });
    localStorage.setItem('token', response.token);
    setUser(response.user);
    localStorage.setItem('user', JSON.stringify(response.user));
    await loadCompanies();
  };

  const logout = () => {
    setUser(null);
    setCompanies([]);
    setCurrentCompany(null);
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    localStorage.removeItem('currentCompany');
    authApi.logout();
  };

  const handleSetCurrentCompany = (company: Company) => {
    setCurrentCompany(company);
    localStorage.setItem('currentCompany', JSON.stringify(company));
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        companies,
        currentCompany,
        setCurrentCompany: handleSetCurrentCompany,
        reloadCompanies: loadCompanies,
        login,
        register,
        logout,
        isLoading,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
