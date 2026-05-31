use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, PartialEq)]
pub enum Tab {
    Directory,
    SavedStations,
    Broadcast,
    Guest,
}

#[derive(Debug, Clone, PartialEq)]
pub enum BroadcastMode {
    None,
    Microphone,
    External,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Station {
    #[serde(alias = "ID", alias = "id")]
    pub id: String,
    #[serde(alias = "Slug", alias = "slug")]
    pub slug: String,
    #[serde(alias = "DisplayName", alias = "display_name")]
    pub display_name: String,
    #[serde(alias = "Description", alias = "description", default)]
    pub description: String,
    #[serde(alias = "Genre", alias = "genre", default)]
    pub genre: String,
    #[serde(alias = "Status", alias = "status")]
    pub status: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SavedStation {
    pub station_key: String,
    pub name: String,
    pub description: String,
    pub genre: String,
}

#[derive(Debug, Clone)]
pub enum SessionState {
    Connecting,
    Connected(String),
    Failed,
}

#[derive(Debug, Clone)]
pub enum ApiEvent {
    SessionReceived(String),
    DirectoryLoaded(Vec<Station>),
    StationLookedUp(SavedStation),
    OwnerDashboard(Station),
    BroadcastStopped,
    BroadcastError(String),
    RoomCreated { room_id: String, code: String, link: String },
    RoomJoined(RoomInfo),
    RoomUpdated(Vec<GuestInfo>),
    ViewerCount { station_id: String, count: u32 },
    Error(String),
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GuestInfo {
    pub id: String,
    pub name: String,
    pub muted_by_host: bool,
    pub muted_self: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RoomInfo {
    pub room_id: String,
    pub guest_id: String,
    pub station_id: String,
    pub guests: Vec<GuestInfo>,
}