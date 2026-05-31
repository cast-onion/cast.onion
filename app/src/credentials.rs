use serde::{Deserialize, Serialize};
use std::path::PathBuf;

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub struct SavedCredential {
    pub label: String,
    pub station_name: String,
    pub access_token: String,
    pub station_key: String,
}

#[derive(Debug, Default, Serialize, Deserialize)]
pub struct CredentialStore {
    pub credentials: Vec<SavedCredential>,
}

impl CredentialStore {
    fn path() -> PathBuf {
        let base = dirs_next::config_dir()
            .unwrap_or_else(|| PathBuf::from("."));
        base.join("cast-onion").join("credentials.json")
    }

    pub fn load() -> Self {
        let path = Self::path();
        if !path.exists() {
            return Self::default();
        }
        std::fs::read_to_string(&path)
            .ok()
            .and_then(|s| serde_json::from_str(&s).ok())
            .unwrap_or_default()
    }

    pub fn save(&self) {
        let path = Self::path();
        if let Some(parent) = path.parent() {
            let _ = std::fs::create_dir_all(parent);
        }
        if let Ok(json) = serde_json::to_string_pretty(self) {
            let _ = std::fs::write(&path, json);
        }
    }

    pub fn add(&mut self, cred: SavedCredential) {
        if let Some(existing) = self.credentials.iter_mut().find(|c| c.label == cred.label) {
            *existing = cred;
        } else {
            self.credentials.push(cred);
        }
        self.save();
    }

    pub fn remove(&mut self, label: &str) {
        self.credentials.retain(|c| c.label != label);
        self.save();
    }
}