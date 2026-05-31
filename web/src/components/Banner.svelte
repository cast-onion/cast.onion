<script lang="ts">
  export let type: 'error' | 'success' | 'info' = 'info'
  export let message: string = ''
</script>

{#if message}
  <div class="banner {type}">{message}</div>
{/if}

<style>
  .banner {
    padding: 12px 16px;
    border-radius: var(--radius);
    margin-bottom: 20px;
    font-size: 13px;
    border: 1px solid var(--border);
  }

  .error {
    border-color: var(--danger);
    color: var(--danger);
    background: #1a0000;
  }
  .success {
    border-color: var(--ok);
    color: var(--ok);
    background: #001a0d;
  }
  .info {
    border-color: var(--border);
    color: var(--muted);
    background: var(--bg2);
  }
</style>
