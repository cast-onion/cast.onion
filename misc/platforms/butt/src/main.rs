use futures_util::StreamExt;
use tokio_tungstenite::connect_async;
use tokio_tungstenite::tungstenite::Message;

pub async fn acquire_session(api_base: &str) -> Result<String, String> {
    let ws_url = api_base
        .replace("https://", "wss://")
        .replace("http://", "ws://");
    let ws_url = format!("{}/v1/ws", ws_url);

    let (mut ws, _) = connect_async(&ws_url)
        .await
        .map_err(|e| format!("ws connect failed: {e}"))?;

    while let Some(msg) = ws.next().await {
        match msg {
            Ok(Message::Text(text)) => {
                if let Ok(val) = serde_json::from_str::<serde_json::Value>(&text) {
                    if val["type"] == "session" {
                        if let Some(id) = val["session_id"].as_str() {
                            return Ok(id.to_string());
                        }
                    }
                }
            }
            Err(e) => return Err(format!("ws error: {e}")),
            _ => {}
        }
    }

    Err("ws closed without sending session".into())
}