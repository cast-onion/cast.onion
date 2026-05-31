<script lang="ts">
  const pages = [
    { name: "Home", href: "/" },
    { name: "Client", href: "/client" },
    { name: "Codebase", href: "/codebase" },
    { name: "API", href: "/api" },
    { name: "Installers", href: "/installers" }
  ];
</script>

<div class="layout">
  <header class="header">
    <div class="logo">
      <img class="logo-icon" src="/icon.png" alt="logo" />
      <span>cast.onion</span>
    </div>
  </header>

  <div class="main">
    <aside class="sidebar">
      {#each pages as page}
        <a class="link" href={page.href}>{page.name}</a>
      {/each}
    </aside>

    <section class="content">
      <slot />
    </section>
  </div>

  <footer class="footer">
    &copy; 2026 cast.onion
  </footer>
</div>

<style>
  .layout {
    display: grid;
    grid-template-rows: 50px 1fr 40px;
    height: 100vh;
    background: #0f0f10;
    color: #eaeaea;
    font-family: system-ui, sans-serif;
  }

  .header {
    display: flex;
    align-items: center;
    padding: 0 16px;
    border-bottom: 1px solid #2a2a2a;
    background: #0f0f10;
  }

  .logo {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
  }

  .logo-icon {
    width: 40px;
    height: 40px;
	transform: translateY(2px);
  }

  .main {
    display: grid;
    grid-template-columns: 240px 1fr;
    min-height: 0;
  }

  .sidebar {
    border-right: 1px solid #2a2a2a;
    padding: 12px;
    display: flex;
    flex-direction: column;
    gap: 6px;
    background: #141414;
  }

  .link {
    padding: 6px 8px;
    border-radius: 6px;
    text-decoration: none;
    color: #eaeaea;
  }

  .link:hover {
    background: #222;
  }

  .content {
    padding: 24px;
    overflow: auto;
  }

  .footer {
    border-top: 1px solid #2a2a2a;
    display: flex;
    align-items: center;
    padding: 0 16px;
    font-size: 12px;
    color: #aaa;
    background: #0f0f10;
  }

	:global(html, body) {
	margin: 0;
	padding: 0;
	height: 100%;
	background: #0f0f10;
	}
</style>