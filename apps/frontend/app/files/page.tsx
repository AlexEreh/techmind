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

export default function FilesPage() {
    const { user, currentCompany, isLoading: authLoading } = useAuth();
    const router = useRouter();
    const [folders, setFolders] = useState<Folder[]>([]);
    const [documents, setDocuments] = useState<Document[]>([]);
    const [selectedDocument, setSelectedDocument] = useState<Document | null>(null);
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

    const loadFolderTree = async () => {
        if (!currentCompany) return;

        setIsLoading(true);
        try {
            console.log("Loading folders for company:", currentCompany.id);
            const data = await foldersApi.getAllCompanyFolders(currentCompany.id);
            setFolders(data.folders || []);
            console.log("Loaded folders:", data.folders);
            //setDocuments(data.documents || []);
        } catch (error) {
            console.error('Failed to load folders:', error);
        } finally {
            setIsLoading(false);
        }
    };

    const handleDocumentSelect = async (doc: Document) => {
        setSelectedDocument(doc);
        // Load full document details
        try {
            const fullDoc = await documentsApi.getById(doc.id);
            setSelectedDocument(fullDoc);
        } catch (error) {
            console.error('Failed to load document:', error);
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
        <div className="flex h-screen bg-background">
            {/* Left Sidebar - Commands */}
            <Sidebar currentView="files" />

            {/* Folder Tree */}
            <div
                className="border-r border-divider overflow-y-auto"
                style={{ width: `${treeWidth}px` }}
            >
                <FolderTree
                    folders={folders}
                    documents={documents}
                    onDocumentSelect={handleDocumentSelect}
                    onRefresh={loadFolderTree}
                    isLoading={isLoading}
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

            {/* File Preview */}
            <div
                className="flex-1 border-r border-divider overflow-y-auto"
                style={{ width: `${previewWidth}px` }}
            >
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
                <FileInfo document={selectedDocument} onUpdate={loadFolderTree} />
            </div>
        </div>
    );
}

