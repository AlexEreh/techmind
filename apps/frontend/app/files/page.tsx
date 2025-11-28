'use client';

import { useState, useEffect, useRef } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import { Sidebar } from '@/components/files/Sidebar';
import { FolderTree } from '@/components/files/FolderTree';
import { FilePreview } from '@/components/files/FilePreview';
import { FileInfo } from '@/components/files/FileInfo';
import { foldersApi } from '@/lib/api/folders';
import { documentsApi } from '@/lib/api/documents';
import { Document, Folder } from '@/lib/api/types';
import { Spinner } from '@heroui/spinner';

export default function FilesPage() {
    const { user, currentCompany, isLoading: authLoading } = useAuth();
    const router = useRouter();
    const [folders, setFolders] = useState<Folder[]>([]);
    const [selectedDocument, setSelectedDocument] = useState<Document | null>(null);
    const [selectedFolderId, setSelectedFolderId] = useState<string | null>(null);
    const [folderDocuments, setFolderDocuments] = useState<Document[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [treeWidth, setTreeWidth] = useState(300);
    const [isDragging, setIsDragging] = useState(false);
    const fileInputRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
        if (!authLoading && !user) {
            router.push('/login');
        }
    }, [user, authLoading, router]);

    useEffect(() => {
        if (currentCompany) {
            loadFolderTree();
        }
    }, [currentCompany]);

    // Автоматически выбираем корень при первой загрузке
    useEffect(() => {
        if (currentCompany && selectedFolderId === null) {
            handleFolderSelect(null);
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [currentCompany]);

    const loadFolderTree = async () => {
        if (!currentCompany) return;
        setIsLoading(true);
        try {
            const data = await foldersApi.getAllCompanyFolders(currentCompany.id);
            setFolders(data.folders || []);
        } catch (error) {
            console.error('Failed to load folders:', error);
        } finally {
            setIsLoading(false);
        }
    };

    // Лениво подгружать документы только для выбранной папки
    const handleFolderSelect = async (folderId: string | null) => {
        setSelectedFolderId(folderId);
        setSelectedDocument(null);
        setIsLoading(true);
        try {
            if (folderId === null) {
                // Корень: получить все документы компании и отфильтровать по folder_id == null/undefined/''
                if (!currentCompany) throw new Error('Компания не выбрана');
                const data = await documentsApi.getByCompany(currentCompany.id);
                setFolderDocuments((data.documents || []).filter(doc => doc.folder_id === null || doc.folder_id === undefined || doc.folder_id === ''));
            } else {
                const data = await documentsApi.getByFolder(folderId);
                setFolderDocuments(data.documents || []);
            }
        } catch (error) {
            setFolderDocuments([]);
            console.error('Failed to load documents for folder:', error);
        } finally {
            setIsLoading(false);
        }
    };

    const handleDocumentSelect = async (doc: Document) => {
        setSelectedDocument(doc);
        try {
            const fullDoc = await documentsApi.getById(doc.id);
            setSelectedDocument(fullDoc);
        } catch (error) {
            console.error('Failed to load document:', error);
        }
    };

    const handleDocumentUpdate = async () => {
        // Перезагружаем текущий документ после изменений (например, добавления/удаления тегов)
        if (selectedDocument) {
            try {
                const updatedDoc = await documentsApi.getById(selectedDocument.id);
                setSelectedDocument(updatedDoc);
            } catch (error) {
                console.error('Failed to reload document:', error);
            }
        }
    };

    const handleFilesUpload = async (files: FileList) => {
        if (!currentCompany) return;
        try {
            for (const file of Array.from(files)) {
                await documentsApi.upload({
                    company_id: currentCompany.id,
                    name: file.name,
                    file,
                    folder_id: selectedFolderId || undefined,
                });
            }
            // После загрузки обновить содержимое текущей папки
            handleFolderSelect(selectedFolderId);
        } catch (e) {
            alert('Ошибка загрузки файлов');
        }
    };

    const handleUploadClick = () => {
        fileInputRef.current?.click();
    };

    const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files.length > 0) {
            handleFilesUpload(e.target.files);
            e.target.value = ''; // Сбрасываем input для повторной загрузки
        }
    };

    // Drag & Drop handlers
    const handleDragEnter = (e: React.DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setIsDragging(true);
    };

    const handleDragLeave = (e: React.DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
        if (e.currentTarget === e.target) {
            setIsDragging(false);
        }
    };

    const handleDragOver = (e: React.DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
    };

    const handleDrop = (e: React.DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setIsDragging(false);

        const files = e.dataTransfer.files;
        if (files && files.length > 0) {
            handleFilesUpload(files);
        }
    };

    if (authLoading || !user || !currentCompany) {
        return (
            <div className="flex items-center justify-center h-screen">
                <Spinner size="lg" />
            </div>
        );
    }

    return (
        <div
            className="flex h-screen bg-background relative"
            onDragEnter={handleDragEnter}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
        >
            {/* Hidden file input */}
            <input
                ref={fileInputRef}
                type="file"
                multiple
                className="hidden"
                onChange={handleFileInputChange}
            />

            {/* Drag & Drop overlay */}
            {isDragging && (
                <div className="absolute inset-0 z-50 bg-primary/10 border-4 border-dashed border-primary flex items-center justify-center">
                    <div className="text-center">
                        <p className="text-2xl font-bold text-primary">Перетащите файлы сюда</p>
                        <p className="text-default-500">
                            {selectedFolderId
                                ? `Загрузка в: ${folders.find(f => f.id === selectedFolderId)?.name || 'выбранную папку'}`
                                : 'Загрузка в корень'
                            }
                        </p>
                    </div>
                </div>
            )}

            {/* Left Sidebar - Commands */}
            <Sidebar currentView="files" onUploadClick={handleUploadClick} />

            {/* Folder Tree with Documents */}
            <div
                className="border-r border-divider overflow-hidden relative"
                style={{ width: `${treeWidth}px` }}
            >
                <FolderTree
                    folders={folders}
                    documents={folderDocuments}
                    onFolderSelect={handleFolderSelect}
                    onDocumentSelect={handleDocumentSelect}
                    selectedFolderId={selectedFolderId}
                    selectedDocumentId={selectedDocument?.id}
                    isLoading={isLoading}
                    showRoot
                    onFolderCreated={loadFolderTree}
                />
                <div
                    className="absolute top-0 right-0 w-1 h-full cursor-col-resize hover:bg-primary z-10"
                    onMouseDown={(e) => {
                        const startX = e.clientX;
                        const startWidth = treeWidth;
                        const handleMouseMove = (e: MouseEvent) => {
                            const newWidth = startWidth + (e.clientX - startX);
                            setTreeWidth(Math.max(250, Math.min(600, newWidth)));
                        };
                        const handleMouseUp = () => {
                            document.removeEventListener('mousemove', handleMouseMove);
                            document.removeEventListener('mouseup', handleMouseUp);
                        };
                        document.addEventListener('mousemove', handleMouseMove);
                        document.addEventListener('mouseup', handleMouseUp);
                    }}
                />
            </div>

            {/* File Preview - Center */}
            <div className="flex-1 border-r border-divider overflow-hidden">
                <FilePreview document={selectedDocument} />
            </div>

            {/* File Info */}
            <div className="w-80 overflow-y-auto">
                <FileInfo document={selectedDocument} onUpdate={handleDocumentUpdate} />
            </div>
        </div>
    );
}
