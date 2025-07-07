CREATE TABLE users (
    id uuid NOT NULL PRIMARY KEY,             -- GUID пользователя
    refresh_token_hash VARCHAR(255),         -- Хеш refresh токена
    refresh_token_expires_at TIMESTAMP,      -- Время истечения refresh токена
    user_agent VARCHAR(500),                 -- Для проверки безопасности
    ip_address VARCHAR(45),                  -- Для webhook уведомлений
    is_active BOOLEAN DEFAULT TRUE,          -- Для деавторизации
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);