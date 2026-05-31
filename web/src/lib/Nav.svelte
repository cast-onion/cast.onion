<script>
  import { session } from './session.js'
  import { page } from './router.js'

  export let current = ''
</script>

<header>
  <a href="/" class="wordmark" on:click|preventDefault={() => page.set('home')}>cast.onion</a>
  <nav>
    <a
      href="/directory"
      class:active={current === 'directory'}
      on:click|preventDefault={() => page.set('directory')}>directory</a
    >
    <a
      href="/apply"
      class:active={current === 'apply'}
      on:click|preventDefault={() => page.set('apply')}>broadcast</a
    >
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
  }

  .wordmark:hover {
    text-decoration: none;
    color: var(--accent2);
  }

  nav {
    display: flex;
    gap: 24px;
  }

  nav a {
    color: var(--muted);
    font-size: 13px;
    letter-spacing: 0.03em;
    transition: color 0.15s;
  }

  nav a:hover,
  nav a.active {
    color: var(--text);
    text-decoration: none;
  }

  .session-status {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
  }

  .dot {
    font-size: 10px;
  }
  .ok {
    color: var(--ok);
  }
  .muted {
    color: var(--muted);
  }
  .sid {
    color: var(--ok);
  }
  .muted-text {
    color: var(--muted);
  }
</style>
