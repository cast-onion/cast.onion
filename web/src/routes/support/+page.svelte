<script>
  const donateLink = "https://streamelements.com/itsbr0dyy/tip"
</script>

<main>
  <div class="hero">
    <div class="logo">support cast.onion</div>

    <p class="tagline">
      keep the network online.
    </p>

    <p class="desc">
      cast.onion is fully independent — no ads, no tracking, no investors.<br />
      if you enjoy the platform, consider supporting it with a one-time contribution.
    </p>

    <div class="actions">
      <a href={donateLink} target="_blank" class="btn-primary">
        donate
      </a>
    </div>

    <div class="install-section">
      <p class="install-label">what you're supporting</p>

      <div class="install-box">
        <span class="install-desc">
          infrastructure, bandwidth, uptime, and continued development of the platform.
        </span>
      </div>
    </div>
  </div>

  <div class="features">
    <div class="feature">
      <span class="feature-title">no ads</span>
      <span class="feature-desc">
        your support keeps the client experience clean and uninterrupted.
      </span>
    </div>

    <div class="feature">
      <span class="feature-title">independent</span>
      <span class="feature-desc">
        no outside funding — fully supported by users.
      </span>
    </div>

    <div class="feature">
      <span class="feature-title">direct impact</span>
      <span class="feature-desc">
        every contribution goes straight into running the network.
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
    margin-bottom: 20px;
  }

  .btn-primary {
    display: inline-block;
    background: var(--accent);
    color: #000;
    font-size: 13px;
    font-weight: bold;
    padding: 12px 28px;
    border-radius: var(--radius);
    letter-spacing: 0.05em;
    text-decoration: none;
    transition: background 0.15s;
  }

  .btn-primary:hover {
    background: var(--accent2);
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
    background: var(--bg2);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 16px 20px;
  }

  .install-desc {
    font-size: 13px;
    color: var(--muted);
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