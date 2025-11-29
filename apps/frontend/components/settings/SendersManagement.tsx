'use client';

import { useState, useEffect } from 'react';
import { Button } from '@heroui/button';
import { Input } from '@heroui/input';
import { Card, CardHeader, CardBody } from '@heroui/card';
import { Modal, ModalContent, ModalHeader, ModalBody, ModalFooter, useDisclosure } from '@heroui/modal';
import { Table, TableHeader, TableColumn, TableBody, TableRow, TableCell } from '@heroui/table';
import { Sender } from '@/lib/api/types';
import { sendersApi } from '@/lib/api/senders';
import { useAuth } from '@/contexts/AuthContext';
import { PlusIcon, TrashIcon, EditIcon } from '@/components/icons';
import { Spinner } from '@heroui/spinner';

export const SendersManagement: React.FC = () => {
  const [senders, setSenders] = useState<Sender[]>([]);
  const [newSenderName, setNewSenderName] = useState('');
  const [newSenderEmail, setNewSenderEmail] = useState('');
  const [editingSender, setEditingSender] = useState<Sender | null>(null);
  const [editName, setEditName] = useState('');
  const [editEmail, setEditEmail] = useState('');
  const [isCreating, setIsCreating] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const { currentCompany } = useAuth();
  const { isOpen: isCreateOpen, onOpen: onCreateOpen, onClose: onCreateClose } = useDisclosure();
  const { isOpen: isEditOpen, onOpen: onEditOpen, onClose: onEditClose } = useDisclosure();

  useEffect(() => {
    loadSenders();
  }, [currentCompany]);

  const loadSenders = async () => {
    if (!currentCompany) return;
    setIsLoading(true);
    try {
      const { senders: loadedSenders } = await sendersApi.getByCompany(currentCompany.id);
      setSenders(loadedSenders);
    } catch (error) {
      console.error('Failed to load senders:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateSender = async () => {
    if (!currentCompany || !newSenderName.trim()) return;

    setIsCreating(true);
    try {
      await sendersApi.create({
        company_id: currentCompany.id,
        name: newSenderName.trim(),
        email: newSenderEmail.trim() || undefined,
      });
      setNewSenderName('');
      setNewSenderEmail('');
      onCreateClose();
      await loadSenders();
    } catch (error) {
      console.error('Failed to create sender:', error);
    } finally {
      setIsCreating(false);
    }
  };

  const handleEditClick = (sender: Sender) => {
    setEditingSender(sender);
    setEditName(sender.name);
    setEditEmail(sender.email || '');
    onEditOpen();
  };

  const handleUpdateSender = async () => {
    if (!editingSender || !editName.trim()) return;

    setIsUpdating(true);
    try {
      await sendersApi.update(editingSender.id, {
        name: editName.trim(),
        email: editEmail.trim() || undefined,
      });
      onEditClose();
      setEditingSender(null);
      setEditName('');
      setEditEmail('');
      await loadSenders();
    } catch (error) {
      console.error('Failed to update sender:', error);
    } finally {
      setIsUpdating(false);
    }
  };

  const handleDeleteSender = async (senderId: string) => {
    if (!confirm('Вы уверены, что хотите удалить этого контрагента?')) return;

    try {
      await sendersApi.delete(senderId);
      await loadSenders();
    } catch (error) {
      console.error('Failed to delete sender:', error);
    }
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center p-8">
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className="py-4 space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-semibold">Управление контрагентами</h3>
        <Button color="default" startContent={<PlusIcon />} onPress={onCreateOpen}>
          Добавить контрагента
        </Button>
      </div>

      <Table aria-label="Таблица контрагентов">
        <TableHeader>
          <TableColumn>НАЗВАНИЕ</TableColumn>
          <TableColumn>EMAIL</TableColumn>
          <TableColumn>ДЕЙСТВИЯ</TableColumn>
        </TableHeader>
        <TableBody emptyContent="Контрагенты не созданы">
          {senders.map((sender) => (
            <TableRow key={sender.id}>
              <TableCell>{sender.name}</TableCell>
              <TableCell>{sender.email || '—'}</TableCell>
              <TableCell>
                <div className="flex gap-2">
                  <Button
                    isIconOnly
                    size="sm"
                    variant="light"
                    onPress={() => handleEditClick(sender)}
                  >
                    <EditIcon className="w-4 h-4" />
                  </Button>
                  <Button
                    isIconOnly
                    size="sm"
                    variant="light"
                    color="danger"
                    onPress={() => handleDeleteSender(sender.id)}
                  >
                    <TrashIcon className="w-4 h-4" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>

      {/* Create Modal */}
      <Modal isOpen={isCreateOpen} onClose={onCreateClose}>
        <ModalContent>
          <ModalHeader>Добавить контрагента</ModalHeader>
          <ModalBody>
            <Input
              label="Название"
              placeholder="Введите название"
              value={newSenderName}
              onChange={(e) => setNewSenderName(e.target.value)}
              isRequired
              autoFocus
            />
            <Input
              label="Email"
              placeholder="Введите email (необязательно)"
              value={newSenderEmail}
              onChange={(e) => setNewSenderEmail(e.target.value)}
              type="email"
            />
          </ModalBody>
          <ModalFooter>
            <Button variant="light" onPress={onCreateClose}>
              Отмена
            </Button>
            <Button
              color="default"
              onPress={handleCreateSender}
              isLoading={isCreating}
              isDisabled={!newSenderName.trim()}
            >
              Создать
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>

      {/* Edit Modal */}
      <Modal isOpen={isEditOpen} onClose={onEditClose}>
        <ModalContent>
          <ModalHeader>Редактировать контрагента</ModalHeader>
          <ModalBody>
            <Input
              label="Название"
              placeholder="Введите название"
              value={editName}
              onChange={(e) => setEditName(e.target.value)}
              isRequired
              autoFocus
            />
            <Input
              label="Email"
              placeholder="Введите email (необязательно)"
              value={editEmail}
              onChange={(e) => setEditEmail(e.target.value)}
              type="email"
            />
          </ModalBody>
          <ModalFooter>
            <Button variant="light" onPress={onEditClose}>
              Отмена
            </Button>
            <Button
              color="default"
              onPress={handleUpdateSender}
              isLoading={isUpdating}
              isDisabled={!editName.trim()}
            >
              Сохранить
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </div>
  );
};

