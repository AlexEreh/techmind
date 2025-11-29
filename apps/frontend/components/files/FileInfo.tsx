'use client';

import { useState, useEffect } from 'react';
import { Document, Tag } from '@/lib/api/types';
import { Button } from '@heroui/button';
import { Chip } from '@heroui/chip';
import { Divider } from '@heroui/divider';
import { tagsApi } from '@/lib/api/tags';
import { documentsApi } from '@/lib/api/documents';
import { useAuth } from '@/contexts/AuthContext';
import { TrashIcon } from '@/components/icons';

interface FileInfoProps {
  document: Document | null;
  onUpdate: () => void;
  onDelete?: () => void;
}

export const FileInfo: React.FC<FileInfoProps> = ({ document, onUpdate, onDelete }) => {
  const [allTags, setAllTags] = useState<Tag[]>([]);
  const [isDeleting, setIsDeleting] = useState(false);
  const { currentCompany } = useAuth();

  useEffect(() => {
    if (currentCompany) {
      loadTags();
    }
  }, [currentCompany]);

  const loadTags = async () => {
    if (!currentCompany) return;
    try {
      const { tags } = await tagsApi.getByCompany(currentCompany.id);
      setAllTags(tags);
    } catch (error) {
      console.error('Failed to load tags:', error);
    }
  };

  const handleAddTag = async (tagId: string) => {
    if (!document) return;
    try {
      await tagsApi.addToDocument(document.id, tagId);
      onUpdate();
    } catch (error) {
      console.error('Failed to add tag:', error);
    }
  };

  const handleRemoveTag = async (tagId: string) => {
    if (!document) return;
    try {
      await tagsApi.removeFromDocument(document.id, tagId);
      onUpdate();
    } catch (error) {
      console.error('Failed to remove tag:', error);
    }
  };

  const handleDeleteDocument = async () => {
    if (!document) return;

    const confirmed = confirm(`Вы уверены, что хотите удалить файл "${document.name}"? Это действие нельзя отменить.`);
    if (!confirmed) return;

    setIsDeleting(true);
    try {
      await documentsApi.delete(document.id);
      if (onDelete) {
        onDelete();
      }
    } catch (error) {
      console.error('Failed to delete document:', error);
      alert('Не удалось удалить файл');
    } finally {
      setIsDeleting(false);
    }
  };

  if (!document) {
    return (
      <div className="p-4">
        <p className="text-default-400 text-center">Информация о файле</p>
      </div>
    );
  }

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes} Б`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} КБ`;
    return `${(bytes / (1024 * 1024)).toFixed(2)} МБ`;
  };

  const formatDate = (date: string) => {
    return new Date(date).toLocaleString('ru-RU');
  };

  return (
    <div className="p-4 space-y-4">
      <div>
        <h3 className="text-lg font-semibold mb-2">Информация о файле</h3>
        <Divider />
      </div>

      <div>
        <p className="text-sm text-default-500 mb-1">Имя файла</p>
        <p className="text-sm font-medium break-all">{document.name}</p>
      </div>

      <div>
        <p className="text-sm text-default-500 mb-1">Размер</p>
        <p className="text-sm">{formatFileSize(document.file_size)}</p>
      </div>

      <div>
        <p className="text-sm text-default-500 mb-1">Тип</p>
        <p className="text-sm">{document.mime_type}</p>
      </div>

      <div>
        <p className="text-sm text-default-500 mb-1">Дата создания</p>
        <p className="text-sm">{formatDate(document.created_at)}</p>
      </div>

      <Divider />

      <div>
        <p className="text-sm text-default-500 mb-2">Теги</p>
        <div className="flex flex-wrap gap-2 mb-2">
          {document.tags?.map((tag) => (
            <Chip
              key={tag.id}
              onClose={() => handleRemoveTag(tag.id)}
              variant="flat"
              color="default"
            >
              {tag.name}
            </Chip>
          ))}
        </div>
        <div className="flex flex-wrap gap-2">
          {allTags
            .filter((tag) => !document.tags?.find((t) => t.id === tag.id))
            .map((tag) => (
              <Chip
                key={tag.id}
                onClick={() => handleAddTag(tag.id)}
                variant="bordered"
                className="cursor-pointer hover:bg-default-100"
              >
                + {tag.name}
              </Chip>
            ))}
        </div>
      </div>

      {document.download_url && (
        <>
          <Divider />
          <div className="space-y-2">
            <Button
              as="a"
              href={document.download_url}
              download
              color="default"
              className="w-full"
            >
              Скачать файл
            </Button>
            <Button
              color="danger"
              variant="flat"
              className="w-full"
              startContent={<TrashIcon />}
              onPress={handleDeleteDocument}
              isLoading={isDeleting}
            >
              Удалить файл
            </Button>
          </div>
        </>
      )}
    </div>
  );
};

