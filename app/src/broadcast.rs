use audiopus::{coder::Encoder as OpusEncoder, Application, Channels, SampleRate};
use cpal::traits::{DeviceTrait, HostTrait, StreamTrait};
use ogg::writing::PacketWriter;
use std::io::{Cursor, Read, Write};
use std::net::TcpStream;
use std::sync::atomic::{AtomicBool, Ordering};
use std::sync::{Arc, Mutex};

const API_HOST: &str = "localhost";
const API_PORT: u16 = 5000;
const OPUS_SAMPLE_RATE: u32 = 48000;
const OPUS_FRAME_SIZE: usize = 960;
const OPUS_CHANNELS: usize = 2;

pub fn list_input_devices() -> Vec<String> {
    if let Ok(out) = std::process::Command::new("pactl").args(["list", "sources"]).output() {
        let text = String::from_utf8_lossy(&out.stdout);
        let mut names = Vec::new();
        let mut current_name = String::new();
        let mut current_desc = String::new();
        for line in text.lines() {
            let line = line.trim();
            if line.starts_with("Name:") {
                current_name = line.trim_start_matches("Name:").trim().to_string();
                current_desc.clear();
            } else if line.starts_with("Description:") {
                current_desc = line.trim_start_matches("Description:").trim().to_string();
            } else if line.is_empty() && !current_name.is_empty() && !current_desc.is_empty() {
                if !current_name.contains(".monitor") {
                    names.push(format!("{} [{}]", current_desc, current_name));
                }
                current_name.clear();
                current_desc.clear();
            }
        }
        if !current_name.is_empty() && !current_desc.is_empty() && !current_name.contains(".monitor") {
            names.push(format!("{} [{}]", current_desc, current_name));
        }
        if !names.is_empty() { return names; }
    }
    let host = cpal::default_host();
    host.input_devices()
        .map(|devs| devs.filter_map(|d| d.name().ok()).collect())
        .unwrap_or_default()
}

pub fn extract_pactl_name(display: &str) -> String {
    if let Some(start) = display.rfind('[') {
        if let Some(end) = display.rfind(']') {
            return display[start+1..end].to_string();
        }
    }
    display.to_string()
}

pub struct Broadcaster {
    pub running: Arc<AtomicBool>,
}

impl Broadcaster {
    pub fn start(
        station_id: &str,
        station_key: &str,
        device_name: &str,
        on_error: impl Fn(String) + Send + 'static,
        on_stop: impl Fn() + Send + 'static,
    ) -> Result<Self, String> {
        let running = Arc::new(AtomicBool::new(true));
        let running_clone = running.clone();

        let station_id = station_id.to_string();
        let key_str = station_key.trim().to_string();
        let device_name = device_name.to_string();

        let mut tcp = TcpStream::connect(format!("{}:{}", API_HOST, API_PORT))
            .map_err(|e| format!("cannot connect to cast.onion: {}", e))?;

        let request = format!(
            "POST /v1/broadcast/{} HTTP/1.1\r\nHost: {}:{}\r\nX-Station-Key: {}\r\nContent-Type: audio/ogg\r\nExpect: 100-continue\r\n\r\n",
            station_id, API_HOST, API_PORT, key_str
        );
        tcp.write_all(request.as_bytes())
            .map_err(|e| format!("request failed: {}", e))?;

        tcp.set_read_timeout(Some(std::time::Duration::from_secs(5))).ok();
        let mut buf = [0u8; 512];
        let n = tcp.read(&mut buf).map_err(|e| format!("response failed: {}", e))?;
        let resp = String::from_utf8_lossy(&buf[..n]);
        if !resp.contains("100") && !resp.contains("200") {
            return Err(format!("server rejected: {}", resp.trim()));
        }
        tcp.set_read_timeout(None).ok();

        let tcp = Arc::new(Mutex::new(tcp));
        let tcp_clone = tcp.clone();

        std::thread::spawn(move || {
            let audio_host = cpal::default_host();
            let clean = extract_pactl_name(&device_name);

            let device = audio_host
                .input_devices()
                .ok()
                .and_then(|mut devs| devs.find(|d| {
                    d.name().ok().map(|n| n == clean || n.contains(&clean)).unwrap_or(false)
                }))
                .or_else(|| audio_host.default_input_device());

            let device = match device {
                Some(d) => d,
                None => { on_error("no input device found".into()); return; }
            };

            let config = cpal::StreamConfig {
                channels: OPUS_CHANNELS as u16,
                sample_rate: cpal::SampleRate(OPUS_SAMPLE_RATE),
                buffer_size: cpal::BufferSize::Fixed(OPUS_FRAME_SIZE as u32),
            };

            let encoder: OpusEncoder = match OpusEncoder::new(SampleRate::Hz48000, Channels::Stereo, Application::Audio) {
                Ok(e) => e,
                Err(e) => { on_error(format!("opus encoder error: {:?}", e)); return; }
            };
            let encoder = Arc::new(Mutex::new(encoder));
            let encoder_clone = encoder.clone();

            let audio_buf: Arc<Mutex<Vec<i16>>> = Arc::new(Mutex::new(Vec::new()));
            let audio_buf_clone = audio_buf.clone();
            let mut packet_no: u64 = 0;

            let input_stream = device.build_input_stream(
                &config,
                move |data: &[f32], _| {
                    let samples: Vec<i16> = data.iter()
                        .map(|&s| (s.clamp(-1.0, 1.0) * 0.8 * i16::MAX as f32) as i16)
                        .collect();
                    if let Ok(mut buf) = audio_buf_clone.lock() {
                        buf.extend_from_slice(&samples);
                    }
                },
                |e| eprintln!("stream error: {}", e),
                None,
            );

            let input_stream = match input_stream {
                Ok(s) => s,
                Err(_) => {
                    let fc = device.default_input_config().unwrap();
                    let ab = audio_buf.clone();
                    device.build_input_stream(
                        &fc.into(),
                        move |data: &[f32], _| {
                            let s: Vec<i16> = data.iter()
                                .map(|&s| (s.clamp(-1.0, 1.0) * 0.8 * i16::MAX as f32) as i16)
                                .collect();
                            if let Ok(mut b) = ab.lock() { b.extend_from_slice(&s); }
                        },
                        |e| eprintln!("stream error: {}", e),
                        None,
                    ).unwrap()
                }
            };

            if let Err(e) = input_stream.play() {
                on_error(format!("stream play failed: {}", e));
                return;
            }

            let mut ogg_buf = Vec::new();

            while running_clone.load(Ordering::Relaxed) {
                let frame: Vec<i16> = {
                    let mut buf = audio_buf.lock().unwrap();
                    let needed = OPUS_FRAME_SIZE * OPUS_CHANNELS;
                    if buf.len() < needed {
                        std::thread::sleep(std::time::Duration::from_millis(5));
                        continue;
                    }
                    buf.drain(..needed).collect()
                };

                let mut opus_out = vec![0u8; 4096];
                let n = {
                    let mut enc = encoder_clone.lock().unwrap();
                    match enc.encode(&frame, &mut opus_out) {
                        Ok(n) => n,
                        Err(_) => continue,
                    }
                };
                opus_out.truncate(n);

                ogg_buf.clear();
                {
                    let cursor = Cursor::new(&mut ogg_buf);
                    let mut writer = PacketWriter::new(cursor);
                    let granule = packet_no * OPUS_FRAME_SIZE as u64;
                    let _ = writer.write_packet(
                        opus_out.into_boxed_slice(),
                        0x1234,
                        ogg::writing::PacketWriteEndInfo::NormalPacket,
                        granule,
                    );
                }
                packet_no += 1;

                if !ogg_buf.is_empty() {
                    if let Ok(mut t) = tcp_clone.lock() {
                        if t.write_all(&ogg_buf).is_err() { break; }
                    }
                }
            }

            on_stop();
        });

        Ok(Broadcaster { running })
    }

    pub fn stop(&self) {
        self.running.store(false, Ordering::Relaxed);
    }
}