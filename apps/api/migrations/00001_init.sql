-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ===========================
-- users
-- ===========================
CREATE TABLE users
(
    id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name     TEXT NOT NULL,
    email    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

-- ===========================
-- companies
-- ===========================
CREATE TABLE companies
(
    id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL
);

-- ===========================
-- company_users
-- ===========================
CREATE TABLE company_users
(
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      UUID NOT NULL,
    company_id UUID NOT NULL,
    role         INT  NOT NULL,

    CONSTRAINT fk_company_users_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_company_users_company FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
);

CREATE INDEX idx_company_users_user_id ON company_users (user_id);
CREATE INDEX idx_company_users_company_id ON company_users (company_id);

-- ===========================
-- folders
-- ===========================
CREATE TABLE folders
(
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id       UUID   NOT NULL,
    parent_folder_id UUID             DEFAULT NULL,
    name             TEXT   NOT NULL,
    size             BIGINT NOT NULL  DEFAULT 0, -- размер всех документов в байтах
    count            INT    NOT NULL  DEFAULT 0, -- количество документов

    CONSTRAINT fk_folders_company FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE,
    CONSTRAINT fk_folders_parent FOREIGN KEY (parent_folder_id) REFERENCES folders (id) ON DELETE SET NULL
);

CREATE INDEX idx_folders_company_id ON folders (company_id);
CREATE INDEX idx_folders_parent_folder_id ON folders (parent_folder_id);

-- ===========================
-- senders
-- ===========================
CREATE TABLE senders
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL,
    name       TEXT NOT NULL,
    email      TEXT             DEFAULT NULL,

    CONSTRAINT fk_senders_company FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
);

CREATE INDEX idx_senders_company_id ON senders (company_id);

-- ===========================
-- documents
-- ===========================
CREATE TABLE documents
(
    id         UUID PRIMARY KEY   DEFAULT uuid_generate_v4(),
    company_id UUID      NOT NULL,
    folder_id  UUID               DEFAULT NULL,
    name       TEXT      NOT NULL,
    file_path  TEXT      NOT NULL,
    preview_file_path TEXT      DEFAULT NULL,
    file_size  BIGINT    NOT NULL,
    mime_type  TEXT      NOT NULL,
    checksum   TEXT      NOT NULL,
    sender_id  UUID               DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_documents_company FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE,
    CONSTRAINT fk_documents_folder FOREIGN KEY (folder_id) REFERENCES folders (id) ON DELETE SET NULL,
    CONSTRAINT fk_documents_sender FOREIGN KEY (sender_id) REFERENCES senders (id) ON DELETE SET NULL
);

CREATE INDEX idx_documents_company_id ON documents (company_id);
CREATE INDEX idx_documents_folder_id ON documents (folder_id);
CREATE INDEX idx_documents_sender_id ON documents (sender_id);

-- ===========================
-- tags
-- ===========================
CREATE TABLE tags
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL,
    name       TEXT NOT NULL,

    CONSTRAINT fk_tags_company FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE,
    CONSTRAINT uq_tags_company_name UNIQUE (company_id, name)
);

CREATE INDEX idx_tags_company_id ON tags (company_id);

-- ===========================
-- document_tags
-- ===========================
CREATE TABLE document_tags
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id UUID NOT NULL,
    tag_id      UUID NOT NULL,

    CONSTRAINT fk_document_tags_document FOREIGN KEY (document_id) REFERENCES documents (id) ON DELETE CASCADE,
    CONSTRAINT fk_document_tags_tag FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
);

CREATE INDEX idx_document_tags_document_id ON document_tags (document_id);
CREATE INDEX idx_document_tags_tag_id ON document_tags (tag_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
