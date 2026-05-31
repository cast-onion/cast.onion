import { writable, get } from 'svelte/store'
import { browser } from '$app/environment'

export const session = writable<string | null>(null)
export const viewerCounts = writable<Record<string, number>>({})

export const API_BASE = 'http://localhost:5000'
const WS_URL = 'ws://localhost:5000/v1/ws'

export function connectSession(): void {
  if (!browser) return

  const ws = new WebSocket(WS_URL)

  ws.onopen = () => console.log('[cast.onion] ws connected')

  ws.onmessage = (e: MessageEvent) => {
    try {
      const msg = JSON.parse(e.data as string)
      if (msg.type === 'session') {
        session.set(msg.session_id as string)
      } else if (msg.type === 'viewer_count') {
        viewerCounts.update(c => ({ ...c, [msg.station_id]: msg.count }))
      }
    } catch {}
  }

  ws.onclose = () => {
    session.set(null)
    setTimeout(connectSession, 3000)
  }

  ws.onerror = () => ws.close()
}

export async function apiFetch(path: string, options: RequestInit = {}): Promise<Response> {
  const sid = get(session)
  return fetch(`${API_BASE}/v1${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'X-Session-ID': sid ?? '',
      ...((options.headers as Record<string, string>) ?? {}),
    }
  })
}