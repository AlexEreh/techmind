-- +goose Up
-- +goose StatementBegin
INSERT INTO companies (name) VALUES ('Тестовая компания');

INSERT INTO users (name, email, password) VALUES ('Пользователь тестовый', 'sus@sus.sus', '$2a$10$3H.kakqB78tQIEde8dGtSennNoiDdf9AAV2uACO3b9Wevjz1BHj22');

INSERT INTO company_users (user_id, company_id, role)
VALUES (
    (SELECT id FROM users WHERE email = 'sus@sus.sus'),
    (SELECT id FROM companies WHERE name = 'Тестовая компания'),
    1
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd