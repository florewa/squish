<script>
  import { createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()

  export let queue = []

  function formatSize(bytes) {
    if (!bytes) return '—'
    if (bytes < 1024) return bytes + ' B'
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(0) + ' KB'
    return (bytes / 1024 / 1024).toFixed(1) + ' MB'
  }

  function savings(item) {
    if (!item.newSize || !item.oldSize) return null
    const pct = Math.round((1 - item.newSize / item.oldSize) * 100)
    if (pct > 0) return { sign: '−', value: pct, grew: false }
    if (pct < 0) return { sign: '+', value: Math.abs(pct), grew: true }
    return { sign: '', value: 0, grew: false }
  }

  const typeIcon = { image: '🖼', video: '🎬' }

  const statusColor = {
    pending:    'text-white/30',
    converting: 'text-violet-400',
    done:       'text-emerald-400',
    error:      'text-red-400',
  }

  const statusLabel = {
    pending:    'Ожидание',
    converting: 'Конвертирую...',
    done:       'Готово',
    error:      'Ошибка',
  }
</script>

<div class="px-5 py-3 space-y-1.5">
  {#each queue as item (item.id)}
    <div class="flex items-center gap-3 px-3 py-2.5 rounded-xl
                bg-white/[0.03] hover:bg-white/[0.05] transition-colors
                border border-white/[0.04]">

      <span class="text-base flex-shrink-0">{typeIcon[item.type] ?? '📄'}</span>

      <div class="flex-1 min-w-0">
        <div class="flex items-center justify-between gap-2 mb-1.5">
          <span class="text-sm text-white/80 truncate font-medium" title={item.path}>
            {item.name}
          </span>
          <span class="text-xs flex-shrink-0 {statusColor[item.status]}">
            {#if item.status === 'done' && savings(item) !== null}
              {@const s = savings(item)}
              <span class={s.grew ? 'text-orange-400' : ''}>{s.sign}{s.value}% · {formatSize(item.oldSize)} → {formatSize(item.newSize)}</span>
            {:else}
              {statusLabel[item.status]}{item.passLabel ? ` (проход ${item.passLabel})` : ''}
            {/if}
          </span>
        </div>

        <!-- Прогресс-бар -->
        <div class="h-0.5 rounded-full bg-white/5 overflow-hidden">
          {#if item.status === 'converting'}
            <div class="h-full rounded-full bg-violet-500 transition-all duration-300"
                 style="width: {item.progress}%"></div>
          {:else if item.status === 'done'}
            <div class="h-full rounded-full bg-emerald-400" style="width: 100%"></div>
          {:else if item.status === 'error'}
            <div class="h-full rounded-full bg-red-500" style="width: 100%"></div>
          {/if}
        </div>

        {#if item.status === 'error'}
          <p class="text-xs text-red-400/70 truncate mt-1">{item.error}</p>
        {/if}
      </div>

      <!-- Кнопки справа -->
      <div class="flex items-center gap-1 flex-shrink-0">
        {#if item.status === 'done'}
          <button
            on:click={() => dispatch('reveal', item.output)}
            class="w-6 h-6 flex items-center justify-center rounded-md
                   text-white/25 hover:text-white/80 hover:bg-white/10 transition-all"
            title="Показать в проводнике"
          >
            <!-- иконка папки со стрелкой -->
            <svg width="14" height="14" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M2 4.5C2 3.67 2.67 3 3.5 3H6.38L7.88 4.5H12.5C13.33 4.5 14 5.17 14 6V11.5C14 12.33 13.33 13 12.5 13H3.5C2.67 13 2 12.33 2 11.5V4.5Z" stroke="currentColor" stroke-width="1.2" stroke-linejoin="round"/>
              <path d="M10 8.5L12 6.5M12 6.5L10 4.5M12 6.5H7" stroke="currentColor" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </button>
        {/if}
        {#if item.status !== 'converting'}
          <button
            on:click={() => dispatch('remove', item.id)}
            class="w-6 h-6 flex items-center justify-center rounded-md
                   text-white/25 hover:text-white/80 hover:bg-white/10
                   transition-all text-lg leading-none"
            title="Удалить из очереди"
          >×</button>
        {:else}
          <div class="w-6"></div>
        {/if}
      </div>

    </div>
  {/each}
</div>
