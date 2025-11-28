'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import { Sidebar } from '@/components/files/Sidebar';
import { SearchPanel } from '@/components/search/SearchPanel';
import { FilePreview } from '@/components/files/FilePreview';
import { FileInfo } from '@/components/files/FileInfo';
import { documentsApi } from '@/lib/api/documents';
import { Document } from '@/lib/api/types';
import { Spinner } from '@heroui/spinner';

export default function SearchPage() {
    const { user, currentCompany, isLoading: authLoading } = useAuth();
    const router = useRouter();
    const [searchResults, setSearchResults] = useState<Document[]>([]);
    const [selectedDocument, setSelectedDocument] = useState<Document | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [highlightedIds, setHighlightedIds] = useState<string[]>([]);

    useEffect(() => {
        if (!authLoading && !user) {
            router.push('/login');
        }
    }, [user, authLoading, router]);

    const handleSearch = async (query: string, tagIds: string[], senderId?: string) => {
        if (!currentCompany) return;

        setIsLoading(true);
        try {
            const { documents } = await documentsApi.search({
                company_id: currentCompany.id,
                query: query || undefined,
                tag_ids: tagIds.length > 0 ? tagIds : undefined,
                sender_id: senderId || undefined,
            });
            setSearchResults(documents);
            setHighlightedIds(documents.map(d => d.id));
        } catch (error) {
            console.error('Search failed:', error);
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

    if (authLoading || !user || !currentCompany) {
        return (
            <div className="flex items-center justify-center h-screen">
                <Spinner size="lg" />
            </div>
        );
    }

    return (
        <div className="flex h-screen bg-background">
            <Sidebar currentView="search" />

            <div className="w-80 border-r border-divider overflow-y-auto">
                <SearchPanel
                    onSearch={handleSearch}
                    results={searchResults}
                    onDocumentSelect={handleDocumentSelect}
                    isLoading={isLoading}
                    highlightedIds={highlightedIds}
                />
            </div>

            <div className="flex-1 border-r border-divider overflow-y-auto">
                <FilePreview document={selectedDocument} />
            </div>

            <div className="w-80 overflow-y-auto">
                <FileInfo document={selectedDocument} onUpdate={() => {}} />
            </div>
        </div>
    );
}

