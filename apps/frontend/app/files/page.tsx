'use client';

import { useState, useEffect } from 'react';
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
import { Dropzone } from '@/components/files/Dropzone';

export default function FilesPage() {
    const { user, currentCompany, isLoading: authLoading } = useAuth();
    const router = useRouter();
    const [folders, setFolders] = useState<Folder[]>([]);
    const [selectedDocument, setSelectedDocument] = useState<Document | null>(null);
    const [selectedFolderId, setSelectedFolderId] = useState<string | null>(null);
    const [folderDocuments, setFolderDocuments] = useState<Document[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [treeWidth, setTreeWidth] = useState(300);
    const [previewWidth, setPreviewWidth] = useState(600);

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
            if (selectedFolderId) {
                handleFolderSelect(selectedFolderId);
            }
        } catch (e) {
            alert('Ошибка загрузки файлов');
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
        <div className="flex flex-col h-screen bg-background">
            <div className="p-4">
                <Dropzone onFilesAdded={handleFilesUpload} />
            </div>
            <div className="flex flex-1 min-h-0">
                {/* Left Sidebar - Commands */}
                <Sidebar currentView="files" />
                {/* Folder Tree */}
                <div
                    className="border-r border-divider overflow-y-auto"
                    style={{ width: `${treeWidth}px` }}
                >
                    <FolderTree
                        folders={folders}
                        onFolderSelect={handleFolderSelect}
                        selectedFolderId={selectedFolderId}
                        isLoading={isLoading}
                        showRoot
                    />
                    <div
                        className="absolute top-0 right-0 w-1 h-full cursor-col-resize hover:bg-primary"
                        onMouseDown={(e) => {
                            const startX = e.clientX;
                            const startWidth = treeWidth;
                            const handleMouseMove = (e: MouseEvent) => {
                                const newWidth = startWidth + (e.clientX - startX);
                                setTreeWidth(Math.max(200, Math.min(600, newWidth)));
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
                {/* File Preview + список файлов выбранной папки */}
                <div
                    className="flex-1 border-r border-divider overflow-y-auto"
                    style={{ width: `${previewWidth}px` }}
                >
                    {/* Список файлов выбранной папки */}
                    <div className="p-4">
                        {selectedFolderId === null && folderDocuments.length > 0 && (
                            <ul>
                                {folderDocuments.map(doc => (
                                    <li key={doc.id} className="cursor-pointer hover:underline" onClick={() => handleDocumentSelect(doc)}>
                                        {doc.name}
                                    </li>
                                ))}
                            </ul>
                        )}
                        {selectedFolderId === null && folderDocuments.length === 0 && !isLoading && (
                            <div className="text-gray-400">В корне нет файлов</div>
                        )}
                        {selectedFolderId && folderDocuments.length > 0 && (
                            <ul>
                                {folderDocuments.map(doc => (
                                    <li key={doc.id} className="cursor-pointer hover:underline" onClick={() => handleDocumentSelect(doc)}>
                                        {doc.name}
                                    </li>
                                ))}
                            </ul>
                        )}
                        {selectedFolderId && folderDocuments.length === 0 && !isLoading && (
                            <div className="text-gray-400">В папке нет файлов</div>
                        )}
                    </div>
                    <FilePreview document={selectedDocument} />
                    <div
                        className="absolute top-0 right-0 w-1 h-full cursor-col-resize hover:bg-primary"
                        onMouseDown={(e) => {
                            const startX = e.clientX;
                            const startWidth = previewWidth;
                            const handleMouseMove = (e: MouseEvent) => {
                                const newWidth = startWidth + (e.clientX - startX);
                                setPreviewWidth(Math.max(400, Math.min(1000, newWidth)));
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
                {/* File Info */}
                <div className="w-80 overflow-y-auto">
                    <FileInfo document={selectedDocument} onUpdate={() => { if (selectedFolderId) handleFolderSelect(selectedFolderId); }} />
                </div>
            </div>
        </div>
    );
}
