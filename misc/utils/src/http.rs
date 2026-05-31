use serde::de::DeserializeOwned;

pub struct Client {
    inner: reqwest::Client,
    base: String,
    session_id: Option<String>,
}

impl Client {
    pub fn new(base: &str) -> Self {
        Client {
            inner: reqwest::Client::new(),
            base: base.trim_end_matches('/').to_string(),
            session_id: None,
        }
    }

    pub fn with_session(mut self, session_id: impl Into<String>) -> Self {
        self.session_id = Some(session_id.into());
        self
    }

    pub async fn get<T: DeserializeOwned>(&self, path: &str) -> Result<T, String> {
        let mut req = self.inner.get(format!("{}{}", self.base, path));
        if let Some(sid) = &self.session_id {
            req = req.header("X-Session-ID", sid);
        }
        let r = req.send().await.map_err(|e| e.to_string())?;
        if !r.status().is_success() {
            return Err(format!("HTTP {}", r.status()));
        }
        r.json::<T>().await.map_err(|e| e.to_string())
    }

    pub async fn post<B: serde::Serialize, T: DeserializeOwned>(
        &self,
        path: &str,
        body: &B,
    ) -> Result<T, String> {
        let mut req = self
            .inner
            .post(format!("{}{}", self.base, path))
            .json(body);
        if let Some(sid) = &self.session_id {
            req = req.header("X-Session-ID", sid);
        }
        let r = req.send().await.map_err(|e| e.to_string())?;
        if !r.status().is_success() {
            return Err(format!("HTTP {}", r.status()));
        }
        r.json::<T>().await.map_err(|e| e.to_string())
    }
}