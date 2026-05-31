CREATE TABLE access_tokens (
    id         CHAR(36)     PRIMARY KEY,
    station_id CHAR(36)     NOT NULL REFERENCES stations(id),
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    revoked    TINYINT(1)   NOT NULL DEFAULT 0,
    created_at DATETIME     NOT NULL DEFAULT NOW(),
    revoked_at DATETIME
);
