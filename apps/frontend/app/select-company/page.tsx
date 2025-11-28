'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import { Card, CardBody, CardHeader } from '@heroui/card';
import { Button } from '@heroui/button';
import { Input } from '@heroui/input';
import { Spinner } from '@heroui/spinner';
import { Divider } from '@heroui/divider';
import { Modal, ModalContent, ModalHeader, ModalBody, ModalFooter, useDisclosure } from '@heroui/modal';
import { companyApi } from '@/lib/api/company';
import { PlusIcon } from '@/components/icons';

export default function SelectCompanyPage() {
  const { user, companies, currentCompany, setCurrentCompany, reloadCompanies, isLoading: authLoading } = useAuth();
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [isCreating, setIsCreating] = useState(false);
  const [newCompanyName, setNewCompanyName] = useState('');
  const { isOpen, onOpen, onClose } = useDisclosure();

  useEffect(() => {
    if (!authLoading && !user) {
      router.push('/login');
      return;
    }

    // Если компания уже выбрана, перенаправляем на файлы
    if (currentCompany) {
      router.push('/files');
    }
  }, [user, currentCompany, authLoading, router]);

  const handleSelectCompany = async (companyId: string) => {
    setIsLoading(true);
    const company = companies.find(c => c.id === companyId);
    if (company) {
      setCurrentCompany(company);
      router.push('/files');
    }
    setIsLoading(false);
  };

  const handleCreateCompany = async () => {
    if (!newCompanyName.trim()) return;

    setIsCreating(true);
    try {
      const newCompany = await companyApi.createCompany(newCompanyName.trim());
      await reloadCompanies();
      setNewCompanyName('');
      onClose();

      // Автоматически выбираем созданную компанию
      setCurrentCompany(newCompany);
      router.push('/files');
    } catch (error) {
      console.error('Failed to create company:', error);
      alert('Не удалось создать компанию. Попробуйте еще раз.');
    } finally {
      setIsCreating(false);
    }
  };

  if (authLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <Spinner size="lg" />
      </div>
    );
  }

  if (!user) {
    return null;
  }

  return (
    <div className="flex items-center justify-center min-h-screen bg-background p-4">
      <Card className="w-full max-w-2xl">
        <CardHeader className="flex flex-col gap-1 items-center pt-6">
          <h1 className="text-2xl font-bold">Выберите компанию</h1>
          <p className="text-small text-default-500">
            Выберите компанию для работы с документами или создайте новую
          </p>
        </CardHeader>
        <Divider />
        <CardBody className="gap-4 p-6">
          {companies.length === 0 ? (
            <div className="text-center py-8">
              <p className="text-default-500 mb-4">
                У вас пока нет доступа ни к одной компании
              </p>
              <p className="text-sm text-default-400 mb-6">
                Создайте новую компанию или попросите администратора пригласить вас
              </p>
              <Button
                color="primary"
                startContent={<PlusIcon />}
                onPress={onOpen}
                size="lg"
              >
                Создать компанию
              </Button>
            </div>
          ) : (
            <>
              <div className="flex justify-between items-center mb-2">
                <h3 className="text-lg font-semibold">Ваши компании</h3>
                <Button
                  color="primary"
                  variant="flat"
                  startContent={<PlusIcon />}
                  onPress={onOpen}
                  size="sm"
                >
                  Создать компанию
                </Button>
              </div>
              <div className="grid gap-3">
                {companies.map((company) => (
                  <Card
                    key={company.id}
                    isPressable
                    isHoverable
                    onPress={() => handleSelectCompany(company.id)}
                    className="border-2 border-transparent hover:border-primary"
                  >
                    <CardBody className="flex flex-row items-center justify-between p-4">
                      <div>
                        <h3 className="text-lg font-semibold">{company.name}</h3>
                        <p className="text-sm text-default-400">ID: {company.id}</p>
                      </div>
                      <Button
                        color="primary"
                        variant="flat"
                        onPress={() => handleSelectCompany(company.id)}
                        isLoading={isLoading}
                      >
                        Выбрать
                      </Button>
                    </CardBody>
                  </Card>
                ))}
              </div>
            </>
          )}
        </CardBody>
      </Card>

      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalContent>
          <ModalHeader>Создать новую компанию</ModalHeader>
          <ModalBody>
            <Input
              label="Название компании"
              placeholder="Введите название компании"
              value={newCompanyName}
              onChange={(e) => setNewCompanyName(e.target.value)}
              autoFocus
              onKeyDown={(e) => {
                if (e.key === 'Enter' && newCompanyName.trim()) {
                  handleCreateCompany();
                }
              }}
            />
            <p className="text-sm text-default-500">
              После создания вы станете администратором этой компании
            </p>
          </ModalBody>
          <ModalFooter>
            <Button variant="light" onPress={onClose}>
              Отмена
            </Button>
            <Button
              color="primary"
              onPress={handleCreateCompany}
              isLoading={isCreating}
              isDisabled={!newCompanyName.trim()}
            >
              Создать
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </div>
  );
}

