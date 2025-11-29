'use client';

import { useState, useEffect } from 'react';
import { Button } from '@heroui/button';
import { Input } from '@heroui/input';
import { Chip } from '@heroui/chip';
import { Modal, ModalContent, ModalHeader, ModalBody, ModalFooter, useDisclosure } from '@heroui/modal';
import { Tag } from '@/lib/api/types';
import { tagsApi } from '@/lib/api/tags';
import { useAuth } from '@/contexts/AuthContext';
import { PlusIcon, TrashIcon } from '@/components/icons';

export const TagsManagement: React.FC = () => {
  const [tags, setTags] = useState<Tag[]>([]);
  const [newTagName, setNewTagName] = useState('');
  const [isCreating, setIsCreating] = useState(false);
  const { currentCompany } = useAuth();
  const { isOpen, onOpen, onClose } = useDisclosure();

  useEffect(() => {
    loadTags();
  }, [currentCompany]);

  const loadTags = async () => {
    if (!currentCompany) return;
    try {
      const { tags: loadedTags } = await tagsApi.getByCompany(currentCompany.id);
      setTags(loadedTags);
    } catch (error) {
      console.error('Failed to load tags:', error);
    }
  };

  const handleCreateTag = async () => {
    if (!currentCompany || !newTagName.trim()) return;

    setIsCreating(true);
    try {
      await tagsApi.create({
        company_id: currentCompany.id,
        name: newTagName.trim(),
      });
      setNewTagName('');
      onClose();
      await loadTags();
    } catch (error) {
      console.error('Failed to create tag:', error);
    } finally {
      setIsCreating(false);
    }
  };

  const handleDeleteTag = async (tagId: string) => {
    if (!confirm('Вы уверены, что хотите удалить этот тег?')) return;

    try {
      await tagsApi.delete(tagId);
      await loadTags();
    } catch (error) {
      console.error('Failed to delete tag:', error);
    }
  };

  return (
    <div className="py-4 space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-semibold">Управление тегами</h3>
        <Button color="default" startContent={<PlusIcon />} onPress={onOpen}>
          Создать тег
        </Button>
      </div>

      <div className="flex flex-wrap gap-3">
        {tags.map((tag) => (
          <Chip
            key={tag.id}
            variant="flat"
            color="default"
            onClose={() => handleDeleteTag(tag.id)}
            size="lg"
          >
            {tag.name}
          </Chip>
        ))}
        {tags.length === 0 && (
          <p className="text-default-400">Теги не созданы</p>
        )}
      </div>

      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalContent>
          <ModalHeader>Создать новый тег</ModalHeader>
          <ModalBody>
            <Input
              label="Название тега"
              placeholder="Введите название"
              value={newTagName}
              onChange={(e) => setNewTagName(e.target.value)}
              autoFocus
            />
          </ModalBody>
          <ModalFooter>
            <Button variant="light" onPress={onClose}>
              Отмена
            </Button>
            <Button
              color="default"
              onPress={handleCreateTag}
              isLoading={isCreating}
              isDisabled={!newTagName.trim()}
            >
              Создать
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </div>
  );
};

