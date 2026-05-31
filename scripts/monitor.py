#!/usr/bin/env python3
"""
cast.onion stream health monitor
checks all active stations and logs viewer counts and live status
"""

import json
import time
import urllib.request
import urllib.error
from datetime import datetime

API_BASE = "http://localhost:5000/v1"
POLL_INTERVAL = 10  # seconds


def fetch(path):
    try:
        with urllib.request.urlopen(f"{API_BASE}{path}", timeout=5) as r:
            return json.loads(r.read())
    except urllib.error.URLError as e:
        return {"error": str(e)}
    except Exception as e:
        return {"error": str(e)}


def log(msg):
    ts = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    print(f"[{ts}] {msg}")


def monitor():
    log("cast.onion stream health monitor starting")
    log(f"polling every {POLL_INTERVAL}s → {API_BASE}")
    print()

    while True:
        stations = fetch("/stations")
        if isinstance(stations, dict) and "error" in stations:
            log(f"api unreachable: {stations['error']}")
            time.sleep(POLL_INTERVAL)
            continue

        if not stations:
            log("no active stations")
        else:
            for s in stations:
                sid = s.get("ID", s.get("id", "?"))
                name = s.get("DisplayName", s.get("display_name", "?"))
                status_data = fetch(f"/status/{sid}")
                live = status_data.get("live", False)
                viewers = status_data.get("viewers", 0)
                status = "🔴 live" if live else "⚫ offline"
                log(f"{status}  {name:<30} viewers: {viewers}")

        print()
        time.sleep(POLL_INTERVAL)


if __name__ == "__main__":
    try:
        monitor()
    except KeyboardInterrupt:
        log("monitor stopped")