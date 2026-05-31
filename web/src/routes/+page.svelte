<script>
  import { goto } from '$app/navigation'
</script>

<main>
  <div class="hero">
    <div class="logo">cast.onion</div>
    <p class="tagline">a private, application-based internet radio network.</p>
    <p class="desc">
      no public access. no algorithms. stations are hand-approved.<br />
      connect, listen, and broadcast — entirely on your terms.
    </p>

    <div class="actions">
      <button class="btn-primary" on:click={() => goto('/directory')}>
        browse stations
      </button>

      <button class="btn-secondary" on:click={() => goto('/apply')}>
        apply to broadcast
      </button>
    </div>

    <div class="install-section">
      <p class="install-label">get the app</p>
      <div class="install-box">
        <span class="install-desc">
          download the cast.onion desktop app to broadcast and listen natively.
        </span>
        <a href="http://localhost:1050/app/cast-onion-app" class="btn-install" download>
          download for linux
        </a>
      </div>
    </div>
  </div>

  <div class="features">
    <div class="feature">
      <span class="feature-icon">
        <span class="icon-box">
          <img src="src/assets/key.webp" alt="" class="icon" />
        </span>
      </span>
      <span class="feature-title">private by design</span>
      <span class="feature-desc">
        every connection requires a valid session. nothing is publicly accessible.
      </span>
    </div>

    <div class="feature">
      <span class="feature-icon">
        <span class="icon-box">
          <img src="src/assets/broadcast.webp" alt="" class="icon" />
        </span>
      </span>
      <span class="feature-title">broadcast anything</span>
      <span class="feature-desc">
        stream from your mic or any audio source. opus codec, 48kHz stereo.
      </span>
    </div>

    <div class="feature">
      <span class="feature-icon">
        <span class="icon-box">
          <img src="src/assets/thumbsup.webp" alt="" class="icon" />
        </span>
      </span>
      <span class="feature-title">guest rooms</span>
      <span class="feature-desc">
        invite up to 4 guests to join your broadcast with a single link.
      </span>
    </div>
  </div>
</main>

<style>
  main {
    max-width: 860px;
    margin: 0 auto;
    padding: 80px 32px 48px;
  }

  .hero {
    margin-bottom: 72px;
  }

  .logo {
    font-size: 36px;
    color: var(--accent);
    font-weight: bold;
    letter-spacing: 0.04em;
    margin-bottom: 16px;
  }

  .tagline {
    font-size: 16px;
    color: var(--text);
    margin-bottom: 12px;
  }

  .desc {
    font-size: 13px;
    color: var(--muted);
    line-height: 1.8;
    margin-bottom: 36px;
  }

  .actions {
    display: flex;
    gap: 12px;
    margin-bottom: 48px;
  }

  .btn-primary {
    background: var(--accent);
    color: #000;
    font-size: 13px;
    font-weight: bold;
    padding: 11px 24px;
    border-radius: var(--radius);
    letter-spacing: 0.03em;
    transition: background 0.15s;
  }

  .btn-primary:hover {
    background: var(--accent2);
  }

  .btn-secondary {
    background: var(--bg2);
    color: var(--muted);
    font-size: 13px;
    padding: 11px 24px;
    border: 1px solid var(--border);
    border-radius: var(--radius);
    transition: color 0.15s, border-color 0.15s;
  }

  .btn-secondary:hover {
    color: var(--text);
    border-color: #444;
  }

  .install-section {
    border-top: 1px solid var(--border);
    padding-top: 28px;
  }

  .install-label {
    font-size: 11px;
    color: var(--muted);
    letter-spacing: 0.08em;
    text-transform: uppercase;
    margin-bottom: 12px;
  }

  .install-box {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: var(--bg2);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 16px 20px;
  }

  .install-desc {
    font-size: 13px;
    color: var(--muted);
  }

  .btn-install {
    background: var(--bg3);
    border: 1px solid var(--border);
    color: var(--text);
    font-size: 12px;
    font-family: var(--font);
    padding: 8px 18px;
    border-radius: var(--radius);
    white-space: nowrap;
    transition: border-color 0.15s, color 0.15s;
    text-decoration: none;
    display: inline-block;
  }

  .btn-install:hover {
    border-color: var(--accent);
    color: var(--accent);
  }

  .features {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 2px;
  }

  .feature {
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding: 20px;
    background: var(--bg2);
    border: 1px solid var(--border);
    border-radius: var(--radius);
  }

  .feature-title {
    font-size: 13px;
    color: var(--text);
  }

  .feature-desc {
    font-size: 12px;
    color: var(--muted);
    line-height: 1.6;
  }
</style>