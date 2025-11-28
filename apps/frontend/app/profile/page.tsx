'use client';

import { useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import { Sidebar } from '@/components/files/Sidebar';
import { Button } from '@heroui/button';
import { Card, CardBody, CardHeader } from '@heroui/card';
import { Divider } from '@heroui/divider';
import { Spinner } from '@heroui/spinner';

export default function ProfilePage() {
  const { user, currentCompany, logout, isLoading: authLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!authLoading && !user) {
      router.push('/login');
    }
  }, [user, authLoading, router]);

  const handleLogout = async () => {
    await logout();
    router.push('/login');
  };

  if (authLoading || !user) {
    return (
      <div className="flex items-center justify-center h-screen">
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className="flex h-screen bg-background">
      <Sidebar currentView="profile" />

      <div className="flex-1 overflow-y-auto p-6">
        <h1 className="text-2xl font-bold mb-6">Профиль пользователя</h1>

        <div className="max-w-2xl space-y-4">
          <Card>
            <CardHeader>
              <h2 className="text-xl font-semibold">Информация о пользователе</h2>
            </CardHeader>
            <Divider />
            <CardBody className="space-y-4">
              <div>
                <p className="text-sm text-default-500 mb-1">Имя</p>
                <p className="text-lg font-medium">{user.name}</p>
              </div>
              <div>
                <p className="text-sm text-default-500 mb-1">Email</p>
                <p className="text-lg font-medium">{user.email}</p>
              </div>
              <div>
                <p className="text-sm text-default-500 mb-1">ID пользователя</p>
                <p className="text-sm text-default-400 font-mono">{user.id}</p>
              </div>
            </CardBody>
          </Card>

          {currentCompany && (
            <Card>
              <CardHeader>
                <h2 className="text-xl font-semibold">Текущая компания</h2>
              </CardHeader>
              <Divider />
              <CardBody className="space-y-4">
                <div>
                  <p className="text-sm text-default-500 mb-1">Название</p>
                  <p className="text-lg font-medium">{currentCompany.name}</p>
                </div>
                <div>
                  <p className="text-sm text-default-500 mb-1">ID компании</p>
                  <p className="text-sm text-default-400 font-mono">{currentCompany.id}</p>
                </div>
              </CardBody>
            </Card>
          )}

          <Card>
            <CardHeader>
              <h2 className="text-xl font-semibold">Действия</h2>
            </CardHeader>
            <Divider />
            <CardBody>
              <Button
                color="danger"
                variant="flat"
                onPress={handleLogout}
                className="w-full"
              >
                Выйти из аккаунта
              </Button>
            </CardBody>
          </Card>
        </div>
      </div>
    </div>
  );
}

