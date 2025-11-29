'use client';

import { useState } from 'react';
import { Button } from '@heroui/button';
import { Spinner } from '@heroui/spinner';
import { Divider } from '@heroui/divider';
import { Folder, Document } from '@/lib/api/types';
import { ChevronDownIcon, ChevronRightIcon, FolderIcon, FileIcon, UploadIcon } from '@/components/icons';

interface FolderTreeProps {
  folders: Folder[];
  documents?: Document[];
  onFolderSelect: (folderId: string | null) => void;
  onDocumentSelect?: (document: Document) => void;
  onUploadClick?: () => void;
  selectedFolderId: string | null;
  selectedDocumentId?: string | null;
  isLoading: boolean;
  highlightedIds?: string[];
  showRoot?: boolean;
}

export const FolderTree: React.FC<FolderTreeProps> = ({
  folders,
  documents = [],
  onFolderSelect,
  onDocumentSelect,
  onUploadClick,
  selectedFolderId,
  selectedDocumentId,
  isLoading,
  highlightedIds = [],
  showRoot = false,
}) => {
  const [expandedFolders, setExpandedFolders] = useState<Set<string>>(new Set());

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
          className={`flex items-center py-1 hover:bg-default-100 cursor-pointer ${selectedFolderId === folder.id ? 'bg-primary/20' : ''}`}
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
    <div className="flex flex-col h-full">
      <div className="flex-1 overflow-y-auto">
        {/* Folder tree */}
        <div className="mb-2">
          {showRoot && (
            <div className={`flex items-center py-1 px-2 hover:bg-default-100 cursor-pointer ${selectedFolderId === null ? 'bg-primary/20' : ''}`}>
              <button
                className="flex items-center w-full text-left"
                onClick={() => onFolderSelect(null)}
              >
                <span className="w-6" />
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
            <div className="px-2">{buildTree()}</div>
          )}
        </div>

        <Divider className="my-2" />

        {/* Documents list */}
        <div className="px-2">
          <p className="text-xs font-semibold text-default-500 uppercase mb-2 px-2">
            Файлы {selectedFolderId ? 'в папке' : 'в корне'}
          </p>
          {documents.length > 0 ? (
            <div className="space-y-1">
              {documents.map((doc) => (
                <div
                  key={doc.id}
                  className={`flex items-center py-2 px-2 rounded cursor-pointer hover:bg-default-100 ${
                    selectedDocumentId === doc.id ? 'bg-primary/20' : ''
                  }`}
                  onClick={() => onDocumentSelect?.(doc)}
                >
                  <FileIcon className="w-4 h-4 mr-2 flex-shrink-0" />
                  <span className="text-sm truncate">{doc.name}</span>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-sm text-default-400 px-2">
              {isLoading ? 'Загрузка...' : 'Нет файлов'}
            </p>
          )}
        </div>
      </div>

      {onUploadClick && (
        <div className="p-2 border-t border-divider">
          <Button
            fullWidth
            color="primary"
            variant="flat"
            startContent={<UploadIcon />}
            onPress={onUploadClick}
            size="sm"
          >
            Загрузить документ
          </Button>
        </div>
      )}
    </div>
  );
};

