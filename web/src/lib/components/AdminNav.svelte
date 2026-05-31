<script lang="ts">
  import { adminUser } from '$lib/admin'
  import { API_BASE } from '$lib/session'
  import { page } from '$lib/router'

  export let toastMessage = ''
  export let toastType: 'ok' | 'deny' = 'ok'

  async function restart(): Promise<void> {
    if (!confirm('restart server?')) return
    await fetch(`${API_BASE}/admin/restart`, { method: 'POST', credentials: 'include' })
  }

  function logout(): void {
    fetch(`${API_BASE}/admin/logout`, { credentials: 'include' })
    adminUser.set(null)
    page.set('admin')
  }
</script>

<header>
  <span class="wordmark">cast.onion / admin</span>
  <div class="right">
    {#if toastMessage}
      <span class="toast {toastType}">{toastMessage}</span>
    {/if}
    <span class="user">{$adminUser}</span>
    <button class="btn-restart" on:click={restart}>restart server</button>
    <button class="btn-logout" on:click={logout}>logout</button>
  </div>
</header>

<style>
  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 28px;
    border-bottom: 1px solid var(--border);
    position: sticky;
    top: 0;
    background: var(--bg);
    z-index: 10;
  }

  .wordmark {
    color: var(--accent);
    font-size: 15px;
    font-weight: bold;
  }

  .right {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .user {
    color: var(--muted);
    font-size: 13px;
  }

  .toast {
    font-size: 12px;
    font-weight: bold;
    padding: 3px 10px;
    border-radius: 3px;
    animation: fadeIn 0.2s ease;
  }

  .ok {
    background: #001a0d;
    color: var(--ok);
    border: 1px solid var(--ok);
  }
  .deny {
    background: #1a0000;
    color: var(--danger);
    border: 1px solid var(--danger);
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(-4px);
    }
    to {
      opacity: 1;
    }
  }

  .btn-restart,
  .btn-logout {
    background: var(--bg3);
    border: 1px solid var(--border);
    color: var(--muted);
    font-size: 12px;
    padding: 5px 12px;
    transition:
      color 0.15s,
      border-color 0.15s;
  }

  .btn-restart:hover {
    color: var(--danger);
    border-color: var(--danger);
  }
  .btn-logout:hover {
    color: var(--text);
  }
</style>
