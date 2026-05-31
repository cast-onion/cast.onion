import { writable } from 'svelte/store'
import type { Page } from '../types'

function getInitialPage(): Page {
  if (typeof window === 'undefined') return 'home'
  const path = window.location.pathname
  if (path.startsWith('/directory')) return 'directory'
  if (path.startsWith('/apply')) return 'apply'
  if (path.startsWith('/admin/dashboard')) return 'admin-dashboard'
  if (path.startsWith('/admin')) return 'admin'
  return 'home'
}

export const page = writable<Page>(getInitialPage())

const pathMap: Record<Page, string> = {
  home: '/',
  directory: '/directory',
  apply: '/apply',
  admin: '/admin',
  'admin-dashboard': '/admin/dashboard',
}

page.subscribe(p => {
  if (typeof window === 'undefined') return
  window.history.pushState({}, '', pathMap[p] ?? '/')
})