pub mod db;
pub mod types;

pub use db::LogDb;
pub use types::{EventKind, LogEvent};

use std::sync::Arc;
use tokio::sync::mpsc;

#[derive(Clone)]
pub struct Logger {
    tx: mpsc::UnboundedSender<LogEvent>,
}

impl Logger {
    pub fn new(db: Arc<LogDb>) -> Self {
        let (tx, mut rx) = mpsc::unbounded_channel::<LogEvent>();

        tokio::spawn(async move {
            while let Some(event) = rx.recv().await {
                let kind = event.kind.to_string();
                let msg = event.message.clone();
                if let Err(e) = db.insert(&event) {
                    eprintln!("[logger] db insert failed for {kind}: {e}");
                } else {
                    println!("[logger] {kind}: {msg}");
                }
            }
        });

        Logger { tx }
    }

    pub fn log(&self, event: LogEvent) {
        let _ = self.tx.send(event);
    }

    pub fn session_created(&self, session_id: &str, ip: &str) {
        self.log(
            LogEvent::new(EventKind::SessionCreated, format!("session created from {ip}"))
                .with_entity("session", session_id),
        );
    }

    pub fn application_submitted(&self, app_id: &str, station_name: &str, email: &str) {
        self.log(
            LogEvent::new(EventKind::ApplicationSubmitted, format!("application submitted for '{station_name}'"))
                .with_entity("application", app_id)
                .with_actor(email),
        );
    }

    pub fn application_approved(&self, app_id: &str, station_id: &str, admin: &str) {
        self.log(
            LogEvent::new(EventKind::ApplicationApproved, format!("application approved → station {station_id}"))
                .with_entity("application", app_id)
                .with_actor(admin),
        );
    }

    pub fn application_denied(&self, app_id: &str, admin: &str, reason: &str) {
        self.log(
            LogEvent::new(EventKind::ApplicationDenied, format!("application denied: {reason}"))
                .with_entity("application", app_id)
                .with_actor(admin),
        );
    }

    pub fn station_suspended(&self, station_id: &str, admin: &str, reason: &str) {
        self.log(
            LogEvent::new(EventKind::StationSuspended, format!("station suspended: {reason}"))
                .with_entity("station", station_id)
                .with_actor(admin),
        );
    }

    pub fn station_revoked(&self, station_id: &str, admin: &str, reason: &str) {
        self.log(
            LogEvent::new(EventKind::StationRevoked, format!("station revoked: {reason}"))
                .with_entity("station", station_id)
                .with_actor(admin),
        );
    }

    pub fn broadcast_started(&self, station_id: &str) {
        self.log(
            LogEvent::new(EventKind::BroadcastStarted, "broadcast started")
                .with_entity("station", station_id),
        );
    }

    pub fn broadcast_ended(&self, station_id: &str) {
        self.log(
            LogEvent::new(EventKind::BroadcastEnded, "broadcast ended")
                .with_entity("station", station_id),
        );
    }

    pub fn listener_joined(&self, station_id: &str, count: u32) {
        self.log(
            LogEvent::new(EventKind::ListenerJoined, format!("listener joined — {count} total"))
                .with_entity("station", station_id),
        );
    }

    pub fn listener_left(&self, station_id: &str, count: u32) {
        self.log(
            LogEvent::new(EventKind::ListenerLeft, format!("listener left — {count} remaining"))
                .with_entity("station", station_id),
        );
    }

    pub fn error(&self, context: &str, message: &str) {
        self.log(
            LogEvent::new(EventKind::Error, format!("[{context}] {message}")),
        );
    }
}