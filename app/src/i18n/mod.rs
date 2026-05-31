use std::collections::HashMap;

pub struct Lang {
    strings: HashMap<String, String>,
}

impl Lang {
    pub fn load(code: &str) -> Self {
        let path = format!("lang/{}.json", code);
        let content = std::fs::read_to_string(&path)
            .or_else(|_| std::fs::read_to_string("lang/en.json"))
            .unwrap_or_else(|_| "{}".into());
        let strings: HashMap<String, String> =
            serde_json::from_str(&content).unwrap_or_default();
        Lang { strings }
    }

    pub fn t<'a>(&'a self, key: &'a str) -> &'a str {
        self.strings.get(key).map(|s| s.as_str()).unwrap_or(key)
    }
}

pub fn available_langs() -> Vec<(&'static str, &'static str)> {
    vec![("en", "English")]
}