CREATE TABLE admin_actions (
    id          CHAR(36)    PRIMARY KEY,
    admin_id    CHAR(36)    NOT NULL REFERENCES users(id),
    action      ENUM('approve','deny','suspend','revoke','unsuspend') NOT NULL,
    target_type ENUM('application','station') NOT NULL,
    target_id   CHAR(36)    NOT NULL,
    reason      TEXT,
    created_at  DATETIME    NOT NULL DEFAULT NOW()
);
