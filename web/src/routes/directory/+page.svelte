<script lang="ts">
  import { onMount } from 'svelte'
  // import Nav from '../../components/Nav.svelte';
  import StationCard from '../../components/StationCard.svelte';
  import Banner from '../../components/Banner.svelte';
  import { session, viewerCounts, apiFetch } from '$lib/session'
  import type { Station } from '../../types'

  let stations: Station[] = []
  let loading = true
  let error = ''

  async function loadStations(): Promise<void> {
    if (!$session) return
    loading = true
    error = ''
    try {
      const r = await apiFetch('/stations')
      if (!r.ok) throw new Error('failed')
      stations = ((await r.json()) as Station[]) ?? []
    } catch {
      error = 'failed to load stations'
    } finally {
      loading = false
    }
  }

  $: if ($session) loadStations()

  onMount(() => {
    if ($session) loadStations()
  })
</script>

<main>
  <div class="page-header">
    <h1>stations</h1>
    <p class="subtext">active stations on the network</p>
  </div>

  <Banner type="error" message={error} />

  {#if loading}
    <div class="empty">tuning in<span class="dots"></span></div>
  {:else if stations.length === 0 && !error}
    <div class="empty">no active stations</div>
  {:else}
    <div class="station-list">
      {#each stations as station (station.ID)}
        <StationCard {station} viewerCount={$viewerCounts[station.ID] ?? 0} />
      {/each}
    </div>
  {/if}
</main>

<style>
  main {
    max-width: 900px;
    margin: 0 auto;
    padding: 48px 32px;
  }

  .page-header {
    margin-bottom: 32px;
  }

  h1 {
    font-size: 22px;
    font-weight: normal;
    color: var(--text);
    margin-bottom: 6px;
  }

  .subtext {
    font-size: 13px;
    color: var(--muted);
  }

  .empty {
    color: var(--muted);
    font-size: 13px;
    padding: 40px 0;
  }

  .dots::after {
    content: '';
    animation: dots 1.2s steps(4, end) infinite;
  }

  @keyframes dots {
    0% {
      content: '';
    }
    25% {
      content: '.';
    }
    50% {
      content: '..';
    }
    75% {
      content: '...';
    }
  }

  .station-list {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
</style>
