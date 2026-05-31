<script lang="ts">
  export let name = ''
  export let stationKey = ''
  export let accessToken = ''
  export let onClose: () => void = () => {}
</script>

<div class="overlay">
  <div class="modal">
    <div class="title">station approved</div>
    <div class="subtitle">save these credentials — they will not be shown again</div>

    <div class="station-name">{name}</div>

    <div class="cred-label">station key</div>
    <div class="cred-value">{stationKey}</div>

    <div class="cred-label">access token</div>
    <div class="cred-value">{accessToken}</div>

    <button class="btn" on:click={onClose}>i have saved these</button>
  </div>
</div>

<style>
  .overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.75);
    z-index: 100;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .modal {
    background: var(--bg2);
    border: 1px solid #333;
    border-radius: 6px;
    padding: 32px;
    max-width: 520px;
    width: 90%;
    font-family: var(--font);
  }

  .title {
    color: var(--accent);
    font-size: 15px;
    font-weight: bold;
    margin-bottom: 4px;
  }
  .subtitle {
    color: var(--muted);
    font-size: 12px;
    margin-bottom: 24px;
  }
  .station-name {
    color: var(--text);
    font-size: 14px;
    margin-bottom: 20px;
  }

  .cred-label {
    color: var(--muted);
    font-size: 11px;
    margin-bottom: 6px;
  }

  .cred-value {
    color: var(--text);
    font-size: 11px;
    word-break: break-all;
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 3px;
    padding: 10px 12px;
    margin-bottom: 16px;
  }

  .btn {
    background: var(--accent);
    color: #000;
    font-size: 12px;
    font-weight: bold;
    padding: 10px 20px;
    margin-top: 8px;
    transition: background 0.15s;
  }

  .btn:hover {
    background: var(--accent2);
  }
</style>
