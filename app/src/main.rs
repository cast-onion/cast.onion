mod app;
mod i18n;
mod api;
mod state;
mod broadcast;
mod player;
mod guest_broadcast;
mod credentials;

use eframe::egui;

fn main() -> eframe::Result<()> {
    let join_code = std::env::args()
        .nth(1)
        .and_then(|arg| {
            if arg.starts_with("cast://join/") {
                Some(arg.trim_start_matches("cast://join/").to_string())
            } else {
                None
            }
        });

    let options = eframe::NativeOptions {
        viewport: egui::ViewportBuilder::default()
            .with_title("cast.onion")
            .with_inner_size([920.0, 620.0])
            .with_min_inner_size([700.0, 480.0]),
        ..Default::default()
    };
    eframe::run_native(
        "cast.onion",
        options,
        Box::new(move |cc| Box::new(app::CastApp::new(cc, join_code.clone()))),
    )
}