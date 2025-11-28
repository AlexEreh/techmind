'use client';

import { useState } from 'react';
import { Button } from '@heroui/button';
import { Spinner } from '@heroui/spinner';
import { Dropdown, DropdownTrigger, DropdownMenu, DropdownItem } from '@heroui/dropdown';
import { Document, Folder } from '@/lib/api/types';
import { ChevronDownIcon, ChevronRightIcon, FolderIcon, FileIcon, MoreVerticalIcon } from '@/components/icons';
import { useAuth } from '@/contexts/AuthContext';

interface FolderTreeProps {
  folders: Folder[];
  documents: Document[];
  onDocumentSelect: (doc: Document) => void;
  onRefresh: () => void;
  isLoading: boolean;
  highlightedIds?: string[];
}

export const FolderTree: React.FC<FolderTreeProps> = ({
  folders,
  documents,
  onDocumentSelect,
  onRefresh,
  isLoading,
  highlightedIds = [],
}) => {
  const [expandedFolders, setExpandedFolders] = useState<Set<string>>(new Set());
  const { currentCompany } = useAuth();

  const toggleFolder = (folderId: string) => {
    setExpandedFolders((prev) => {
      const next = new Set(prev);
      if (next.has(folderId)) {
        next.delete(folderId);
      } else {
        next.add(folderId);
      }
      return next;
    });
  };

  const buildTree = (parentId?: string, level = 0): JSX.Element[] => {
    const childFolders = folders.filter((f) => f.parent_folder_id === parentId);
    const childDocs = documents.filter((d) => d.folder_id === parentId);
    const elements: JSX.Element[] = [];

    childFolders.forEach((folder) => {
      const isExpanded = expandedFolders.has(folder.id);
      const isHighlighted = highlightedIds.includes(folder.id);

      elements.push(
        <div key={folder.id}>
          <div
            className={`flex items-center gap-2 px-2 py-1 hover:bg-default-100 cursor-pointer ${
              isHighlighted ? 'bg-warning-50' : ''
            }`}
            style={{ paddingLeft: `${level * 16 + 8}px` }}
          >
            <Button
              isIconOnly
              size="sm"
              variant="light"
              onPress={() => toggleFolder(folder.id)}
              className="min-w-6 w-6 h-6"
            >
              {isExpanded ? (
                <ChevronDownIcon className="w-4 h-4" />
              ) : (
                <ChevronRightIcon className="w-4 h-4" />
              )}
            </Button>
            <FolderIcon className="w-4 h-4 text-warning" />
            <span className="flex-1 text-sm">{folder.name}</span>
            <Dropdown>
              <DropdownTrigger>
                <Button isIconOnly size="sm" variant="light" className="min-w-6 w-6 h-6">
                  <MoreVerticalIcon className="w-4 h-4" />
                </Button>
              </DropdownTrigger>
              <DropdownMenu aria-label="Folder actions">
                <DropdownItem key="rename">Переименовать</DropdownItem>
                <DropdownItem key="delete" className="text-danger" color="danger">
                  Удалить
                </DropdownItem>
              </DropdownMenu>
            </Dropdown>
          </div>
          {isExpanded && buildTree(folder.id, level + 1)}
        </div>
      );
    });

    childDocs.forEach((doc) => {
      const isHighlighted = highlightedIds.includes(doc.id);

      elements.push(
        <div
          key={doc.id}
          className={`flex items-center gap-2 px-2 py-1 hover:bg-default-100 cursor-pointer ${
            isHighlighted ? 'bg-warning-50' : ''
          }`}
          style={{ paddingLeft: `${level * 16 + 40}px` }}
          onClick={() => onDocumentSelect(doc)}
        >
          <FileIcon className="w-4 h-4 text-primary" />
          <span className="flex-1 text-sm truncate">{doc.name}</span>
        </div>
      );
    });

    return elements;
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center p-8">
        <Spinner />
      </div>
    );
  }

  return (
    <div className="p-2">
      <div className="flex items-center justify-between mb-4 px-2">
        <h2 className="text-lg font-semibold">Файлы</h2>
        <Button size="sm" variant="light" onPress={onRefresh}>
          Обновить
        </Button>
      </div>
      {buildTree()}
      {folders.length === 0 && documents.length === 0 && (
        <div className="text-center text-default-400 py-8">
          Нет файлов и папок
        </div>
      )}
    </div>
  );
};
