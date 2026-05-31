<script lang="ts">
  import '../app.css'
  import { page } from '$app/stores'

  let session: string | null = null

  setTimeout(() => {
    session = Math.random().toString(36).slice(2)
  }, 1000)
</script>

<header>
  <a href="/" class="wordmark">cast.onion</a>

  <nav>
    <a href="/directory" class:active={$page.url.pathname === '/directory'}>
      directory
    </a>
    <a href="/apply" class:active={$page.url.pathname === '/apply'}>
      broadcast
    </a>
    <a href="/support" class:active={$page.url.pathname === '/support'}>
      tip
    </a>
  </nav>

  <div class="session-status">
    {#if session}
      <span class="dot ok">●</span>
      <span class="sid">{session.slice(0, 8)}</span>
    {:else}
      <span class="dot muted">●</span>
      <span class="muted-text">connecting...</span>
    {/if}
  </div>
</header>

<slot />