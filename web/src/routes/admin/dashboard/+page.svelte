<script lang="ts">
  import { onMount } from 'svelte'
  import AdminNav from '$lib/components/AdminNav.svelte'
  import StatusBadge from '$lib/components/StatusBadge.svelte'
  import CredentialsModal from '$lib/components/CredentialsModal.svelte'
  import { adminUser } from '$lib/admin'
  import { API_BASE } from '$lib/session'
  import { page } from '$lib/router'
  import type { Application, Station, AdminAction, ApprovalResult } from '../../../types'

  let applications: Application[] = []
  let stations: Station[] = []
  let actions: AdminAction[] = []
  let loading = true

  let toastMessage = ''
  let toastType: 'ok' | 'deny' = 'ok'
  let toastTimer: ReturnType<typeof setTimeout>

  let creds: { name: string; stationKey: string; accessToken: string } | null = null

  if (!$adminUser) {
    page.set('admin')
  }

  async function load(): Promise<void> {
    loading = true
    try {
      const [appsR, stationsR, actionsR] = await Promise.all([
        fetch(`${API_BASE}/admin/dashboard`, { credentials: 'include' }),
        fetch(`${API_BASE}/admin/dashboard`, { credentials: 'include' }),
        fetch(`${API_BASE}/admin/dashboard`, { credentials: 'include' }),
      ])
      const dash = await appsR.json()
      applications = dash.Applications ?? []
      stations = dash.Stations ?? []
      actions = dash.Actions ?? []
    } catch {
      applications = []
    } finally {
      loading = false
    }
  }

  function toast(msg: string, type: 'ok' | 'deny'): void {
    toastMessage = msg
    toastType = type
    clearTimeout(toastTimer)
    toastTimer = setTimeout(() => {
      toastMessage = ''
    }, 4000)
  }

  async function approve(app: Application): Promise<void> {
    const r = await fetch(`${API_BASE}/admin/approve/${app.ID}`, {
      method: 'POST',
      credentials: 'include',
    })
    if (!r.ok) {
      toast('approve failed', 'deny')
      return
    }
    const data = (await r.json()) as ApprovalResult
    applications = applications.filter(a => a.ID !== app.ID)
    toast(`approved: ${app.StationName}`, 'ok')
    creds = {
      name: app.StationName,
      stationKey: data.station_key,
      accessToken: data.access_token,
    }
  }

  async function deny(app: Application): Promise<void> {
    const r = await fetch(`${API_BASE}/admin/deny/${app.ID}`, {
      method: 'POST',
      credentials: 'include',
    })
    if (!r.ok) {
      toast('deny failed', 'deny')
      return
    }
    applications = applications.filter(a => a.ID !== app.ID)
    toast(`denied: ${app.StationName}`, 'deny')
  }

  async function suspend(station: Station): Promise<void> {
    await fetch(`${API_BASE}/admin/suspend/${station.ID}`, {
      method: 'POST',
      credentials: 'include',
    })
    stations = stations.map(s => (s.ID === station.ID ? { ...s, Status: 'suspended' } : s))
    toast(`suspended: ${station.DisplayName}`, 'deny')
  }

  async function revoke(station: Station): Promise<void> {
    if (!confirm(`revoke ${station.DisplayName}? this cannot be undone.`)) return
    await fetch(`${API_BASE}/admin/revoke/${station.ID}`, {
      method: 'POST',
      credentials: 'include',
    })
    stations = stations.map(s => (s.ID === station.ID ? { ...s, Status: 'revoked' } : s))
    toast(`revoked: ${station.DisplayName}`, 'deny')
  }

  async function unsuspend(station: Station): Promise<void> {
    await fetch(`${API_BASE}/admin/unsuspend/${station.ID}`, {
      method: 'POST',
      credentials: 'include',
    })
    stations = stations.map(s => (s.ID === station.ID ? { ...s, Status: 'active' } : s))
    toast(`unsuspended: ${station.DisplayName}`, 'ok')
  }

  function timeAgo(ts: string): string {
    const d = Date.now() - new Date(ts).getTime()
    if (d < 60000) return 'just now'
    if (d < 3600000) return `${Math.floor(d / 60000)}m ago`
    if (d < 86400000) return `${Math.floor(d / 3600000)}h ago`
    return new Date(ts).toLocaleDateString()
  }

  onMount(load)
</script>

<AdminNav {toastMessage} {toastType} />

{#if creds}
  <CredentialsModal
    name={creds.name}
    stationKey={creds.stationKey}
    accessToken={creds.accessToken}
    onClose={() => (creds = null)}
  />
{/if}

{#if loading}
  <div class="loading">loading dashboard...</div>
{:else}
  <div class="layout">
    <div class="panel">
      <div class="panel-title">applications</div>
      {#if applications.length === 0}
        <div class="empty">no applications yet</div>
      {:else}
        <table>
          <thead>
            <tr><th>station</th><th>contact</th><th>status</th><th>actions</th></tr>
          </thead>
          <tbody>
            {#each applications as app (app.ID)}
              <tr>
                <td>
                  <div class="cell-main">{app.StationName}</div>
                  {#if app.Genre}<div class="cell-sub">{app.Genre}</div>{/if}
                </td>
                <td class="muted">{app.ContactEmail}</td>
                <td><StatusBadge status={app.Status} /></td>
                <td class="actions-cell">
                  {#if app.Status === 'pending'}
                    <button class="btn-ok" on:click={() => approve(app)}>approve</button>
                    <button class="btn-danger" on:click={() => deny(app)}>deny</button>
                  {:else}
                    <span class="reviewed">reviewed</span>
                  {/if}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/if}
    </div>

    <div class="panel">
      <div class="panel-title">stations</div>
      {#if stations.length === 0}
        <div class="empty">no stations yet</div>
      {:else}
        <table>
          <thead>
            <tr><th>name</th><th>status</th><th>actions</th></tr>
          </thead>
          <tbody>
            {#each stations as station (station.ID)}
              <tr>
                <td class="cell-main">{station.DisplayName}</td>
                <td><StatusBadge status={station.Status} /></td>
                <td class="actions-cell">
                  {#if station.Status === 'active'}
                    <button class="btn-warn" on:click={() => suspend(station)}>suspend</button>
                  {/if}
                  {#if station.Status === 'suspended'}
                    <button class="btn-ok" on:click={() => unsuspend(station)}>unsuspend</button>
                  {/if}
                  {#if station.Status !== 'revoked'}
                    <button class="btn-danger" on:click={() => revoke(station)}>revoke</button>
                  {/if}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/if}
    </div>

    <div class="panel panel-full">
      <div class="panel-title">audit log</div>
      {#if actions.length === 0}
        <div class="empty">no actions logged yet</div>
      {:else}
        <table>
          <thead>
            <tr><th>time</th><th>admin</th><th>action</th><th>target</th><th>reason</th></tr>
          </thead>
          <tbody>
            {#each actions as action (action.ID)}
              <tr class="log-row">
                <td class="muted">{timeAgo(action.CreatedAt)}</td>
                <td class="muted">{action.AdminID}</td>
                <td>{action.Action}</td>
                <td class="mono muted">{action.TargetType} / {action.TargetID.slice(0, 8)}...</td>
                <td class="muted">{action.Reason}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/if}
    </div>
  </div>
{/if}

<style>
  .loading {
    color: var(--muted);
    font-size: 13px;
    padding: 48px 32px;
  }

  .layout {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1px;
    background: var(--border);
    min-height: calc(100vh - 53px);
  }

  .panel {
    background: var(--bg);
    padding: 24px 28px;
    overflow-y: auto;
  }

  .panel-full {
    grid-column: 1 / -1;
  }

  .panel-title {
    font-size: 11px;
    color: var(--muted);
    letter-spacing: 0.08em;
    text-transform: uppercase;
    margin-bottom: 16px;
  }

  .empty {
    color: var(--muted);
    font-size: 12px;
    padding: 16px 0;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th {
    text-align: left;
    font-size: 11px;
    color: var(--muted);
    padding: 0 0 8px 0;
    border-bottom: 1px solid var(--border);
    font-weight: normal;
  }

  td {
    padding: 10px 0;
    border-bottom: 1px solid var(--bg3);
    vertical-align: top;
    font-size: 13px;
    padding-right: 12px;
  }

  .cell-main {
    color: var(--text);
  }
  .cell-sub {
    font-size: 11px;
    color: var(--muted);
    margin-top: 2px;
  }
  .muted {
    color: var(--muted);
  }
  .mono {
    font-family: var(--font);
    font-size: 11px;
  }
  .reviewed {
    font-size: 11px;
    color: var(--muted);
  }

  .actions-cell {
    display: flex;
    gap: 6px;
    justify-content: flex-end;
    align-items: center;
  }

  .btn-ok,
  .btn-danger,
  .btn-warn {
    padding: 4px 10px;
    font-size: 11px;
    font-weight: bold;
    border-radius: 3px;
    transition: opacity 0.15s;
  }

  .btn-ok:hover,
  .btn-danger:hover,
  .btn-warn:hover {
    opacity: 0.85;
  }

  .btn-ok {
    background: var(--ok);
    color: #000;
  }
  .btn-danger {
    background: var(--danger);
    color: #fff;
  }
  .btn-warn {
    background: var(--warn);
    color: #000;
  }

  .log-row td {
    font-size: 12px;
  }
</style>
