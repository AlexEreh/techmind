# Настройка Gotenberg Connection

## Описание

Gotenberg клиент теперь интегрирован в приложение через Uber FX dependency injection и автоматически инициализируется при старте приложения.

## Конфигурация

### dev.yml / prod.yml

```yaml
gotenberg:
  url: "http://localhost:3000"  # URL Gotenberg сервиса
  enabled: true                  # Включить/выключить Gotenberg
  timeout: 60                    # Таймаут в секундах (по умолчанию 60)
```

### Параметры

- **url** - URL Gotenberg API (обязательный)
- **enabled** - Флаг включения/выключения сервиса. Если `false`, клиент не будет создан
- **timeout** - Таймаут HTTP запросов в секундах (по умолчанию 60)

## Docker Compose

Gotenberg добавлен в `docker-compose.dev.yml`:

```yaml
gotenberg:
  image: gotenberg/gotenberg:8
  ports:
    - "3000:3000"
  command:
    - "gotenberg"
    - "--api-timeout=60s"
    - "--log-level=info"
```

### Запуск

```bash
# Запустить все сервисы включая Gotenberg
docker-compose -f docker-compose.dev.yml up -d

# Проверить статус Gotenberg
curl http://localhost:3000/health
```

## Использование в коде

Gotenberg клиент автоматически инжектируется в сервисы через FX:

```go
// В document service
func NewService(
	documentRepo repo.DocumentRepository,
	documentTagRepo repo.DocumentTagRepository,
	tagRepo repo.TagRepository,
	folderRepo repo.FolderRepository,
	minioClient *minio.Client,
	gotenbergClient *gotenberg.Client, // Автоматически инжектируется
) service.DocumentService {
	// ...
}
```

## Отключение Gotenberg

Если Gotenberg не требуется, можно отключить его:

```yaml
gotenberg:
  url: "http://localhost:3000"
  enabled: false  # Отключить
  timeout: 60
```

При `enabled: false`:
- Gotenberg клиент не будет создан (`nil`)
- Функции, требующие Gotenberg, будут возвращать ошибку
- Приложение продолжит работать без конвертации документов

## Архитектура

```
app/app.go
  ├── fx.Provide(gotenberg.New)  // Регистрация провайдера
  └── connections/gotenberg/
      └── gotenberg.go            // Инициализация клиента
          ├── Читает config
          ├── Создает HTTP client с таймаутом
          ├── Создает Gotenberg client
          └── Регистрирует lifecycle hooks
```

## Проверка работоспособности

### 1. Проверка конфигурации

```bash
# Проверить что Gotenberg запущен
curl http://localhost:3000/health

# Ожидаемый ответ:
{
  "status": "up",
  "details": {
    "chromium": {"status": "up"},
    "libreoffice": {"status": "up"}
  }
}
```

### 2. Проверка версии

```bash
curl http://localhost:3000/version
# Ожидаемый ответ: 8.x.x
```

### 3. Тест конвертации

```bash
# Создать тестовый HTML файл
echo "<h1>Test</h1>" > test.html

# Конвертировать в PDF
curl --request POST \
  --url http://localhost:3000/forms/chromium/convert/html \
  --form files=@test.html \
  --output test.pdf
```

## Troubleshooting

### Gotenberg не запускается

```bash
# Проверить логи
docker-compose -f docker-compose.dev.yml logs gotenberg

# Перезапустить контейнер
docker-compose -f docker-compose.dev.yml restart gotenberg
```

### Таймауты при конвертации

Увеличьте таймаут в конфигурации:

```yaml
gotenberg:
  url: "http://localhost:3000"
  enabled: true
  timeout: 120  # Увеличить до 120 секунд
```

### Ошибка "connection refused"

Проверьте что:
1. Gotenberg запущен: `docker ps | grep gotenberg`
2. URL правильный в конфиге
3. Порты доступны: `netstat -an | grep 3000`

## Production настройка

Для production окружения в `prod.yml`:

```yaml
gotenberg:
  url: "http://gotenberg:3000"  # Внутренний DNS в Docker
  enabled: true
  timeout: 90
```

В docker-compose.yml:

```yaml
gotenberg:
  image: gotenberg/gotenberg:8
  restart: always
  networks:
    - backend
  command:
    - "gotenberg"
    - "--api-timeout=90s"
    - "--log-level=error"
  deploy:
    resources:
      limits:
        memory: 2G
      reservations:
        memory: 1G
```

## См. также

- [Gotenberg Documentation](https://gotenberg.dev/)
- [Preview Generation Guide](./PREVIEW_GENERATION.md)
- [API Documentation](./API_GUIDE.md)

