CREATE TABLE users (
    id uuid NOT NULL PRIMARY KEY,             -- GUID пользователя
    refresh_token_hash VARCHAR(255) NOT NULL,         -- Хеш refresh токена
    refresh_token_expires_at TIMESTAMP NOT NULL,      -- Время истечения refresh токена
    user_agent VARCHAR(500) NOT NULL,                 -- Для проверки безопасности
    ip_address VARCHAR(45) NOT NULL,                  -- Для webhook уведомлений
    is_active BOOLEAN NOT NULL DEFAULT TRUE,          -- Для деавторизации
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);