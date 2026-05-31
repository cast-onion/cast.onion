use crate::types::LogEvent;
use mysql::prelude::*;
use mysql::*;

pub struct LogDb {
    pool: Pool,
}

impl LogDb {
    pub fn new(dsn: &str) -> Result<Self, mysql::Error> {
        let pool = Pool::new(dsn)?;
        Ok(LogDb { pool })
    }

    pub fn migrate(&self) -> Result<(), mysql::Error> {
        let mut conn = self.pool.get_conn()?;
        conn.query_drop(r"
            CREATE TABLE IF NOT EXISTS system_logs (
                id          CHAR(36)     PRIMARY KEY,
                kind        VARCHAR(64)  NOT NULL,
                entity_type VARCHAR(64),
                entity_id   CHAR(36),
                actor       VARCHAR(128),
                message     TEXT         NOT NULL,
                metadata    JSON,
                created_at  DATETIME(3)  NOT NULL DEFAULT NOW(3),
                INDEX idx_kind       (kind),
                INDEX idx_entity     (entity_type, entity_id),
                INDEX idx_created_at (created_at)
            )
        ")?;
        Ok(())
    }

    pub fn insert(&self, event: &LogEvent) -> Result<(), mysql::Error> {
        let mut conn = self.pool.get_conn()?;
        let metadata_str = event.metadata.as_ref().map(|m| m.to_string());
        conn.exec_drop(
            r"INSERT INTO system_logs
              (id, kind, entity_type, entity_id, actor, message, metadata, created_at)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
            (
                &event.id,
                event.kind.to_string(),
                &event.entity_type,
                &event.entity_id,
                &event.actor,
                &event.message,
                metadata_str,
                event.created_at.naive_utc(),
            ),
        )?;
        Ok(())
    }

    pub fn query(
        &self,
        kind: Option<&str>,
        entity_id: Option<&str>,
        limit: u32,
    ) -> Result<Vec<LogEvent>, mysql::Error> {
        let mut conn = self.pool.get_conn()?;

        let rows: Vec<(String, String, Option<String>, Option<String>, Option<String>, String, Option<String>, String)> =
            match (kind, entity_id) {
                (Some(k), Some(eid)) => conn.exec(
                    "SELECT id, kind, entity_type, entity_id, actor, message, metadata, created_at FROM system_logs WHERE kind = ? AND entity_id = ? ORDER BY created_at DESC LIMIT ?",
                    (k, eid, limit),
                )?,
                (Some(k), None) => conn.exec(
                    "SELECT id, kind, entity_type, entity_id, actor, message, metadata, created_at FROM system_logs WHERE kind = ? ORDER BY created_at DESC LIMIT ?",
                    (k, limit),
                )?,
                _ => conn.exec(
                    "SELECT id, kind, entity_type, entity_id, actor, message, metadata, created_at FROM system_logs ORDER BY created_at DESC LIMIT ?",
                    (limit,),
                )?,
            };

        Ok(rows.into_iter().map(|(id, kind_str, et, eid, actor, msg, meta, ts)| {
            LogEvent {
                id,
                kind: serde_json::from_str(&format!("\"{}\"", kind_str))
                    .unwrap_or(crate::types::EventKind::Error),
                entity_type: et,
                entity_id: eid,
                actor,
                message: msg,
                metadata: meta.and_then(|m| serde_json::from_str(&m).ok()),
                created_at: chrono::NaiveDateTime::parse_from_str(&ts, "%Y-%m-%d %H:%M:%S%.f")
                    .map(|dt| chrono::DateTime::from_naive_utc_and_offset(dt, Utc))
                    .unwrap_or_else(|_| Utc::now()),
            }
        }).collect())
    }
}