<script lang="ts">
  import type { Station } from '../types'

  export let station: Station
  export let viewerCount: number = 0
</script>

<div class="card">
  <div class="info">
    <span class="name">{station.DisplayName}</span>
    {#if station.Description}
      <span class="desc">{station.Description}</span>
    {/if}
  </div>
  <div class="meta">
    {#if viewerCount > 0}
      <span class="viewers">👂 {viewerCount}</span>
    {/if}
    {#if station.Genre}
      <span class="genre">{station.Genre}</span>
    {/if}
  </div>
</div>

<style>
  .card {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 20px;
    background: var(--bg2);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    transition: border-color 0.15s;
  }

  .card:hover {
    border-color: #333;
  }

  .info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .name {
    font-size: 14px;
    color: var(--text);
  }
  .desc {
    font-size: 12px;
    color: var(--muted);
  }

  .meta {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .viewers {
    font-size: 12px;
    color: var(--muted);
  }

  .genre {
    font-size: 12px;
    color: var(--accent);
    padding: 2px 8px;
    border: 1px solid var(--accent);
    border-radius: 2px;
    opacity: 0.7;
  }
</style>
