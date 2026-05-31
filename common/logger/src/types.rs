use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum EventKind {
    SessionCreated,
    SessionExpired,
    ApplicationSubmitted,
    ApplicationApproved,
    ApplicationDenied,
    StationCreated,
    StationUpdated,
    StationSuspended,
    StationUnsuspended,
    StationRevoked,
    BroadcastStarted,
    BroadcastEnded,
    ListenerJoined,
    ListenerLeft,
    RoomCreated,
    RoomJoined,
    RoomLeft,
    GuestMuted,
    GuestUnmuted,
    TokenIssued,
    TokenRevoked,
    AdminLogin,
    AdminAction,
    EmailSent,
    EmailFailed,
    Error,
}

impl std::fmt::Display for EventKind {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let s = serde_json::to_string(self).unwrap_or_default();
        write!(f, "{}", s.trim_matches('"'))
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LogEvent {
    pub id: String,
    pub kind: EventKind,
    pub entity_type: Option<String>,
    pub entity_id: Option<String>,
    pub actor: Option<String>,
    pub message: String,
    pub metadata: Option<serde_json::Value>,
    pub created_at: DateTime<Utc>,
}

impl LogEvent {
    pub fn new(kind: EventKind, message: impl Into<String>) -> Self {
        LogEvent {
            id: uuid::Uuid::new_v4().to_string(),
            kind,
            entity_type: None,
            entity_id: None,
            actor: None,
            message: message.into(),
            metadata: None,
            created_at: Utc::now(),
        }
    }

    pub fn with_entity(mut self, entity_type: impl Into<String>, entity_id: impl Into<String>) -> Self {
        self.entity_type = Some(entity_type.into());
        self.entity_id = Some(entity_id.into());
        self
    }

    pub fn with_actor(mut self, actor: impl Into<String>) -> Self {
        self.actor = Some(actor.into());
        self
    }

    pub fn with_metadata(mut self, meta: serde_json::Value) -> Self {
        self.metadata = Some(meta);
        self
    }
}