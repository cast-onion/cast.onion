CREATE TABLE applications (
    id            CHAR(36)     PRIMARY KEY,
    session_id    CHAR(36)     NOT NULL REFERENCES sessions(id),
    contact_email VARCHAR(255) NOT NULL,
    station_name  VARCHAR(128) NOT NULL,
    description   TEXT         NOT NULL,
    genre         VARCHAR(64),
    notes         TEXT,
    status        ENUM('pending','approved','denied') NOT NULL DEFAULT 'pending',
    reviewed_by   CHAR(36)     REFERENCES users(id),
    reviewed_at   DATETIME,
    station_id    CHAR(36)     REFERENCES stations(id),
    created_at    DATETIME     NOT NULL DEFAULT NOW()
);
