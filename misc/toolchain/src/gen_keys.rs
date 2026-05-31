use base64::{engine::general_purpose::URL_SAFE_NO_PAD, Engine};
use rand::RngCore;
use sha2::{Digest, Sha256};

fn generate(prefix: &str, n: usize) -> (String, String) {
    let mut bytes = vec![0u8; n];
    rand::thread_rng().fill_bytes(&mut bytes);
    let raw = format!("{}{}", prefix, URL_SAFE_NO_PAD.encode(&bytes));
    let hash = hex::encode(Sha256::digest(raw.as_bytes()));
    (raw, hash)
}

fn main() {
    let args: Vec<String> = std::env::args().collect();
    let kind = args.get(1).map(|s| s.as_str()).unwrap_or("sk");

    match kind {
        "sk" | "station-key" => {
            let (raw, hash) = generate("sk_", 32);
            println!("station key");
            println!("  raw  : {raw}");
            println!("  hash : {hash}");
        }
        "at" | "access-token" => {
            let (raw, hash) = generate("at_", 48);
            println!("access token");
            println!("  raw  : {raw}");
            println!("  hash : {hash}");
        }
        _ => {
            eprintln!("usage: gen-key [sk|at]");
            std::process::exit(1);
        }
    }
}