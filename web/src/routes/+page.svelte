<script lang="ts">
  import { ShieldCheck, RefreshCw, AlertTriangle } from '@lucide/svelte'

  const workloads = [
    {
      id: 1,
      pod_name: 'torchi-db-0',
      namespace: 'default',
      secret_name: 'torchi-db-secret',
      type: 'database',
      db_type: 'postgres',
      status: 'active',
      rotation_days: 30,
      last_rotated_at: '2026-03-12T00:00:00Z',
    },
    {
      id: 2,
      pod_name: 'redis-0',
      namespace: 'default',
      secret_name: 'redis-secret',
      type: 'database',
      db_type: 'redis',
      status: 'active',
      rotation_days: 14,
      last_rotated_at: null,
    },
    {
      id: 3,
      pod_name: 'payment-api',
      namespace: 'production',
      secret_name: 'stripe-secret',
      type: 'manual',
      db_type: null,
      status: 'active',
      rotation_days: 90,
      last_rotated_at: '2025-12-01T00:00:00Z',
    },
  ]

  function daysUntilRotation(lastRotatedAt: string | null, rotationDays: number): number | null {
    if (!lastRotatedAt) return null
    const next = new Date(lastRotatedAt).getTime() + rotationDays * 86400_000
    return Math.ceil((next - Date.now()) / 86400_000)
  }

  function rotationStatus(lastRotatedAt: string | null, rotationDays: number) {
    const days = daysUntilRotation(lastRotatedAt, rotationDays)
    if (days === null) return 'never'
    if (days < 0) return 'overdue'
    if (days <= 7) return 'soon'
    return 'ok'
  }

  function formatDate(dateStr: string | null) {
    if (!dateStr) return '—'
    return new Date(dateStr).toLocaleDateString('en-US', {
      year: 'numeric', month: 'short', day: 'numeric',
    })
  }

  const pendingRotation = workloads.filter(w => {
    const s = rotationStatus(w.last_rotated_at, w.rotation_days)
    return s === 'overdue' || s === 'soon' || s === 'never'
  }).length
</script>

<!-- Stats -->
<div class="grid grid-cols-3 gap-4 mb-6">
  <div class="stat bg-base-100 rounded-box shadow-sm">
    <div class="stat-figure text-primary"><ShieldCheck size={28} /></div>
    <div class="stat-title">Managed Workloads</div>
    <div class="stat-value">{workloads.length}</div>
    <div class="stat-desc">across all namespaces</div>
  </div>
  <div class="stat bg-base-100 rounded-box shadow-sm">
    <div class="stat-figure text-success"><RefreshCw size={28} /></div>
    <div class="stat-title">Active</div>
    <div class="stat-value text-success">{workloads.filter(w => w.status === 'active').length}</div>
    <div class="stat-desc">secrets healthy</div>
  </div>
  <div class="stat bg-base-100 rounded-box shadow-sm">
    <div class="stat-figure text-warning"><AlertTriangle size={28} /></div>
    <div class="stat-title">Pending Rotation</div>
    <div class="stat-value text-warning">{pendingRotation}</div>
    <div class="stat-desc">require attention</div>
  </div>
</div>

<!-- Table -->
<div class="bg-base-100 rounded-box shadow-sm overflow-hidden">
  <div class="px-5 py-4 border-b border-base-200">
    <h2 class="font-semibold text-sm">Workloads</h2>
  </div>
  <div class="overflow-x-auto">
    <table class="table table-zebra text-sm">
      <thead>
        <tr class="text-xs text-base-content/50">
          <th>Pod</th>
          <th>Namespace</th>
          <th>Secret</th>
          <th>Type</th>
          <th>Last Rotated</th>
          <th>Next Rotation</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        {#each workloads as w}
          {@const days = daysUntilRotation(w.last_rotated_at, w.rotation_days)}
          {@const rStatus = rotationStatus(w.last_rotated_at, w.rotation_days)}
          <tr class="cursor-pointer hover">
            <td class="font-mono">{w.pod_name}</td>
            <td><span class="badge badge-ghost badge-sm">{w.namespace}</span></td>
            <td class="font-mono text-base-content/70">{w.secret_name}</td>
            <td>
              {#if w.db_type}
                <span class="badge badge-info badge-sm">{w.db_type}</span>
              {:else}
                <span class="badge badge-ghost badge-sm">{w.type}</span>
              {/if}
            </td>
            <td class="text-base-content/60">{formatDate(w.last_rotated_at)}</td>
            <td>
              {#if rStatus === 'never'}
                <span class="text-warning text-xs">never rotated</span>
              {:else if rStatus === 'overdue'}
                <span class="text-error text-xs font-medium">{Math.abs(days!)}d overdue</span>
              {:else if rStatus === 'soon'}
                <span class="text-warning text-xs">{days}d left</span>
              {:else}
                <span class="text-base-content/50 text-xs">{days}d left</span>
              {/if}
            </td>
            <td>
              {#if rStatus === 'overdue' || rStatus === 'never'}
                <span class="badge badge-error badge-sm">needs rotation</span>
              {:else if rStatus === 'soon'}
                <span class="badge badge-warning badge-sm">soon</span>
              {:else}
                <span class="badge badge-success badge-sm">healthy</span>
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>
