CREATE TABLE stations (
    id           CHAR(36)     PRIMARY KEY,
    slug         VARCHAR(64)  NOT NULL UNIQUE,
    display_name VARCHAR(128) NOT NULL,
    description  TEXT,
    genre        VARCHAR(64),
    website_url  VARCHAR(512),
    art_key      VARCHAR(512),
    status       ENUM('active','suspended','revoked') NOT NULL DEFAULT 'active',
    created_at   DATETIME     NOT NULL DEFAULT NOW(),
    updated_at   DATETIME     NOT NULL DEFAULT NOW() ON UPDATE NOW()
);
