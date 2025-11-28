'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function SettingsPage() {
  const router = useRouter();

  useEffect(() => {
    // Редирект на страницу тегов
    router.replace('/tags');
  }, [router]);

  return null;
}

