use chrono::{DateTime, Duration, Utc};

pub fn now() -> DateTime<Utc> {
    Utc::now()
}

pub fn now_plus(seconds: i64) -> DateTime<Utc> {
    Utc::now() + Duration::seconds(seconds)
}

pub fn is_expired(ts: &DateTime<Utc>) -> bool {
    *ts < Utc::now()
}

pub fn human_ago(ts: &DateTime<Utc>) -> String {
    let d = Utc::now().signed_duration_since(*ts);
    let secs = d.num_seconds();
    if secs < 60 { return "just now".into(); }
    if secs < 3600 { return format!("{}m ago", secs / 60); }
    if secs < 86400 { return format!("{}h ago", secs / 3600); }
    format!("{}d ago", secs / 86400)
}