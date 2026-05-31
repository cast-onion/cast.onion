<script lang="ts">
  import Nav from '../../components/Nav.svelte'
  import Banner from '../../components/Banner.svelte'
  import { session, apiFetch } from '../../lib/session'

  let form = { station_name: '', contact_email: '', genre: '', description: '', notes: '' }
  let submitting = false
  let submitted = false
  let submittedID = ''
  let error = ''

  async function submit(): Promise<void> {
    if (!form.station_name || !form.contact_email || !form.description) {
      error = 'please fill in all required fields'
      return
    }

    if (!$session) {
      await new Promise<void>(resolve => {
        const unsub = session.subscribe(s => {
          if (s) { unsub(); resolve(); }
        })
        setTimeout(() => { unsub(); resolve(); }, 5000)
      })
    }

    if (!$session) {
      error = 'could not establish a session — check your connection and refresh'
      return
    }

    submitting = true
    error = ''

    try {
      const r = await apiFetch('/apply', { method: 'POST', body: JSON.stringify(form) })
      if (!r.ok) {
        const text = await r.text()
        throw new Error(text || 'server error')
      }
      const data = await r.json()
      submittedID = (data.ID ?? data.id ?? '') as string
      submitted = true
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'submission failed — please try again'
    } finally {
      submitting = false
    }
  }
</script>

<main>
  <div class="page-header">
    <h1>apply to broadcast</h1>
    <p class="subtext">applications are reviewed manually. approved hosts receive a station key and access token.</p>
  </div>

  {#if submitted}
    <div class="success-box">
      <div class="success-title">application submitted</div>
      {#if submittedID}<div class="success-id">id: {submittedID}</div>{/if}
      <p class="success-desc">your application is under review. you'll receive a follow-up email within 48 hours.</p>
    </div>
  {:else}
    <Banner type="error" message={error} />

    {#if !$session}
      <Banner type="info" message="connecting to server..." />
    {/if}

    <form class="form" on:submit|preventDefault={submit}>
      <div class="field">
        <label for="station_name">station name <span class="req">*</span></label>
        <input id="station_name" type="text" bind:value={form.station_name} placeholder="e.g. late night static" required />
      </div>
      <div class="field">
        <label for="contact_email">contact email <span class="req">*</span></label>
        <input id="contact_email" type="email" bind:value={form.contact_email} placeholder="you@somewhere.net" required />
      </div>
      <div class="field">
        <label for="genre">genre</label>
        <input id="genre" type="text" bind:value={form.genre} placeholder="e.g. ambient, jazz, noise" />
      </div>
      <div class="field">
        <label for="description">what will you broadcast? <span class="req">*</span></label>
        <textarea id="description" bind:value={form.description} rows={4} placeholder="tell us about your station" required></textarea>
      </div>
      <div class="field">
        <label for="notes">anything else?</label>
        <textarea id="notes" bind:value={form.notes} rows={3} placeholder="optional"></textarea>
      </div>
      <button type="submit" class="btn-submit" disabled={submitting}>
        {submitting ? 'submitting...' : 'submit application'}
      </button>
    </form>
  {/if}
</main>

<style>
  main { max-width: 900px; margin: 0 auto; padding: 48px 32px; }
  .page-header { margin-bottom: 32px; }
  h1 { font-size: 22px; font-weight: normal; color: var(--text); margin-bottom: 6px; }
  .subtext { font-size: 13px; color: var(--muted); }
  .form { display: flex; flex-direction: column; gap: 22px; max-width: 560px; }
  .field { display: flex; flex-direction: column; gap: 8px; }
  label { font-size: 12px; color: var(--muted); letter-spacing: 0.03em; }
  .req { color: var(--accent); }
  textarea { resize: vertical; }
  .btn-submit {
    background: var(--accent); color: #000; font-size: 13px; font-weight: bold;
    padding: 11px 24px; align-self: flex-start; letter-spacing: 0.03em; transition: background 0.15s;
  }
  .btn-submit:hover:not(:disabled) { background: var(--accent2); }
  .btn-submit:disabled { opacity: 0.4; cursor: not-allowed; }
  .success-box {
    background: var(--bg2); border: 1px solid var(--ok);
    border-radius: var(--radius); padding: 24px; max-width: 480px;
  }
  .success-title { color: var(--ok); font-size: 15px; margin-bottom: 8px; }
  .success-id { font-size: 11px; color: var(--muted); margin-bottom: 12px; }
  .success-desc { font-size: 13px; color: var(--muted); line-height: 1.6; }
</style>