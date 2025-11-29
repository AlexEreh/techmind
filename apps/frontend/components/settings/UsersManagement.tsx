'use client';

import { useState, useEffect } from 'react';
import { Button } from '@heroui/button';
import { Input } from '@heroui/input';
import { Select, SelectItem } from '@heroui/select';
import { Table, TableHeader, TableColumn, TableBody, TableRow, TableCell } from '@heroui/table';
import { Modal, ModalContent, ModalHeader, ModalBody, ModalFooter, useDisclosure } from '@heroui/modal';
import { CompanyUserWithDetails } from '@/lib/api/types';
import { companyApi } from '@/lib/api/company';
import { useAuth } from '@/contexts/AuthContext';
import { PlusIcon, TrashIcon } from '@/components/icons';

const ROLES = [
  { value: 0, label: 'Просмотр' },
  { value: 1, label: 'Редактор' },
  { value: 2, label: 'Администратор' },
];

export const UsersManagement: React.FC = () => {
  const [users, setUsers] = useState<CompanyUserWithDetails[]>([]);
  const [inviteEmail, setInviteEmail] = useState('');
  const [inviteRole, setInviteRole] = useState('1');
  const [isInviting, setIsInviting] = useState(false);
  const { currentCompany } = useAuth();
  const { isOpen, onOpen, onClose } = useDisclosure();

  useEffect(() => {
    loadUsers();
  }, [currentCompany]);

  const loadUsers = async () => {
    if (!currentCompany) return;
    try {
      const { users: loadedUsers } = await companyApi.getCompanyUsers(currentCompany.id);
      setUsers(loadedUsers);
    } catch (error) {
      console.error('Failed to load users:', error);
    }
  };

  const handleInviteUser = async () => {
    if (!currentCompany || !inviteEmail.trim()) return;

    setIsInviting(true);
    try {
      await companyApi.inviteUser(currentCompany.id, inviteEmail, parseInt(inviteRole));
      setInviteEmail('');
      onClose();
      await loadUsers();
    } catch (error) {
      console.error('Failed to invite user:', error);
    } finally {
      setIsInviting(false);
    }
  };

  const handleUpdateRole = async (userId: string, newRole: number) => {
    if (!currentCompany) return;
    try {
      await companyApi.updateUserRole(userId, newRole);
      await loadUsers();
    } catch (error) {
      console.error('Failed to update role:', error);
    }
  };

  const handleRemoveUser = async (userId: string) => {
    if (!currentCompany || !confirm('Удалить пользователя из компании?')) return;

    try {
      await companyApi.removeUser(userId);
      await loadUsers();
    } catch (error) {
      console.error('Failed to remove user:', error);
    }
  };

  const getRoleLabel = (role: number) => {
    return ROLES.find((r) => r.value === role)?.label || 'Неизвестно';
  };

  return (
    <div className="py-4 space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-semibold">Управление пользователями</h3>
        <Button color="default" startContent={<PlusIcon />} onPress={onOpen}>
          Пригласить пользователя
        </Button>
      </div>

      <Table aria-label="Users table">
        <TableHeader>
          <TableColumn>EMAIL</TableColumn>
          <TableColumn>ИМЯ</TableColumn>
          <TableColumn>РОЛЬ</TableColumn>
          <TableColumn>ДЕЙСТВИЯ</TableColumn>
        </TableHeader>
        <TableBody emptyContent="Нет пользователей">
          {users.map((user) => (
            <TableRow key={user.company_user_id}>
              <TableCell>{user.email || '-'}</TableCell>
              <TableCell>{user.name || '-'}</TableCell>
              <TableCell>
                <Select
                  size="sm"
                  selectedKeys={[user.role.toString()]}
                  onChange={(e) => handleUpdateRole(user.company_user_id, parseInt(e.target.value))}
                  className="w-40"
                >
                  {ROLES.map((role) => (
                    <SelectItem key={role.value}>
                      {role.label}
                    </SelectItem>
                  ))}
                </Select>
              </TableCell>
              <TableCell>
                <Button
                  isIconOnly
                  size="sm"
                  color="danger"
                  variant="light"
                  onPress={() => handleRemoveUser(user.company_user_id)}
                >
                  <TrashIcon />
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>

      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalContent>
          <ModalHeader>Пригласить пользователя</ModalHeader>
          <ModalBody>
            <Input
              label="Email пользователя"
              type="email"
              placeholder="user@example.com"
              value={inviteEmail}
              onChange={(e) => setInviteEmail(e.target.value)}
              autoFocus
            />
            <Select
              label="Роль"
              selectedKeys={[inviteRole]}
              onChange={(e) => setInviteRole(e.target.value)}
            >
              {ROLES.map((role) => (
                <SelectItem key={role.value}>
                  {role.label}
                </SelectItem>
              ))}
            </Select>
          </ModalBody>
          <ModalFooter>
            <Button variant="light" onPress={onClose}>
              Отмена
            </Button>
            <Button
              color="default"
              onPress={handleInviteUser}
              isLoading={isInviting}
              isDisabled={!inviteEmail.trim()}
            >
              Пригласить
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </div>
  );
};
