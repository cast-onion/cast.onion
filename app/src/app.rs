use crate::api;
use crate::broadcast::Broadcaster;
use crate::credentials::{CredentialStore, SavedCredential};
use crate::guest_broadcast::GuestBroadcaster;
use crate::player::Player;
use crate::state::{ApiEvent, BroadcastMode, GuestInfo, RoomInfo, SavedStation, SessionState, Station, Tab};
use eframe::egui;
use std::sync::mpsc::{self, Receiver, Sender};

pub struct CastApp {
    tx: Sender<ApiEvent>,
    rx: Receiver<ApiEvent>,
    session: SessionState,
    tab: Tab,

    directory: Vec<Station>,
    directory_loading: bool,
    listening_to: Option<String>,
    player: Option<Player>,

    saved_stations: Vec<SavedStation>,
    station_key_input: String,
    station_lookup_status: String,

    access_token_input: String,
    owner_station: Option<Station>,
    owner_status: String,
    broadcast_mode: BroadcastMode,
    broadcaster: Option<Broadcaster>,
    broadcast_status: String,
    station_key_input_broadcast: String,
    input_devices: Vec<String>,
    selected_device: usize,

    room_id: Option<String>,
    room_code: Option<String>,
    room_link: Option<String>,
    room_guests: Vec<GuestInfo>,

    join_code_input: String,
    guest_name_input: String,
    joined_room: Option<RoomInfo>,
    guest_broadcaster: Option<GuestBroadcaster>,
    guest_muted: bool,
    guest_status: String,

    viewer_counts: std::collections::HashMap<String, u32>,
    lang: crate::i18n::Lang,

    cred_store: CredentialStore,
    selected_cred: Option<usize>,
    save_cred_label: String,

    rt: tokio::runtime::Runtime,
}

impl CastApp {
    pub fn new(_cc: &eframe::CreationContext, join_code: Option<String>) -> Self {
        let (tx, rx) = mpsc::channel();
        let rt = tokio::runtime::Runtime::new().unwrap();
        let tx_clone = tx.clone();
        rt.spawn(async move { api::connect_websocket(tx_clone); });
        Self {
            tx, rx,
            session: SessionState::Connecting,
            tab: if join_code.is_some() { Tab::Guest } else { Tab::Directory },
            directory: vec![],
            directory_loading: false,
            listening_to: None,
            player: None,
            saved_stations: vec![],
            station_key_input: String::new(),
            station_lookup_status: String::new(),
            access_token_input: String::new(),
            owner_station: None,
            owner_status: String::new(),
            broadcast_mode: BroadcastMode::None,
            broadcaster: None,
            broadcast_status: String::new(),
            station_key_input_broadcast: String::new(),
            input_devices: crate::broadcast::list_input_devices(),
            selected_device: 0,
            room_id: None,
            room_code: None,
            room_link: None,
            room_guests: vec![],
            join_code_input: join_code.unwrap_or_default(),
            guest_name_input: String::new(),
            joined_room: None,
            guest_broadcaster: None,
            guest_muted: false,
            guest_status: String::new(),
            viewer_counts: std::collections::HashMap::new(),
            lang: crate::i18n::Lang::load("en"),
            cred_store: CredentialStore::load(),
            selected_cred: None,
            save_cred_label: String::new(),
            rt,
        }
    }

    fn poll_events(&mut self) {
        while let Ok(event) = self.rx.try_recv() {
            match event {
                ApiEvent::SessionReceived(id) => {
                    self.session = SessionState::Connected(id.clone());
                    self.directory_loading = true;
                    let tx = self.tx.clone();
                    self.rt.spawn(async move { api::fetch_directory(id, tx); });
                }
                ApiEvent::DirectoryLoaded(stations) => {
                    self.directory = stations;
                    self.directory_loading = false;
                }
                ApiEvent::StationLookedUp(s) => {
                    self.station_lookup_status = format!("saved: {}", s.name);
                    self.saved_stations.push(s);
                    self.station_key_input.clear();
                }
                ApiEvent::OwnerDashboard(station) => {
                    self.owner_station = Some(station);
                    self.owner_status = "connected".into();
                }
                ApiEvent::BroadcastStopped => {
                    self.broadcaster = None;
                    self.broadcast_mode = BroadcastMode::None;
                    self.broadcast_status = "stopped".into();
                }
                ApiEvent::BroadcastError(e) => {
                    self.broadcaster = None;
                    self.broadcast_mode = BroadcastMode::None;
                    self.broadcast_status = format!("error: {}", e);
                }
                ApiEvent::RoomCreated { room_id, code, link } => {
                    self.room_id = Some(room_id);
                    self.room_code = Some(code);
                    self.room_link = Some(link);
                }
                ApiEvent::RoomJoined(info) => {
                    self.room_guests = info.guests.clone();
                    self.joined_room = Some(info);
                    self.guest_status = "joined".into();
                }
                ApiEvent::RoomUpdated(guests) => {
                    self.room_guests = guests;
                }
                ApiEvent::ViewerCount { station_id, count } => {
                    self.viewer_counts.insert(station_id, count);
                }
                ApiEvent::Error(e) => {
                    if matches!(self.session, SessionState::Connecting) {
                        self.session = SessionState::Failed;
                    }
                    self.station_lookup_status = format!("error: {}", e);
                    self.owner_status = format!("error: {}", e);
                    self.guest_status = format!("error: {}", e);
                }
            }
        }
    }
}

impl eframe::App for CastApp {
    fn update(&mut self, ctx: &egui::Context, _frame: &mut eframe::Frame) {
        self.poll_events();
        ctx.request_repaint_after(std::time::Duration::from_millis(100));

        let accent = egui::Color32::from_rgb(200, 255, 0);
        let bg     = egui::Color32::from_rgb(10, 10, 10);
        let bg2    = egui::Color32::from_rgb(17, 17, 17);
        let border = egui::Color32::from_rgb(34, 34, 34);
        let muted  = egui::Color32::from_rgb(90, 90, 90);
        let text   = egui::Color32::from_rgb(224, 224, 224);
        let danger = egui::Color32::from_rgb(255, 68, 68);
        let ok     = egui::Color32::from_rgb(68, 204, 136);
        let warn   = egui::Color32::from_rgb(255, 170, 0);

        let mut style = (*ctx.style()).clone();
        style.visuals.panel_fill = bg;
        style.visuals.window_fill = bg2;
        style.visuals.extreme_bg_color = bg2;
        style.visuals.faint_bg_color = bg2;
        style.visuals.widgets.noninteractive.bg_fill = bg2;
        style.visuals.widgets.noninteractive.fg_stroke = egui::Stroke::new(1.0, text);
        style.visuals.widgets.inactive.bg_fill = bg2;
        style.visuals.widgets.inactive.fg_stroke = egui::Stroke::new(1.0, muted);
        style.visuals.widgets.hovered.bg_fill = egui::Color32::from_rgb(25, 25, 25);
        style.visuals.widgets.active.bg_fill = egui::Color32::from_rgb(30, 30, 30);
        style.visuals.override_text_color = Some(text);
        ctx.set_style(style);

        egui::TopBottomPanel::top("header")
            .frame(egui::Frame::none().fill(bg).inner_margin(egui::Margin::symmetric(20.0, 12.0)))
            .show(ctx, |ui| {
                ui.horizontal(|ui| {
                    ui.label(egui::RichText::new("cast.onion").color(accent).size(16.0).strong());
                    ui.add_space(24.0);
                    for (label, t) in [
                        ("directory", Tab::Directory),
                        ("saved", Tab::SavedStations),
                        ("broadcast", Tab::Broadcast),
                        ("guest", Tab::Guest),
                    ] {
                        let selected = self.tab == t;
                        if ui.selectable_label(selected, egui::RichText::new(label).size(13.0)).clicked() {
                            self.tab = t;
                        }
                        ui.add_space(4.0);
                    }
                    ui.with_layout(egui::Layout::right_to_left(egui::Align::Center), |ui| {
                        match &self.session {
                            SessionState::Connecting => { ui.label(egui::RichText::new("● connecting...").color(muted).size(12.0)); }
                            SessionState::Connected(id) => {
                                let short = &id[..8.min(id.len())];
                                ui.label(egui::RichText::new(format!("● {}", short)).color(ok).size(12.0));
                            }
                            SessionState::Failed => { ui.label(egui::RichText::new("● no connection").color(danger).size(12.0)); }
                        }
                    });
                });
            });

        egui::CentralPanel::default()
            .frame(egui::Frame::none().fill(bg).inner_margin(egui::Margin::symmetric(24.0, 20.0)))
            .show(ctx, |ui| {
                match self.tab {
                    Tab::Directory     => self.show_directory(ui, accent, bg2, border, muted, text, ok, danger),
                    Tab::SavedStations => self.show_saved(ui, accent, bg2, border, muted, text),
                    Tab::Broadcast     => self.show_broadcast(ui, accent, bg2, border, muted, text, ok, danger, warn),
                    Tab::Guest         => self.show_guest(ui, accent, bg2, border, muted, text, ok, danger),
                }
            });
    }
}

impl CastApp {
    fn show_directory(&mut self, ui: &mut egui::Ui, accent: egui::Color32, bg2: egui::Color32, border: egui::Color32, muted: egui::Color32, text: egui::Color32, ok: egui::Color32, danger: egui::Color32) {
        ui.label(egui::RichText::new("stations").color(text).size(18.0));
        ui.add_space(4.0);
        ui.label(egui::RichText::new("active stations on the network").color(muted).size(12.0));
        ui.add_space(16.0);

        if self.directory_loading {
            ui.label(egui::RichText::new("tuning in...").color(muted).size(13.0));
            return;
        }
        if self.directory.is_empty() {
            ui.label(egui::RichText::new("no active stations").color(muted).size(13.0));
            return;
        }

        egui::ScrollArea::vertical().show(ui, |ui| {
            for station in self.directory.clone() {
                let is_playing = self.listening_to.as_deref() == Some(&station.id);
                let frame = egui::Frame::none()
                    .fill(bg2)
                    .stroke(egui::Stroke::new(if is_playing { 1.5 } else { 1.0 }, if is_playing { accent } else { border }))
                    .inner_margin(egui::Margin::symmetric(16.0, 12.0))
                    .rounding(4.0);

                frame.show(ui, |ui| {
                    ui.horizontal(|ui| {
                        ui.vertical(|ui| {
                            ui.horizontal(|ui| {
                                if is_playing { ui.label(egui::RichText::new("▶ ").color(ok).size(12.0)); }
                                ui.label(egui::RichText::new(&station.display_name).color(text).size(14.0));
                            });
                            ui.add_space(2.0);
                            ui.label(egui::RichText::new(&station.description).color(muted).size(12.0));
                        });
                        ui.with_layout(egui::Layout::right_to_left(egui::Align::Center), |ui| {
                            if is_playing {
                                if ui.add(egui::Button::new(egui::RichText::new("stop").color(egui::Color32::WHITE).size(11.0)).fill(danger)).clicked() {
                                    if let Some(p) = &mut self.player { p.stop(); }
                                    self.player = None;
                                    self.listening_to = None;
                                }
                            } else {
                                let sid = station.id.clone();
                                if ui.add(egui::Button::new(egui::RichText::new("listen").color(egui::Color32::BLACK).size(11.0)).fill(accent)).clicked() {
                                    if let Some(p) = &mut self.player { p.stop(); }
                                    match Player::play(&sid) {
                                        Ok(p) => { self.player = Some(p); self.listening_to = Some(station.id.clone()); }
                                        Err(e) => { eprintln!("player error: {}", e); }
                                    }
                                }
                            }
                            ui.add_space(8.0);
                            if let Some(count) = self.viewer_counts.get(&station.id) {
                                if *count > 0 {
                                    ui.label(egui::RichText::new(format!("👂 {}", count)).color(muted).size(11.0));
                                    ui.add_space(4.0);
                                }
                            }
                            if !station.genre.is_empty() {
                                ui.label(egui::RichText::new(&station.genre).color(accent).size(11.0));
                            }
                        });
                    });
                });
                ui.add_space(2.0);
            }
        });
    }

    fn show_saved(&mut self, ui: &mut egui::Ui, accent: egui::Color32, bg2: egui::Color32, border: egui::Color32, muted: egui::Color32, text: egui::Color32) {
        ui.label(egui::RichText::new("saved stations").color(text).size(18.0));
        ui.add_space(16.0);
        let frame = egui::Frame::none().fill(bg2).stroke(egui::Stroke::new(1.0, border)).inner_margin(egui::Margin::symmetric(16.0, 12.0)).rounding(4.0);
        frame.show(ui, |ui| {
            ui.label(egui::RichText::new("station id").color(muted).size(12.0));
            ui.add_space(6.0);
            ui.horizontal(|ui| {
                ui.add(egui::TextEdit::singleline(&mut self.station_key_input).hint_text("enter station id...").desired_width(360.0).font(egui::TextStyle::Monospace));
                ui.add_space(8.0);
                let can = !self.station_key_input.is_empty() && matches!(&self.session, SessionState::Connected(_));
                if ui.add_enabled(can, egui::Button::new(egui::RichText::new("add").size(12.0))).clicked() {
                    if let SessionState::Connected(sid) = &self.session {
                        let tx = self.tx.clone();
                        let key = self.station_key_input.clone();
                        let sid = sid.clone();
                        self.station_lookup_status = "looking up...".into();
                        self.rt.spawn(async move { api::lookup_station_by_key(key, sid, tx); });
                    }
                }
            });
            if !self.station_lookup_status.is_empty() {
                ui.add_space(6.0);
                ui.label(egui::RichText::new(&self.station_lookup_status).color(muted).size(11.0));
            }
        });
        ui.add_space(16.0);
        let mut to_remove = None;
        for (i, s) in self.saved_stations.iter().enumerate() {
            let f = egui::Frame::none().fill(bg2).stroke(egui::Stroke::new(1.0, border)).inner_margin(egui::Margin::symmetric(16.0, 12.0)).rounding(4.0);
            f.show(ui, |ui| {
                ui.horizontal(|ui| {
                    ui.vertical(|ui| {
                        ui.label(egui::RichText::new(&s.name).color(text).size(14.0));
                        ui.label(egui::RichText::new(&s.description).color(muted).size(12.0));
                    });
                    ui.with_layout(egui::Layout::right_to_left(egui::Align::Min), |ui| {
                        if ui.small_button("✕").clicked() { to_remove = Some(i); }
                        if !s.genre.is_empty() { ui.label(egui::RichText::new(&s.genre).color(accent).size(11.0)); }
                    });
                });
            });
            ui.add_space(2.0);
        }
        if let Some(i) = to_remove { self.saved_stations.remove(i); }
        if self.saved_stations.is_empty() {
            ui.label(egui::RichText::new("no saved stations yet").color(muted).size(13.0));
        }
    }

    fn show_broadcast(&mut self, ui: &mut egui::Ui, accent: egui::Color32, bg2: egui::Color32, border: egui::Color32, muted: egui::Color32, text: egui::Color32, ok: egui::Color32, danger: egui::Color32, warn: egui::Color32) {
        ui.label(egui::RichText::new("broadcast").color(text).size(18.0));
        ui.add_space(16.0);

        let token_frame = egui::Frame::none().fill(bg2).stroke(egui::Stroke::new(1.0, border)).inner_margin(egui::Margin::symmetric(16.0, 12.0)).rounding(4.0);
        token_frame.show(ui, |ui| {

            if !self.cred_store.credentials.is_empty() {
                ui.label(egui::RichText::new("saved credentials").color(muted).size(11.0));
                ui.add_space(4.0);

                let selected_label = self.selected_cred
                    .and_then(|i| self.cred_store.credentials.get(i))
                    .map(|c| format!("{} — {}", c.label, c.station_name))
                    .unwrap_or_else(|| "select...".into());

                egui::ComboBox::from_id_source("saved_creds")
                    .selected_text(&selected_label)
                    .width(400.0)
                    .show_ui(ui, |ui| {
                        for (i, cred) in self.cred_store.credentials.iter().enumerate() {
                            let label = format!("{} — {}", cred.label, cred.station_name);
                            if ui.selectable_value(&mut self.selected_cred, Some(i), &label).clicked() {
                                self.access_token_input = cred.access_token.clone();
                                self.station_key_input_broadcast = cred.station_key.clone();
                            }
                        }
                    });

                if let Some(idx) = self.selected_cred {
                    if let Some(cred) = self.cred_store.credentials.get(idx).cloned() {
                        ui.add_space(4.0);
                        if ui.small_button("🗑 remove").clicked() {
                            self.cred_store.remove(&cred.label);
                            self.selected_cred = None;
                        }
                    }
                }

                ui.add_space(12.0);
                ui.separator();
                ui.add_space(12.0);
            }

            ui.label(egui::RichText::new("access token").color(muted).size(12.0));
            ui.add_space(6.0);
            ui.horizontal(|ui| {
                ui.add(egui::TextEdit::singleline(&mut self.access_token_input)
                    .hint_text("at_...")
                    .desired_width(400.0)
                    .password(true)
                    .font(egui::TextStyle::Monospace));
                ui.add_space(8.0);
                let can = !self.access_token_input.is_empty() && matches!(&self.session, SessionState::Connected(_));
                if ui.add_enabled(can, egui::Button::new(egui::RichText::new("connect").size(12.0))).clicked() {
                    if let SessionState::Connected(sid) = &self.session {
                        let tx = self.tx.clone();
                        let token = self.access_token_input.clone();
                        let sid = sid.clone();
                        self.owner_status = "connecting...".into();
                        self.rt.spawn(async move { api::fetch_owner_dashboard(token, sid, tx); });
                    }
                }
            });
            if !self.owner_status.is_empty() {
                ui.add_space(6.0);
                let color = if self.owner_status.starts_with("error") { danger } else { ok };
                ui.label(egui::RichText::new(&self.owner_status).color(color).size(11.0));
            }

            ui.add_space(10.0);

            ui.label(egui::RichText::new("station key (sk_...)").color(muted).size(12.0));
            ui.add_space(4.0);
            ui.add(egui::TextEdit::singleline(&mut self.station_key_input_broadcast)
                .hint_text("sk_...")
                .desired_width(400.0)
                .password(true)
                .font(egui::TextStyle::Monospace));

            if !self.access_token_input.is_empty() && !self.station_key_input_broadcast.is_empty() {
                ui.add_space(10.0);
                ui.horizontal(|ui| {
                    ui.label(egui::RichText::new("save as").color(muted).size(11.0));
                    ui.add_space(6.0);
                    ui.add(egui::TextEdit::singleline(&mut self.save_cred_label)
                        .hint_text("e.g. my station")
                        .desired_width(180.0));
                    ui.add_space(6.0);
                    let can_save = !self.save_cred_label.is_empty();
                    if ui.add_enabled(can_save, egui::Button::new(egui::RichText::new("save").size(11.0))).clicked() {
                        let station_name = self.owner_station
                            .as_ref()
                            .map(|s| s.display_name.clone())
                            .unwrap_or_else(|| "unknown".into());
                        self.cred_store.add(SavedCredential {
                            label: self.save_cred_label.clone(),
                            station_name,
                            access_token: self.access_token_input.clone(),
                            station_key: self.station_key_input_broadcast.clone(),
                        });
                        self.save_cred_label.clear();
                    }
                });
            }
        });

        ui.add_space(16.0);

        if let Some(station) = self.owner_station.clone() {
            let dash_frame = egui::Frame::none().fill(bg2).stroke(egui::Stroke::new(1.0, border)).inner_margin(egui::Margin::symmetric(20.0, 16.0)).rounding(4.0);
            dash_frame.show(ui, |ui| {
                ui.horizontal(|ui| {
                    ui.vertical(|ui| {
                        ui.label(egui::RichText::new(&station.display_name).color(text).size(16.0).strong());
                        if !station.description.is_empty() {
                            ui.label(egui::RichText::new(&station.description).color(muted).size(12.0));
                        }
                        if !station.genre.is_empty() {
                            ui.label(egui::RichText::new(&station.genre).color(accent).size(11.0));
                        }
                    });
                    ui.with_layout(egui::Layout::right_to_left(egui::Align::Min), |ui| {
                        let (sc, st) = match station.status.as_str() {
                            "active" => (ok, "active"), "suspended" => (warn, "suspended"), _ => (danger, "revoked"),
                        };
                        ui.label(egui::RichText::new(st).color(sc).size(12.0));
                    });
                });

                ui.add_space(12.0);
                ui.separator();
                ui.add_space(12.0);

                ui.label(egui::RichText::new("input device").color(muted).size(11.0));
                ui.add_space(4.0);
                let sel = self.input_devices.get(self.selected_device).cloned().unwrap_or_else(|| "default".into());
                egui::ComboBox::from_id_source("input_device").selected_text(&sel).width(400.0).show_ui(ui, |ui| {
                    for (i, name) in self.input_devices.iter().enumerate() {
                        ui.selectable_value(&mut self.selected_device, i, name);
                    }
                });
                ui.add_space(12.0);

                let is_live = self.broadcaster.is_some();

                if !is_live {
                    let can = !self.station_key_input_broadcast.is_empty();
                    if ui.add_enabled(can, egui::Button::new(egui::RichText::new("🎙  microphone").color(egui::Color32::BLACK).size(13.0)).fill(accent)).clicked() {
                        self.broadcast_mode = BroadcastMode::Microphone;
                        let key = self.station_key_input_broadcast.clone();
                        let tx = self.tx.clone();
                        let tx2 = self.tx.clone();
                        let sid = station.id.clone();
                        let raw = self.input_devices.get(self.selected_device).cloned().unwrap_or_default();
                        let dev = crate::broadcast::extract_pactl_name(&raw);
                        match Broadcaster::start(&sid, &key, &dev,
                            move |e| { let _ = tx.send(ApiEvent::BroadcastError(e)); },
                            move || { let _ = tx2.send(ApiEvent::BroadcastStopped); },
                        ) {
                            Ok(b) => { self.broadcaster = Some(b); self.broadcast_status = "on air".into(); }
                            Err(e) => { self.broadcast_status = format!("error: {}", e); self.broadcast_mode = BroadcastMode::None; }
                        }
                    }
                } else {
                    ui.horizontal(|ui| {
                        ui.label(egui::RichText::new("● on air").color(ok).size(14.0).strong());
                        ui.add_space(16.0);
                        if let Some(count) = self.viewer_counts.get(&station.id) {
                            ui.label(egui::RichText::new(format!("👂 {} listening", count)).color(muted).size(12.0));
                            ui.add_space(8.0);
                        }
                        if ui.add(egui::Button::new(egui::RichText::new("stop").color(egui::Color32::WHITE).size(12.0)).fill(danger)).clicked() {
                            if let Some(b) = &self.broadcaster { b.stop(); }
                            self.broadcaster = None;
                            self.broadcast_mode = BroadcastMode::None;
                            self.broadcast_status = "stopped".into();
                        }
                    });
                }

                if !self.broadcast_status.is_empty() {
                    ui.add_space(6.0);
                    let color = if self.broadcast_status.starts_with("error") { danger } else if self.broadcast_status == "on air" { ok } else { muted };
                    ui.label(egui::RichText::new(&self.broadcast_status).color(color).size(11.0));
                }

                ui.add_space(16.0);
                ui.separator();
                ui.add_space(12.0);
                ui.label(egui::RichText::new("guest room").color(muted).size(11.0));
                ui.add_space(8.0);

                if self.room_code.is_none() {
                    if ui.add(egui::Button::new(egui::RichText::new("create invite link").size(12.0)).fill(bg2).stroke(egui::Stroke::new(1.0, border))).clicked() {
                        if let SessionState::Connected(sid) = &self.session {
                            let tx = self.tx.clone();
                            let token = self.access_token_input.clone();
                            let sid = sid.clone();
                            self.rt.spawn(async move { api::create_room(token, sid, tx); });
                        }
                    }
                } else {
                    let code = self.room_code.clone().unwrap_or_default();
                    let link = self.room_link.clone().unwrap_or_default();
                    ui.horizontal(|ui| {
                        ui.label(egui::RichText::new("join code:").color(muted).size(12.0));
                        ui.label(egui::RichText::new(&code).color(accent).size(14.0).monospace().strong());
                    });
                    ui.add_space(4.0);
                    ui.label(egui::RichText::new(&link).color(muted).size(11.0).monospace());
                    ui.add_space(8.0);

                    if !self.room_guests.is_empty() {
                        ui.label(egui::RichText::new("guests").color(muted).size(11.0));
                        ui.add_space(4.0);
                        let room_id = self.room_id.clone().unwrap_or_default();
                        let token = self.access_token_input.clone();
                        let session_id = if let SessionState::Connected(s) = &self.session { s.clone() } else { String::new() };
                        for g in &self.room_guests {
                            let f = egui::Frame::none().fill(bg2).stroke(egui::Stroke::new(1.0, border)).inner_margin(egui::Margin::symmetric(12.0, 8.0)).rounding(4.0);
                            f.show(ui, |ui| {
                                ui.horizontal(|ui| {
                                    let mic_color = if g.muted_by_host || g.muted_self { danger } else { ok };
                                    ui.label(egui::RichText::new("●").color(mic_color).size(12.0));
                                    ui.add_space(6.0);
                                    ui.label(egui::RichText::new(&g.name).color(text).size(13.0));
                                    ui.with_layout(egui::Layout::right_to_left(egui::Align::Center), |ui| {
                                        let mute_label = if g.muted_by_host { "unmute" } else { "mute" };
                                        if ui.small_button(mute_label).clicked() {
                                            api::mute_guest(room_id.clone(), g.id.clone(), !g.muted_by_host, token.clone(), session_id.clone());
                                        }
                                    });
                                });
                            });
                            ui.add_space(2.0);
                        }
                    } else {
                        ui.label(egui::RichText::new("no guests yet — share the link").color(muted).size(12.0));
                    }
                }
            });
        }
    }

    fn show_guest(&mut self, ui: &mut egui::Ui, accent: egui::Color32, bg2: egui::Color32, border: egui::Color32, muted: egui::Color32, text: egui::Color32, ok: egui::Color32, danger: egui::Color32) {
        ui.label(egui::RichText::new("join as guest").color(text).size(18.0));
        ui.add_space(4.0);
        ui.label(egui::RichText::new("enter your invite code or open a cast:// link").color(muted).size(12.0));
        ui.add_space(16.0);

        if self.joined_room.is_none() {
            let frame = egui::Frame::none().fill(bg2).stroke(egui::Stroke::new(1.0, border)).inner_margin(egui::Margin::symmetric(16.0, 12.0)).rounding(4.0);
            frame.show(ui, |ui| {
                ui.label(egui::RichText::new("your name").color(muted).size(12.0));
                ui.add_space(4.0);
                ui.add(egui::TextEdit::singleline(&mut self.guest_name_input).hint_text("e.g. brody").desired_width(300.0));
                ui.add_space(12.0);

                ui.label(egui::RichText::new("invite code").color(muted).size(12.0));
                ui.add_space(4.0);
                ui.horizontal(|ui| {
                    ui.add(egui::TextEdit::singleline(&mut self.join_code_input).hint_text("8 character code...").desired_width(200.0).font(egui::TextStyle::Monospace));
                    ui.add_space(8.0);
                    let can = !self.join_code_input.is_empty() && matches!(&self.session, SessionState::Connected(_));
                    if ui.add_enabled(can, egui::Button::new(egui::RichText::new("join room").color(egui::Color32::BLACK).size(13.0)).fill(accent)).clicked() {
                        if let SessionState::Connected(sid) = &self.session {
                            let tx = self.tx.clone();
                            let code = self.join_code_input.clone();
                            let name = if self.guest_name_input.is_empty() { "guest".into() } else { self.guest_name_input.clone() };
                            let sid = sid.clone();
                            self.guest_status = "joining...".into();
                            self.rt.spawn(async move { api::join_room(code, name, sid, tx); });
                        }
                    }
                });

                if !self.guest_status.is_empty() {
                    ui.add_space(6.0);
                    let color = if self.guest_status.starts_with("error") { danger } else { ok };
                    ui.label(egui::RichText::new(&self.guest_status).color(color).size(11.0));
                }
            });
        } else {
            let room = self.joined_room.clone().unwrap();
            let dash = egui::Frame::none().fill(bg2).stroke(egui::Stroke::new(1.0, border)).inner_margin(egui::Margin::symmetric(20.0, 16.0)).rounding(4.0);
            dash.show(ui, |ui| {
                ui.label(egui::RichText::new("you're in the room").color(ok).size(16.0).strong());
                ui.add_space(4.0);
                ui.label(egui::RichText::new(format!("room: {}", room.room_id)).color(muted).size(11.0).monospace());
                ui.add_space(16.0);

                let is_live = self.guest_broadcaster.is_some();

                ui.label(egui::RichText::new("input device").color(muted).size(11.0));
                ui.add_space(4.0);
                let sel = self.input_devices.get(self.selected_device).cloned().unwrap_or_else(|| "default".into());
                egui::ComboBox::from_id_source("guest_device").selected_text(&sel).width(360.0).show_ui(ui, |ui| {
                    for (i, name) in self.input_devices.iter().enumerate() {
                        ui.selectable_value(&mut self.selected_device, i, name);
                    }
                });
                ui.add_space(12.0);

                if !is_live {
                    if ui.add(egui::Button::new(egui::RichText::new("🎙  go live").color(egui::Color32::BLACK).size(13.0)).fill(accent)).clicked() {
                        let raw = self.input_devices.get(self.selected_device).cloned().unwrap_or_default();
                        let dev = crate::broadcast::extract_pactl_name(&raw);
                        match GuestBroadcaster::start(
                            &room.room_id, &room.guest_id, &dev,
                            |e| eprintln!("guest broadcast error: {}", e),
                            || eprintln!("guest broadcast stopped"),
                        ) {
                            Ok(b) => { self.guest_broadcaster = Some(b); self.guest_status = "broadcasting".into(); }
                            Err(e) => { self.guest_status = format!("error: {}", e); }
                        }
                    }
                } else {
                    ui.horizontal(|ui| {
                        ui.label(egui::RichText::new("● live").color(ok).size(14.0).strong());
                        ui.add_space(12.0);

                        let mute_label = if self.guest_muted { "unmute" } else { "mute" };
                        let mute_color = if self.guest_muted { danger } else { bg2 };
                        if ui.add(egui::Button::new(egui::RichText::new(mute_label).size(12.0)).fill(mute_color).stroke(egui::Stroke::new(1.0, border))).clicked() {
                            self.guest_muted = !self.guest_muted;
                            if let SessionState::Connected(sid) = &self.session {
                                api::self_mute(room.room_id.clone(), room.guest_id.clone(), self.guest_muted, sid.clone());
                            }
                        }

                        ui.add_space(8.0);
                        if ui.add(egui::Button::new(egui::RichText::new("leave").color(egui::Color32::WHITE).size(12.0)).fill(danger)).clicked() {
                            if let Some(b) = &self.guest_broadcaster { b.stop(); }
                            self.guest_broadcaster = None;
                            self.joined_room = None;
                            self.guest_status = String::new();
                            self.guest_muted = false;
                        }
                    });
                }

                ui.add_space(16.0);
                ui.separator();
                ui.add_space(8.0);
                ui.label(egui::RichText::new("guests in room").color(muted).size(11.0));
                ui.add_space(8.0);

                for g in &self.room_guests {
                    let f = egui::Frame::none().fill(bg2).stroke(egui::Stroke::new(1.0, border)).inner_margin(egui::Margin::symmetric(12.0, 8.0)).rounding(4.0);
                    f.show(ui, |ui| {
                        ui.horizontal(|ui| {
                            let mic_color = if g.muted_by_host || g.muted_self { danger } else { ok };
                            ui.label(egui::RichText::new("●").color(mic_color).size(12.0));
                            ui.add_space(6.0);
                            let is_you = g.id == room.guest_id;
                            let label = if is_you { format!("{} (you)", g.name) } else { g.name.clone() };
                            ui.label(egui::RichText::new(label).color(text).size(13.0));
                        });
                    });
                    ui.add_space(2.0);
                }
            });
        }
    }
}