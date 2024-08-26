-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY,
                                     email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL
    );

-- Создание таблицы для хранения Refresh токенов
CREATE TABLE IF NOT EXISTS refresh_tokens (
                                              id SERIAL PRIMARY KEY,
                                              token_hash VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user
    FOREIGN KEY(user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
    );

-- Индекс для ускорения поиска по user_id
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);
