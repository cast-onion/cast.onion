<script lang="ts">
  import { page } from '$app/stores'
  import { goto } from '$app/navigation'
  import { adminUser } from '$lib/admin'

  $: pathname = $page.url.pathname

  $: isLogin = pathname === '/admin'
  $: isDashboard = pathname.startsWith('/admin/dashboard')

  $: if (isDashboard && !$adminUser) {
    goto('/admin')
  }
</script>

<div class="admin-shell">
  <slot />
</div>

<style>
  .admin-shell {
    min-height: 100vh;
    background: var(--bg);
  }
</style>
