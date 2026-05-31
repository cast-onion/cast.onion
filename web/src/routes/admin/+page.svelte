<script lang="ts">
  import { goto } from '$app/navigation'
  import { adminUser } from '$lib/admin'
  import { API_BASE } from '$lib/session'

  let username = ''
  let password = ''
  let error = ''
  let loading = false

  async function login(): Promise<void> {
    if (!username || !password) {
      error = 'enter credentials'
      return
    }

    loading = true
    error = ''

    try {
      const r = await fetch(`${API_BASE}/admin/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: new URLSearchParams({ username, password }),
        credentials: 'include',
        redirect: 'manual',
      })

      if (r.ok || r.status === 302 || r.type === 'opaqueredirect') {
        adminUser.set(username)
        goto('/admin/dashboard')
      } else {
        const text = await r.text()
        error = text.includes('invalid') ? 'invalid credentials' : 'login failed'
      }
    } catch {
      error = 'could not reach server'
    } finally {
      loading = false
    }
  }
</script>

<div class="login-wrap">
  <div class="box">
    <div class="wordmark">cast.onion / admin</div>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    <form on:submit|preventDefault={login}>
      <div class="field">
        <label for="username">username</label>
        <input id="username" type="text" bind:value={username} autocomplete="username" />
      </div>
      <div class="field">
        <label for="password">password</label>
        <input
          id="password"
          type="password"
          bind:value={password}
          autocomplete="current-password"
        />
      </div>
      <button type="submit" class="btn" disabled={loading}>
        {loading ? 'signing in...' : 'sign in'}
      </button>
    </form>
  </div>
</div>

<style>
  .login-wrap {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
  }

  .box {
    width: 320px;
  }

  .wordmark {
    color: var(--accent);
    font-size: 16px;
    font-weight: bold;
    display: block;
    margin-bottom: 32px;
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 16px;
  }

  label {
    font-size: 12px;
    color: var(--muted);
  }

  .error {
    color: var(--danger);
    font-size: 12px;
    margin-bottom: 16px;
    padding: 8px 12px;
    border: 1px solid var(--danger);
    border-radius: var(--radius);
    background: #1a0000;
  }

  .btn {
    background: var(--accent);
    color: #000;
    font-size: 13px;
    font-weight: bold;
    padding: 10px 20px;
    width: 100%;
    margin-top: 8px;
    transition: background 0.15s;
  }

  .btn:hover:not(:disabled) {
    background: var(--accent2);
  }
  .btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
</style>
