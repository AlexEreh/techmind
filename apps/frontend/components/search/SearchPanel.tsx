'use client';

import { useState, useEffect } from 'react';
import { Input } from '@heroui/input';
import { Button } from '@heroui/button';
import { Select, SelectItem } from '@heroui/select';
import { Chip } from '@heroui/chip';
import { Divider } from '@heroui/divider';
import { Document, Tag, Sender } from '@/lib/api/types';
import { tagsApi } from '@/lib/api/tags';
import { sendersApi } from '@/lib/api/senders';
import { useAuth } from '@/contexts/AuthContext';
import { SearchIcon, FileIcon } from '@/components/icons';

interface SearchPanelProps {
  onSearch: (query: string, tagIds: string[], senderId?: string) => void;
  results: Document[];
  onDocumentSelect: (doc: Document) => void;
  isLoading: boolean;
  highlightedIds: string[];
}

export const SearchPanel: React.FC<SearchPanelProps> = ({
  onSearch,
  results,
  onDocumentSelect,
  isLoading,
  highlightedIds,
}) => {
  const [query, setQuery] = useState('');
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [selectedSender, setSelectedSender] = useState<string>('');
  const [allTags, setAllTags] = useState<Tag[]>([]);
  const [allSenders, setAllSenders] = useState<Sender[]>([]);
  const { currentCompany } = useAuth();

  useEffect(() => {
    if (currentCompany) {
      loadFilters();
    }
  }, [currentCompany]);

  const loadFilters = async () => {
    if (!currentCompany) return;
    try {
      const [tagsData, sendersData] = await Promise.all([
        tagsApi.getByCompany(currentCompany.id),
        sendersApi.getByCompany(currentCompany.id),
      ]);
      setAllTags(tagsData.tags);
      setAllSenders(sendersData.senders);
    } catch (error) {
      console.error('Failed to load filters:', error);
    }
  };

  const handleSearch = () => {
    onSearch(query, selectedTags, selectedSender || undefined);
  };

  const toggleTag = (tagId: string) => {
    setSelectedTags((prev) => {
      const updated = prev.includes(tagId) ? prev.filter((id) => id !== tagId) : [...prev, tagId];
      console.log('Selected tags updated:', updated);
      return updated;
    });
  };

  return (
    <div className="p-4 space-y-4">
      <div>
        <h2 className="text-lg font-semibold mb-4">Поиск файлов</h2>
        <Input
          placeholder="Введите запрос..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
          startContent={<SearchIcon className="w-4 h-4" />}
          variant="bordered"
        />
      </div>

      <div>
        <p className="text-sm text-default-500 mb-2">Теги</p>
        <div className="flex flex-wrap gap-2">
          {allTags.map((tag) => (
            <Chip
              key={tag.id}
              onClick={() => toggleTag(tag.id)}
              variant={selectedTags.includes(tag.id) ? 'solid' : 'bordered'}
              color={selectedTags.includes(tag.id) ? 'default' : 'default'}
              className="cursor-pointer"
            >
              {tag.name}
            </Chip>
          ))}
        </div>
      </div>

      <div>
        <p className="text-sm text-default-500 mb-2">Контрагент</p>
        <Select
          placeholder="Выберите контрагента"
          selectedKeys={selectedSender ? [selectedSender] : []}
          onSelectionChange={(keys) => {
            const selected = Array.from(keys)[0] as string;
            setSelectedSender(selected || '');
          }}
          variant="bordered"
          classNames={{
            value: "text-white",
          }}
          className="text-white"
        >
          {allSenders.map((sender) => (
            <SelectItem key={sender.id} textValue={sender.name}>
              {sender.name} {sender.email && `(${sender.email})`}
            </SelectItem>
          ))}
        </Select>
      </div>

      <Button
        color="default"
        onPress={handleSearch}
        isLoading={isLoading}
        className="w-full"
      >
        Поиск
      </Button>

      <Divider />

      <div>
        <p className="text-sm text-default-500 mb-2">
          Результаты: {results.length}
        </p>
        <div className="space-y-1">
          {results.map((doc) => (
            <div
              key={doc.id}
              className={`flex items-start gap-2 p-2 rounded hover:bg-default-100 cursor-pointer ${
                highlightedIds.includes(doc.id) ? 'bg-warning-50' : ''
              }`}
              onClick={() => onDocumentSelect(doc)}
            >
              <FileIcon className="w-4 h-4 text-primary flex-shrink-0 mt-0.5" />
              <div className="flex-1 min-w-0">
                <div className="text-sm truncate">{doc.name}</div>
                {doc.sender_id && (
                  <div className="text-xs text-default-400 truncate">
                    Отправитель: {doc.sender?.name || 'Загрузка...'}
                  </div>
                )}
              </div>
            </div>
          ))}
          {results.length === 0 && !isLoading && (
            <p className="text-sm text-default-400 text-center py-4">
              Нет результатов
            </p>
          )}
        </div>
      </div>
    </div>
  );
};
