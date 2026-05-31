use logger::{LogDb, Logger};
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
}