import React, { useCallback, useRef, useState } from 'react';

interface DropzoneProps {
  onFilesAdded: (files: FileList) => void;
}

export const Dropzone: React.FC<DropzoneProps> = ({ onFilesAdded }) => {
  const [isDragActive, setIsDragActive] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const openFileDialog = () => {
    inputRef.current?.click();
  };

  const onFiles = (files: FileList | null) => {
    if (files && files.length > 0) {
      onFilesAdded(files);
    }
  };

  const handleDrop = useCallback((e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragActive(false);
    onFiles(e.dataTransfer.files);
  }, [onFilesAdded]);

  return (
    <div
      onDragOver={e => {
        e.preventDefault();
        setIsDragActive(true);
      }}
      onDragLeave={e => {
        e.preventDefault();
        setIsDragActive(false);
      }}
      onDrop={handleDrop}
      onClick={openFileDialog}
      className={`border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors ${isDragActive ? 'border-blue-500 bg-blue-50' : 'border-gray-300'}`}
      style={{ minHeight: 120 }}
    >
      <input
        ref={inputRef}
        type="file"
        multiple
        style={{ display: 'none' }}
        onChange={e => onFiles(e.target.files)}
      />
      <p className="text-gray-500">Перетащите файлы сюда или кликните для выбора</p>
    </div>
  );
};

