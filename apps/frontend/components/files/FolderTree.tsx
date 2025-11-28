'use client';

import { useState } from 'react';
import { Button } from '@heroui/button';
import { Spinner } from '@heroui/spinner';
import { Dropdown, DropdownTrigger, DropdownMenu, DropdownItem } from '@heroui/dropdown';
import { Folder } from '@/lib/api/types';
import { ChevronDownIcon, ChevronRightIcon, FolderIcon, MoreVerticalIcon } from '@/components/icons';
import { useAuth } from '@/contexts/AuthContext';

interface FolderTreeProps {
  folders: Folder[];
  onFolderSelect: (folderId: string | null) => void;
  selectedFolderId: string | null;
  isLoading: boolean;
  highlightedIds?: string[];
  showRoot?: boolean;
}

export const FolderTree: React.FC<FolderTreeProps> = ({
  folders,
  onFolderSelect,
  selectedFolderId,
  isLoading,
  highlightedIds = [],
  showRoot = false,
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
    const elements: JSX.Element[] = [];

    childFolders.forEach((folder) => {
      const isExpanded = expandedFolders.has(folder.id);
      const isHighlighted = highlightedIds.includes(folder.id);
      elements.push(
        <div
          key={folder.id}
          style={{ paddingLeft: level * 16 }}
          className={`flex items-center py-1 ${selectedFolderId === folder.id ? 'bg-blue-100' : ''}`}
        >
          <button
            className="flex items-center w-full text-left"
            onClick={() => {
              toggleFolder(folder.id);
              onFolderSelect(folder.id);
            }}
          >
            {isExpanded ? <ChevronDownIcon /> : <ChevronRightIcon />}
            <FolderIcon className="ml-1 mr-2" />
            <span className={isHighlighted ? 'font-bold' : ''}>{folder.name}</span>
          </button>
        </div>
      );
      if (isExpanded) {
        elements.push(...buildTree(folder.id, level + 1));
      }
    });
    return elements;
  };

  return (
    <div className="relative">
      {showRoot && (
        <div className={`flex items-center py-1 hover:bg-default-100 ${selectedFolderId === null ? 'bg-primary/20' : ''}`}>
          <button
            className="flex items-center w-full text-left"
            onClick={() => onFolderSelect(null)}
          >
            <span className="w-6" /> {/* Спейсер для выравнивания с папками с иконками шеврона */}
            <FolderIcon className="ml-1 mr-2" />
            <span className="font-bold">Корень</span>
          </button>
        </div>
      )}
      {isLoading ? (
        <div className="flex items-center justify-center h-32">
          <Spinner size="sm" />
        </div>
      ) : (
        <div>{buildTree()}</div>
      )}
    </div>
  );
};
