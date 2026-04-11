<script lang="ts">
  import '../app.css'
  import { page } from '$app/stores'
  import { KeyRound, Moon, Sun, LayoutList, ScrollText, Settings, ShieldPlus } from '@lucide/svelte'

  let theme = $state<'karden' | 'karden-dark'>('karden')

  function toggleTheme() {
    theme = theme === 'karden' ? 'karden-dark' : 'karden'
    document.documentElement.setAttribute('data-theme', theme)
  }

  const navItems = [
    { href: '/',        label: 'Workloads', icon: LayoutList },
    { href: '/secrets', label: 'Secrets',   icon: ShieldPlus },
    { href: '/audit',   label: 'Audit Log', icon: ScrollText },
    { href: '/settings',label: 'Settings',  icon: Settings },
  ]

  let { children } = $props()
</script>

<div class="min-h-screen bg-base-200 flex flex-col">
  <!-- Navbar -->
  <div class="navbar bg-base-100 shadow-sm px-6 sticky top-0 z-50">
    <div class="flex-1 flex items-center gap-3">
      <KeyRound class="text-primary" size={18} />
      <span class="font-bold tracking-tight">karden</span>
      <span class="text-xs text-base-content/40 hidden sm:block">secret lifecycle manager</span>
    </div>

    <!-- Nav links -->
    <div class="flex-none hidden md:flex">
      <ul class="menu menu-horizontal px-1 text-sm">
        {#each navItems as item}
          <li>
            <a
              href={item.href}
              class:active={$page.url.pathname === item.href}
              class="gap-2"
            >
              <svelte:component this={item.icon} size={15} />
              {item.label}
            </a>
          </li>
        {/each}
      </ul>
    </div>

    <div class="flex-none ml-2">
      <button class="btn btn-ghost btn-sm btn-circle" onclick={toggleTheme}>
        {#if theme === 'karden'}
          <Moon size={15} />
        {:else}
          <Sun size={15} />
        {/if}
      </button>
    </div>
  </div>

  <!-- Page content -->
  <main class="flex-1 p-6 max-w-6xl w-full mx-auto">
    {@render children()}
  </main>
</div>
