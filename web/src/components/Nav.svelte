<script lang="ts">
  import { session } from '../lib/session'
  import { page } from '../lib/router'
  import type { Page } from '../types'

  function go(p: Page) {
    page.set(p)
  }
</script>

<header>
  <button class="wordmark" on:click={() => go('home')}>cast.onion</button>
  <nav>
    <button class:active={$page === 'directory'} on:click={() => go('directory')}>directory</button>
    <button class:active={$page === 'apply'} on:click={() => go('apply')}>broadcast</button>
  </nav>
  <div class="session-status">
    {#if $session}
      <span class="dot ok">●</span>
      <span class="sid">{$session.slice(0, 8)}</span>
    {:else}
      <span class="dot dim">●</span>
      <span class="dim">connecting...</span>
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
    background: none;
    border: none;
    padding: 0;
    transition: color 0.15s;
  }

  .wordmark:hover {
    color: var(--accent2);
  }

  nav {
    display: flex;
    gap: 24px;
  }

  nav button {
    color: var(--muted);
    font-size: 13px;
    letter-spacing: 0.03em;
    background: none;
    border: none;
    padding: 0;
    transition: color 0.15s;
  }

  nav button:hover,
  nav button.active {
    color: var(--text);
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
  .dim {
    color: var(--muted);
  }
  .sid {
    color: var(--ok);
  }
</style>
