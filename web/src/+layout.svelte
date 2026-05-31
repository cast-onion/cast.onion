<script lang="ts">
  import { onMount } from 'svelte'
  import { page } from '$app/stores'
  import { session, connectSession } from '../lib/session'
  import '../app.css'

  onMount(() => {
    connectSession()
  })
</script>

<header>
  <a href="/" class="wordmark">cast.onion</a>

  <nav>
    <a href="/directory" class:active={$page.url.pathname === '/directory'}>directory</a>
    <a href="/apply" class:active={$page.url.pathname === '/apply'}>broadcast</a>
  </nav>

  <div class="session-status">
    {#if $session}
      <span class="dot ok">●</span>
      <span class="sid">{$session.slice(0, 8)}</span>
    {:else}
      <span class="dot muted">●</span>
      <span class="muted-text">connecting...</span>
    {/if}
  </div>
</header>

<slot />

<style>
  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 20px 32px;
    border-bottom: 1px solid var(--border);
  }
  .wordmark {
    font-size: 16px;
    font-weight: bold;
    color: var(--accent);
    letter-spacing: 0.05em;
    text-decoration: none;
  }
  .wordmark:hover { color: var(--accent2); }
  nav { display: flex; gap: 24px; }
  nav a {
    color: var(--muted);
    font-size: 13px;
    text-decoration: none;
    transition: color 0.15s;
  }
  nav a:hover, nav a.active { color: var(--text); }
  .session-status { display: flex; align-items: center; gap: 6px; font-size: 12px; }
  .dot { font-size: 10px; }
  .ok { color: var(--ok); }
  .muted { color: var(--muted); }
  .sid { color: var(--ok); }
  .muted-text { color: var(--muted); }
</style>