CREATE TABLE users (
    id            CHAR(36)     PRIMARY KEY,
    username      VARCHAR(64)  NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role          ENUM('admin') NOT NULL DEFAULT 'admin',
    created_at    DATETIME     NOT NULL DEFAULT NOW()
);
