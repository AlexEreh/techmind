'use client';

import { Button } from '@heroui/button';
import { Tooltip } from '@heroui/tooltip';
import { useRouter } from 'next/navigation';
import {
  FileIcon,
  SearchIcon,
  SettingsIcon,
} from '@/components/icons';

interface SidebarProps {
  currentView: 'files' | 'search' | 'settings';
}

export const Sidebar: React.FC<SidebarProps> = ({ currentView }) => {
  const router = useRouter();

  return (
    <div className="w-16 bg-content1 border-r border-divider flex flex-col items-center py-4 gap-4">
      <Tooltip content="Файлы" placement="right">
        <Button
          isIconOnly
          variant={currentView === 'files' ? 'solid' : 'light'}
          color={currentView === 'files' ? 'primary' : 'default'}
          onPress={() => router.push('/files')}
        >
          <FileIcon className="w-5 h-5" />
        </Button>
      </Tooltip>

      <Tooltip content="Поиск" placement="right">
        <Button
          isIconOnly
          variant={currentView === 'search' ? 'solid' : 'light'}
          color={currentView === 'search' ? 'primary' : 'default'}
          onPress={() => router.push('/search')}
        >
          <SearchIcon className="w-5 h-5" />
        </Button>
      </Tooltip>

      <Tooltip content="Настройки" placement="right">
        <Button
          isIconOnly
          variant={currentView === 'settings' ? 'solid' : 'light'}
          color={currentView === 'settings' ? 'primary' : 'default'}
          onPress={() => router.push('/settings')}
        >
          <SettingsIcon className="w-5 h-5" />
        </Button>
      </Tooltip>
    </div>
  );
};
