<script lang="ts">
  import { page } from '$app/stores';
  import { ScrollText, Settings } from '@lucide/svelte';
  import '../app.css';

  let { children } = $props()
</script>

<div class="min-h-screen flex flex-col bg-base-200">

  <!-- Topbar -->
  <header class="bg-base-100 border-b border-base-300 sticky top-0 z-50">
    <div class="max-w-6xl mx-auto px-4 h-10 flex items-center gap-4">

      <!-- Logo -->
      <a href="/" class="flex items-center gap-2 shrink-0">
        <!-- <KeyRound class="text-primary" size={17} /> -->
        <span class="font-bold tracking-tight text-sm">karden</span>
      </a>

      <!-- Nav -->
      <nav class="flex items-center gap-1 flex-1">
        <a
          href="/"
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm transition-colors
            {$page.url.pathname === '/'
              ? 'bg-primary/10 text-primary font-medium'
              : 'text-base-content/70 hover:bg-base-200'}"
        >
          Secrets
        </a>

        <a
          href="/audit"
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm transition-colors
            {$page.url.pathname === '/audit'
              ? 'bg-primary/10 text-primary font-medium'
              : 'text-base-content/70 hover:bg-base-200'}"
        >
          <ScrollText size={14} />
          Audit Log
        </a>
      </nav>

      <!-- Right: settings -->
      <div class="flex items-center gap-1">
        <a
          href="/settings"
          class="btn btn-ghost btn-sm btn-circle {$page.url.pathname === '/settings' ? 'text-primary' : ''}"
          aria-label="Settings"
        >
          <Settings size={15} />
        </a>
      </div>

    </div>
  </header>

  <!-- Page content -->
  <main class="flex-1 p-4 max-w-6xl w-full mx-auto">
    {@render children()}
  </main>

</div>
