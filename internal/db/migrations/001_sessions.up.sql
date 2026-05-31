CREATE TABLE sessions (
    id          CHAR(36)    PRIMARY KEY,
    created_at  DATETIME    NOT NULL DEFAULT NOW(),
    expires_at  DATETIME    NOT NULL,
    ip_address  VARCHAR(45),
    user_agent  TEXT,
    invalidated TINYINT(1)  NOT NULL DEFAULT 0
);
