'use client';

import { Document } from '@/lib/api/types';
import { Card, CardBody } from '@heroui/card';
import { Image } from '@heroui/image';

interface FilePreviewProps {
  document: Document | null;
}

export const FilePreview: React.FC<FilePreviewProps> = ({ document }) => {
  if (!document) {
    return (
      <div className="flex items-center justify-center h-full text-default-400">
        <div className="text-center">
          <p className="text-lg">Выберите файл для просмотра</p>
        </div>
      </div>
    );
  }

  const isImage = document.mime_type.startsWith('image/');
  const isPdf = document.mime_type === 'application/pdf';

  return (
    <div className="p-4 h-full overflow-auto">
      <Card className="h-full">
        <CardBody className="flex items-center justify-center">
          {document.preview_url && (
            <Image
              src={document.preview_url}
              alt={document.name}
              className="max-w-full max-h-full object-contain"
            />
          )}
          {!document.preview_url && isImage && document.download_url && (
            <Image
              src={document.download_url}
              alt={document.name}
              className="max-w-full max-h-full object-contain"
            />
          )}
          {!document.preview_url && isPdf && (
            <div className="text-center">
              <p className="mb-4">Предпросмотр PDF</p>
              {document.download_url && (
                <iframe
                  src={document.download_url}
                  className="w-full h-[600px] border-0"
                  title={document.name}
                />
              )}
            </div>
          )}
          {!document.preview_url && !isImage && !isPdf && (
            <div className="text-center text-default-400">
              <p>Предпросмотр недоступен</p>
              <p className="text-sm mt-2">{document.mime_type}</p>
            </div>
          )}
        </CardBody>
      </Card>
    </div>
  );
};

