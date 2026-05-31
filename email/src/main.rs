use lettre::message::header::ContentType;
use lettre::transport::smtp::authentication::Credentials;
use lettre::{Message, SmtpTransport, Transport};
use minijinja::{Environment, context};
use serde::Deserialize;
use std::io::{self, Read};

#[derive(Debug, Deserialize)]
#[serde(tag = "type", rename_all = "snake_case")]
enum EmailRequest {
    ApplicationReceived {
        to: String,
        name: String,
        station_name: String,
        genre: String,
        application_id: String,
    },
    ApplicationApproved {
        to: String,
        name: String,
        station_name: String,
        station_key: String,
        access_token: String,
    },
    ApplicationDenied {
        to: String,
        name: String,
        station_name: String,
        reason: Option<String>,
    },
    StationSuspended {
        to: String,
        name: String,
        station_name: String,
        reason: Option<String>,
    },
    StationRevoked {
        to: String,
        name: String,
        station_name: String,
        reason: Option<String>,
    },
}

struct Config {
    smtp_host: String,
    smtp_port: u16,
    smtp_user: String,
    smtp_pass: String,
    from_email: String,
    from_name: String,
}

impl Config {
    fn from_env() -> Self {
        dotenvy::dotenv().ok();
        Config {
            smtp_host:  std::env::var("SMTP_HOST").unwrap_or_else(|_| "smtp.gmail.com".into()),
            smtp_port:  std::env::var("SMTP_PORT").unwrap_or_else(|_| "587".into()).parse().unwrap_or(587),
            smtp_user:  std::env::var("SMTP_USER").expect("SMTP_USER required"),
            smtp_pass:  std::env::var("SMTP_PASS").expect("SMTP_PASS required"),
            from_email: std::env::var("FROM_EMAIL").unwrap_or_else(|_| "no-reply@cast.onion".into()),
            from_name:  std::env::var("FROM_NAME").unwrap_or_else(|_| "cast.ouse logger::{LogDb, Logger};
use std::sync::Arc;

#[tokio::main]
async fn main() {
    dotenvy::dotenv().ok();
    let dsn = std::env::var("DATABASE_DSN").expect("DATABASE_DSN required");

    let db = Arc::new(LogDb::new(&dsn).expect("db connect failed"));
    db.migrate().expect("migration failed");

    let args: Vec<String> = std::env::args().collect();
    let kind    = args.get(1).map(|s| s.as_str());
    let entity  = args.get(2).map(|s| s.as_str());
    let limit: u32 = args.get(3).and_then(|s| s.parse().ok()).unwrap_or(50);

    match db.query(kind, entity, limit) {
        Ok(events) => {
            for e in &events {
                println!("[{}] {} | {} | {}",
                    e.created_at.format("%Y-%m-%d %H:%M:%S"),
                    e.kind,
                    e.entity_id.as_deref().unwrap_or("-"),
                    e.message,
                );
            }
            println!("\n{} events", events.len());
        }
        Err(e) => eprintln!("query failed: {e}"),
    }
}nion".into()),
        }
    }
}

fn load_template(env: &Environment, name: &str, ctx: minijinja::Value) -> (String, String) {
    let raw = env.get_template(name).unwrap().render(ctx).unwrap();
    let mut subject = String::new();
    let mut body = String::new();
    let mut past_subject = false;

    for line in raw.lines() {
        if !past_subject {
            if line.starts_with("Subject: ") {
                subject = line.trim_start_matches("Subject: ").to_string();
            } else if line.is_empty() && !subject.is_empty() {
                past_subject = true;
            }
        } else {
            body.push_str(line);
            body.push('\n');
        }
    }

    (subject, body.trim().to_string())
}

fn send(cfg: &Config, to: &str, subject: &str, body: &str) -> Result<(), String> {
    let from = format!("{} <{}>", cfg.from_name, cfg.from_email);

    let email = Message::builder()
        .from(from.parse().map_err(|e| format!("from parse error: {}", e))?)
        .to(to.parse().map_err(|e| format!("to parse error: {}", e))?)
        .subject(subject)
        .header(ContentType::TEXT_PLAIN)
        .body(body.to_string())
        .map_err(|e| format!("build error: {}", e))?;

    let creds = Credentials::new(cfg.smtp_user.clone(), cfg.smtp_pass.clone());

    let mailer = SmtpTransport::starttls_relay(&cfg.smtp_host)
        .map_err(|e| format!("relay error: {}", e))?
        .port(cfg.smtp_port)
        .credentials(creds)
        .build();

    mailer.send(&email).map_err(|e| format!("send error: {}", e))?;
    Ok(())
}

#[tokio::main]
async fn main() {
    let cfg = Config::from_env();

    let mut input = String::new();
    io::stdin().read_to_string(&mut input).expect("failed to read stdin");

    let request: EmailRequest = serde_json::from_str(&input).expect("invalid JSON input");

    let mut env = Environment::new();
    env.add_template("received.txt", include_str!("../templates/received.txt")).unwrap();
    env.add_template("approved.txt", include_str!("../templates/approved.txt")).unwrap();
    env.add_template("denied.txt",   include_str!("../templates/denied.txt")).unwrap();
    env.add_template("suspended.txt",include_str!("../templates/suspended.txt")).unwrap();
    env.add_template("revoked.txt",  include_str!("../templates/revoked.txt")).unwrap();

    let result = match &request {
        EmailRequest::ApplicationReceived { to, name, station_name, genre, application_id } => {
            let (subject, body) = load_template(&env, "received.txt", context! {
                name, station_name, genre, application_id
            });
            send(&cfg, to, &subject, &body)
        }
        EmailRequest::ApplicationApproved { to, name, station_name, station_key, access_token } => {
            let (subject, body) = load_template(&env, "approved.txt", context! {
                name, station_name, station_key, access_token
            });
            send(&cfg, to, &subject, &body)
        }
        EmailRequest::ApplicationDenied { to, name, station_name, reason } => {
            let (subject, body) = load_template(&env, "denied.txt", context! {
                name, station_name, reason => reason.as_deref().unwrap_or("")
            });
            send(&cfg, to, &subject, &body)
        }
        EmailRequest::StationSuspended { to, name, station_name, reason } => {
            let (subject, body) = load_template(&env, "suspended.txt", context! {
                name, station_name, reason => reason.as_deref().unwrap_or("")
            });
            send(&cfg, to, &subject, &body)
        }
        EmailRequest::StationRevoked { to, name, station_name, reason } => {
            let (subject, body) = load_template(&env, "revoked.txt", context! {
                name, station_name, reason => reason.as_deref().unwrap_or("")
            });
            send(&cfg, to, &subject, &body)
        }
    };

    match result {
        Ok(()) => {
            println!(r#"{{"ok":true}}"#);
        }
        Err(e) => {
            eprintln!("email error: {}", e);
            println!(r#"{{"ok":false,"error":"{}"}}"#, e);
            std::process::exit(1);
        }
    }
}
