'use client';

import { Button } from '@heroui/button';
import { Tooltip } from '@heroui/tooltip';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import {
  FileIcon,
  SearchIcon,
  UploadIcon,
  TagIcon,
  UsersIcon,
  UserIcon,
  FolderIcon,
} from '@/components/icons';

interface SidebarProps {
  currentView: 'files' | 'search' | 'tags' | 'users' | 'profile';
  onUploadClick?: () => void;
}

export const Sidebar: React.FC<SidebarProps> = ({ currentView, onUploadClick }) => {
  const router = useRouter();
  const { currentCompany } = useAuth();

  return (
    <div className="w-16 bg-content1 border-r border-divider flex flex-col items-center py-4 gap-4">
      <Tooltip content="Файлы" placement="right" color={"default"}>
        <Button
          isIconOnly
          variant={currentView === 'files' ? 'solid' : 'light'}
          color={currentView === 'files' ? 'default' : 'default'}
          onPress={() => router.push('/files')}
        >
          <FileIcon className="w-5 h-5" />
        </Button>
      </Tooltip>

      <Tooltip content="Поиск" placement="right">
        <Button
          isIconOnly
          variant={currentView === 'search' ? 'solid' : 'light'}
          color={currentView === 'search' ? 'default' : 'default'}
          onPress={() => router.push('/search')}
        >
          <SearchIcon className="w-5 h-5" />
        </Button>
      </Tooltip>

      <Tooltip content="Теги" placement="right">
        <Button
          isIconOnly
          variant={currentView === 'tags' ? 'solid' : 'light'}
          color={currentView === 'tags' ? 'default' : 'default'}
          onPress={() => router.push('/tags')}
        >
          <TagIcon className="w-5 h-5" />
        </Button>
      </Tooltip>

      <Tooltip content="Пользователи" placement="right">
        <Button
          isIconOnly
          variant={currentView === 'users' ? 'solid' : 'light'}
          color={currentView === 'users' ? 'default' : 'default'}
          onPress={() => router.push('/users')}
        >
          <UsersIcon className="w-5 h-5" />
        </Button>
      </Tooltip>

      <div className="flex-1" />

      {currentCompany && (
        <Tooltip
          content={`Компания: ${currentCompany.name}\nНажмите для смены`}
          placement="right"
        >
          <Button
            isIconOnly
            variant="light"
            onPress={() => router.push('/select-company')}
          >
            <FolderIcon className="w-5 h-5" />
          </Button>
        </Tooltip>
      )}

      <Tooltip content="Профиль" placement="right">
        <Button
          isIconOnly
          variant={currentView === 'profile' ? 'solid' : 'light'}
          color={currentView === 'profile' ? 'default' : 'default'}
          onPress={() => router.push('/profile')}
        >
          <UserIcon className="w-5 h-5" />
        </Button>
      </Tooltip>

      {/*{currentView === 'files' && onUploadClick && (*/}
      {/*  <Tooltip content="Загрузить документ" placement="right">*/}
      {/*    <Button*/}
      {/*      isIconOnly*/}
      {/*      color="success"*/}
      {/*      variant="flat"*/}
      {/*      onPress={onUploadClick}*/}
      {/*    >*/}
      {/*      <UploadIcon className="w-5 h-5" />*/}
      {/*    </Button>*/}
      {/*  </Tooltip>*/}
      {/*)}*/}
    </div>
  );
};
