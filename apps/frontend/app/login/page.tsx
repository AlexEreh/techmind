'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Input } from '@heroui/input';
import { Button } from '@heroui/button';
import { Card, CardBody, CardHeader } from '@heroui/card';
import { Link } from '@heroui/link';
import { useAuth } from '@/contexts/AuthContext';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await login(email, password);
      router.push('/select-company');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Ошибка входа. Проверьте email и пароль.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen">
      <Card className="w-full max-w-md" style={{background: "#242428"}}>
        <CardHeader className="flex flex-col gap-1 items-center pt-6">
          <h1 className="text-2xl font-bold">Вход в систему</h1>
          <p className="text-small text-default-500">Менеджер файлов организации</p>
        </CardHeader>
        <CardBody className="gap-4">
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <Input
              label="Email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="user@example.com"
              isRequired
              variant="bordered"
            />
            <Input
              label="Пароль"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Введите пароль"
              isRequired
              variant="bordered"
            />
            {error && (
              <div className="text-danger text-small">{error}</div>
            )}
            <Button
              type="submit"
              color="default"
              isLoading={isLoading}
              className="w-full"
            >
              Войти
            </Button>
          </form>
          <div className="text-center text-small">
            Нет аккаунта?{' '}
            <Link href="/register" size="sm" style={{color: "#A2A2A2"}}>
              Зарегистрироваться
            </Link>
          </div>
        </CardBody>
      </Card>
    </div>
  );
}

