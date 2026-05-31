use crate::state::{ApiEvent, SavedStation, Station};
use futures_util::StreamExt;
use std::sync::mpsc::Sender;
use tokio_tungstenite::connect_async;
use tokio_tungstenite::tungstenite::Message;

const API_BASE: &str = "http://localhost:5000";
const WS_BASE: &str = "ws://localhost:5000";

pub fn connect_websocket(tx: Sender<ApiEvent>) {
    tokio::spawn(async move {
        match connect_async(format!("{}/v1/ws", WS_BASE)).await {
            Ok((mut ws, _)) => {
                while let Some(msg) = ws.next().await {
                    match msg {
                        Ok(Message::Text(text)) => {
                            if let Ok(val) = serde_json::from_str::<serde_json::Value>(&text) {
                                match val["type"].as_str().unwrap_or("") {
                                    "session" => {
                                        if let Some(id) = val["session_id"].as_str() {
                                            let _ = tx.send(ApiEvent::SessionReceived(id.to_string()));
                                        }
                                    }
                                    "viewer_count" => {
                                        if let (Some(sid), Some(count)) = (
                                            val["station_id"].as_str(),
                                            val["count"].as_u64(),
                                        ) {
                                            let _ = tx.send(ApiEvent::ViewerCount {
                                                station_id: sid.to_string(),
                                                count: count as u32,
                                            });
                                        }
                                    }
                                    _ => {}
                                }
                            }
                        }
                        Err(_) => break,
                        _ => {}
                    }
                }
                let _ = tx.send(ApiEvent::Error("websocket disconnected".into()));
            }
            Err(e) => {
                let _ = tx.send(ApiEvent::Error(format!("ws connect failed: {}", e)));
            }
        }
    });
}

pub fn fetch_directory(session_id: String, tx: Sender<ApiEvent>) {
    tokio::spawn(async move {
        let client = reqwest::Client::new();
        match client
            .get(format!("{}/v1/stations", API_BASE))
            .header("X-Session-ID", &session_id)
            .send()
            .await
        {
            Ok(r) => match r.json::<Vec<Station>>().await {
                Ok(stations) => {
                    let _ = tx.send(ApiEvent::DirectoryLoaded(stations));
                }
                Err(e) => {
                    let _ = tx.send(ApiEvent::Error(format!("parse error: {}", e)));
                }
            },
            Err(e) => {
                let _ = tx.send(ApiEvent::Error(format!("fetch failed: {}", e)));
            }
        }
    });
}

pub fn lookup_station_by_key(station_key: String, session_id: String, tx: Sender<ApiEvent>) {
    tokio::spawn(async move {
        let client = reqwest::Client::new();
        match client
            .get(format!("{}/v1/stations", API_BASE))
            .header("X-Session-ID", &session_id)
            .send()
            .await
        {
            Ok(r) => match r.json::<Vec<Station>>().await {
                Ok(stations) => {
                    if let Some(s) = stations.into_iter().find(|s| s.id == station_key || s.slug == station_key) {
                        let _ = tx.send(ApiEvent::StationLookedUp(SavedStation {
                            station_key: station_key.clone(),
                            name: s.display_name,
                            description: s.description,
                            genre: s.genre,
                        }));
                    } else {
                        let _ = tx.send(ApiEvent::Error("station not found".into()));
                    }
                }
                Err(_) => {
                    let _ = tx.send(ApiEvent::Error("could not parse response".into()));
                }
            },
            Err(_) => {
                let _ = tx.send(ApiEvent::Error("request failed".into()));
            }
        }
    });
}

pub fn fetch_owner_dashboard(access_token: String, session_id: String, tx: Sender<ApiEvent>) {
    tokio::spawn(async move {
        let token = access_token.trim().replace("\r", "").replace("\n", "");
        let session = session_id.trim().replace("\r", "").replace("\n", "");
        let client = reqwest::Client::new();
        match client
            .get(format!("{}/v1/owner/dashboard", API_BASE))
            .header("X-Session-ID", &session)
            .header("X-Access-Token", &token)
            .send()
            .await
        {
            Ok(r) => match r.json::<Station>().await {
                Ok(station) => {
                    let _ = tx.send(ApiEvent::OwnerDashboard(station));
                }
                Err(e) => {
                    let _ = tx.send(ApiEvent::Error(format!("invalid token or parse error: {}", e)));
                }
            },
            Err(e) => {
                let _ = tx.send(ApiEvent::Error(format!("request failed: {}", e)));
            }
        }
    });
}

pub fn create_room(access_token: String, session_id: String, tx: Sender<ApiEvent>) {
    tokio::spawn(async move {
        let client = reqwest::Client::new();
        match client
            .post(format!("{}/v1/room/create", API_BASE))
            .header("X-Session-ID", &session_id)
            .header("X-Access-Token", &access_token)
            .send()
            .await
        {
            Ok(r) => match r.json::<serde_json::Value>().await {
                Ok(v) => {
                    let _ = tx.send(ApiEvent::RoomCreated {
                        room_id: v["room_id"].as_str().unwrap_or("").to_string(),
                        code:    v["code"].as_str().unwrap_or("").to_string(),
                        link:    v["link"].as_str().unwrap_or("").to_string(),
                    });
                }
                Err(e) => { let _ = tx.send(ApiEvent::Error(e.to_string())); }
            },
            Err(e) => { let _ = tx.send(ApiEvent::Error(e.to_string())); }
        }
    });
}

pub fn join_room(code: String, name: String, session_id: String, tx: Sender<ApiEvent>) {
    tokio::spawn(async move {
        let client = reqwest::Client::new();
        match client
            .post(format!("{}/v1/room/join/{}", API_BASE, code))
            .header("X-Session-ID", &session_id)
            .json(&serde_json::json!({"name": name}))
            .send()
            .await
        {
            Ok(r) => match r.json::<crate::state::RoomInfo>().await {
                Ok(info) => { let _ = tx.send(ApiEvent::RoomJoined(info)); }
                Err(e)   => { let _ = tx.send(ApiEvent::Error(format!("join failed: {}", e))); }
            },
            Err(e) => { let _ = tx.send(ApiEvent::Error(e.to_string())); }
        }
    });
}

pub fn mute_guest(room_id: String, guest_id: String, muted: bool, access_token: String, session_id: String) {
    tokio::spawn(async move {
        let client = reqwest::Client::new();
        let _ = client
            .post(format!("{}/v1/room/{}/mute/{}", API_BASE, room_id, guest_id))
            .header("X-Session-ID", &session_id)
            .header("X-Access-Token", &access_token)
            .json(&serde_json::json!({"muted": muted}))
            .send()
            .await;
    });
}

pub fn self_mute(room_id: String, guest_id: String, muted: bool, session_id: String) {
    tokio::spawn(async move {
        let client = reqwest::Client::new();
        let _ = client
            .post(format!("{}/v1/room/{}/selfmute", API_BASE, room_id))
            .header("X-Session-ID", &session_id)
            .header("X-Guest-ID", &guest_id)
            .json(&serde_json::json!({"muted": muted}))
            .send()
            .await;
    });
}