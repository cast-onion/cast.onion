<script lang="ts">
  export let status: string = ''
</script>

<span class="badge {status}">{status}</span>

<style>
  .badge {
    display: inline-block;
    padding: 2px 7px;
    border-radius: 2px;
    font-size: 11px;
    border: 1px solid;
  }

  .pending {
    background: #1a1500;
    color: var(--warn);
    border-color: var(--warn);
  }
  .approved {
    background: #001a0d;
    color: var(--ok);
    border-color: var(--ok);
  }
  .denied {
    background: #1a0000;
    color: var(--danger);
    border-color: var(--danger);
  }
  .active {
    background: #001a0d;
    color: var(--ok);
    border-color: var(--ok);
  }
  .suspended {
    background: #1a1500;
    color: var(--warn);
    border-color: var(--warn);
  }
  .revoked {
    background: #1a0000;
    color: var(--danger);
    border-color: var(--danger);
  }
</style>
