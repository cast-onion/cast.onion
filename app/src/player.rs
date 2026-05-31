use std::process::{Child, Command};

pub struct Player {
    child: Child,
}

impl Player {
    pub fn play(station_id: &str) -> Result<Self, String> {
        let url = format!("http://localhost:5000/v1/listen/{}", station_id);

        let child = Command::new("ffplay")
            .args([
                "-nodisp",
                "-autoexit",
                "-loglevel", "quiet",
                &url,
            ])
            .spawn()
            .map_err(|e| format!("ffplay failed to start: {}", e))?;

        Ok(Player { child })
    }

    pub fn stop(&mut self) {
        let _ = self.child.kill();
    }
}