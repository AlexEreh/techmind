-- +goose Up
-- +goose StatementBegin

-- Добавляем поля created_by и updated_by в таблицу documents (если они еще не существуют)
ALTER TABLE documents
    ADD COLUMN IF NOT EXISTS created_by UUID DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS updated_by UUID DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();

-- Добавляем внешние ключи (если они еще не существуют)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_documents_created_by'
    ) THEN
        ALTER TABLE documents
            ADD CONSTRAINT fk_documents_created_by FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_documents_updated_by'
    ) THEN
        ALTER TABLE documents
            ADD CONSTRAINT fk_documents_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL;
    END IF;
END $$;

-- Добавляем индексы для быстрого поиска (если они еще не существуют)
CREATE INDEX IF NOT EXISTS idx_documents_created_by ON documents (created_by);
CREATE INDEX IF NOT EXISTS idx_documents_updated_by ON documents (updated_by);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Удаляем индексы
DROP INDEX IF EXISTS idx_documents_created_by;
DROP INDEX IF EXISTS idx_documents_updated_by;

-- Удаляем внешние ключи
ALTER TABLE documents
    DROP CONSTRAINT IF EXISTS fk_documents_created_by,
    DROP CONSTRAINT IF EXISTS fk_documents_updated_by;

-- Удаляем колонки
ALTER TABLE documents
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS updated_by,
    DROP COLUMN IF EXISTS updated_at;

-- +goose StatementEnd

