import { writable } from 'svelte/store'
import { API_BASE } from './session'

export const adminUser = writable<string | null>(null)

export async function adminLogin(username: string, password: string): Promise<boolean> {
  const r = await fetch(`${API_BASE}/admin/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: new URLSearchParams({ username, password }),
    credentials: 'include',
  })
  if (r.ok || r.redirected) {
    adminUser.set(username)
    return true
  }
  return false
}

export async function adminFetch(path: string, options: RequestInit = {}): Promise<Response> {
  return fetch(`${API_BASE}${path}`, {
    ...options,
    credentials: 'include',
  })
}
