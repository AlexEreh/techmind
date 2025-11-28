'use client';

import { useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import { Sidebar } from '@/components/files/Sidebar';
import { TagsManagement } from '@/components/settings/TagsManagement';
import { Spinner } from '@heroui/spinner';

export default function TagsPage() {
  const { user, currentCompany, isLoading: authLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!authLoading && !user) {
      router.push('/login');
    }
  }, [user, authLoading, router]);

  if (authLoading || !user || !currentCompany) {
    return (
      <div className="flex items-center justify-center h-screen">
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className="flex h-screen bg-background">
      <Sidebar currentView="tags" />

      <div className="flex-1 overflow-y-auto p-6">
        <h1 className="text-2xl font-bold mb-6">Управление тегами</h1>
        <TagsManagement />
      </div>
    </div>
  );
}

