<script>
  import { createEventDispatcher } from 'svelte'
  export let settings
  export let outputDir = ''

  const dispatch = createEventDispatcher()

  const presets = [
    { id: 'fast',     label: 'Быстро' },
    { id: 'balanced', label: 'Баланс' },
    { id: 'quality',  label: 'Качество' },
  ]

  const resizeOptions = [
    { label: 'Оригинал', w: 0,    h: 0    },
    { label: '4K',       w: 3840, h: 2160 },
    { label: '1080p',    w: 1920, h: 1080 },
    { label: '720p',     w: 1280, h: 720  },
    { label: '480p',     w: 854,  h: 480  },
  ]

  const audioBitrateOptions = [64, 128, 192, 320]
  const threadOptions = [1, 2, 3, 4]

  $: currentResize = resizeOptions.find(r => r.w === settings.maxWidth && r.h === settings.maxHeight) || resizeOptions[0]

  function setResize(opt) {
    settings = { ...settings, maxWidth: opt.w, maxHeight: opt.h }
  }

  let showAdvanced = false
</script>

<div class="px-5 py-3 space-y-3">

  <!-- Качество -->
  <div class="flex items-center gap-4">
    <span class="text-xs text-white/40 w-16 flex-shrink-0">Качество</span>
    <input type="range" min="1" max="100"
      bind:value={settings.quality}
      disabled={settings.lossless}
      class="flex-1 accent-violet-500 h-1 cursor-pointer disabled:opacity-30"
    />
    <span class="text-sm font-mono text-white/70 w-8 text-right">
      {settings.lossless ? '∞' : settings.quality}
    </span>
  </div>

  <!-- Скорость -->
  <div class="flex items-center gap-4">
    <span class="text-xs text-white/40 w-16 flex-shrink-0">Скорость</span>
    <div class="flex gap-1">
      {#each presets as p}
        <button on:click={() => settings = { ...settings, preset: p.id }}
          class="px-3 py-1 rounded-lg text-xs transition-all
            {settings.preset === p.id
              ? 'bg-accent-500 text-white'
              : 'bg-white/5 text-white/40 hover:bg-white/10 hover:text-white/70'}">
          {p.label}
        </button>
      {/each}
    </div>
  </div>

  <!-- Размер -->
  <div class="flex items-center gap-4">
    <span class="text-xs text-white/40 w-16 flex-shrink-0">Размер</span>
    <div class="flex gap-1 flex-wrap">
      {#each resizeOptions as opt}
        <button on:click={() => setResize(opt)}
          class="px-3 py-1 rounded-lg text-xs transition-all
            {currentResize.label === opt.label
              ? 'bg-accent-500 text-white'
              : 'bg-white/5 text-white/40 hover:bg-white/10 hover:text-white/70'}">
          {opt.label}
        </button>
      {/each}
    </div>
  </div>

  <!-- Удалить оригинал -->
  <div class="flex items-center gap-4">
    <span class="text-xs text-white/40 w-16 flex-shrink-0">Оригинал</span>
    <button on:click={() => settings = { ...settings, deleteOriginal: !settings.deleteOriginal }}
      class="flex items-center gap-2 text-xs transition-all
        {settings.deleteOriginal ? 'text-red-400' : 'text-white/40 hover:text-white/60'}">
      <div class="w-8 h-4 rounded-full relative transition-all
        {settings.deleteOriginal ? 'bg-red-500/40' : 'bg-white/10'}">
        <div class="absolute top-0.5 w-3 h-3 rounded-full transition-all
          {settings.deleteOriginal ? 'left-4 bg-red-400' : 'left-0.5 bg-white/40'}"></div>
      </div>
      {settings.deleteOriginal ? 'Удалить после сжатия' : 'Оставить оригинал'}
    </button>
  </div>

  <!-- Куда сохранять -->
  <div class="flex items-center gap-4">
    <span class="text-xs text-white/40 w-16 flex-shrink-0">Куда</span>
    <button on:click={() => dispatch('selectOutputDir')}
      class="flex-1 text-left text-xs px-3 py-1.5 rounded-lg bg-white/5
             hover:bg-white/10 transition-colors border border-white/5
             text-white/50 hover:text-white/80 truncate">
      {outputDir || 'Рядом с оригиналом'}
    </button>
    {#if outputDir}
      <button on:click={() => { settings = { ...settings, outputDir: '' }; outputDir = '' }}
        class="text-white/20 hover:text-white/60 text-lg leading-none">×</button>
    {/if}
  </div>

  <!-- Разделитель + тоггл расширенных -->
  <button
    on:click={() => showAdvanced = !showAdvanced}
    class="flex items-center gap-2 text-xs text-white/25 hover:text-white/50 transition-colors w-full pt-1"
  >
    <div class="flex-1 h-px bg-white/5"></div>
    <span>Дополнительно</span>
    <svg class="w-3 h-3 transition-transform {showAdvanced ? 'rotate-180' : ''}"
      viewBox="0 0 12 12" fill="none">
      <path d="M2 4l4 4 4-4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
    </svg>
    <div class="flex-1 h-px bg-white/5"></div>
  </button>

  <!-- Расширенные настройки -->
  {#if showAdvanced}
    <div class="space-y-3 pb-1">

      <!-- Изображения -->
      <div class="text-[10px] text-white/20 uppercase tracking-wider font-medium">Изображения</div>
      <div class="flex items-center gap-4">
        <span class="text-xs text-white/40 w-16 flex-shrink-0">Режим</span>
        <div class="flex gap-2 flex-wrap">
          <button on:click={() => settings = { ...settings, lossless: !settings.lossless }}
            class="flex items-center gap-1.5 px-3 py-1 rounded-lg text-xs transition-all
              {settings.lossless
                ? 'bg-accent-500 text-white'
                : 'bg-white/5 text-white/40 hover:bg-white/10 hover:text-white/70'}">
            Lossless
          </button>
          <button on:click={() => settings = { ...settings, stripMetadata: !settings.stripMetadata }}
            class="flex items-center gap-1.5 px-3 py-1 rounded-lg text-xs transition-all
              {settings.stripMetadata
                ? 'bg-accent-500 text-white'
                : 'bg-white/5 text-white/40 hover:bg-white/10 hover:text-white/70'}">
            Убрать EXIF
          </button>
        </div>
      </div>

      <!-- Видео -->
      <div class="text-[10px] text-white/20 uppercase tracking-wider font-medium pt-1">Видео</div>
      <div class="flex items-center gap-4">
        <span class="text-xs text-white/40 w-16 flex-shrink-0">Режим</span>
        <div class="flex gap-2 flex-wrap">
          <button on:click={() => settings = { ...settings, twoPass: !settings.twoPass }}
            class="flex items-center gap-1.5 px-3 py-1 rounded-lg text-xs transition-all
              {settings.twoPass
                ? 'bg-accent-500 text-white'
                : 'bg-white/5 text-white/40 hover:bg-white/10 hover:text-white/70'}"
            title="Два прохода — лучше качество при том же размере, но в 2х медленнее">
            Two-pass
          </button>
          <button on:click={() => settings = { ...settings, noAudio: !settings.noAudio }}
            class="flex items-center gap-1.5 px-3 py-1 rounded-lg text-xs transition-all
              {settings.noAudio
                ? 'bg-accent-500 text-white'
                : 'bg-white/5 text-white/40 hover:bg-white/10 hover:text-white/70'}">
            Без аудио
          </button>
        </div>
      </div>

      {#if !settings.noAudio}
        <div class="flex items-center gap-4">
          <span class="text-xs text-white/40 w-16 flex-shrink-0">Аудио</span>
          <div class="flex gap-1">
            {#each audioBitrateOptions as br}
              <button on:click={() => settings = { ...settings, audioBitrate: br }}
                class="px-3 py-1 rounded-lg text-xs transition-all
                  {(settings.audioBitrate || 128) === br
                    ? 'bg-accent-500 text-white'
                    : 'bg-white/5 text-white/40 hover:bg-white/10 hover:text-white/70'}">
                {br}k
              </button>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Общее -->
      <div class="text-[10px] text-white/20 uppercase tracking-wider font-medium pt-1">Обработка</div>
      <div class="flex items-center gap-4">
        <span class="text-xs text-white/40 w-16 flex-shrink-0">Потоков</span>
        <div class="flex gap-1">
          {#each threadOptions as t}
            <button on:click={() => settings = { ...settings, threads: t }}
              class="w-8 py-1 rounded-lg text-xs transition-all
                {(settings.threads || 1) === t
                  ? 'bg-accent-500 text-white'
                  : 'bg-white/5 text-white/40 hover:bg-white/10 hover:text-white/70'}">
              {t}
            </button>
          {/each}
        </div>
        <span class="text-xs text-white/20">файлов одновременно</span>
      </div>

      <div class="flex items-center gap-4">
        <span class="text-xs text-white/40 w-16 flex-shrink-0">Суффикс</span>
        <input
          type="text"
          placeholder="_compressed"
          bind:value={settings.suffix}
          class="flex-1 bg-white/5 border border-white/5 rounded-lg px-3 py-1.5
                 text-xs text-white/70 placeholder-white/20
                 focus:outline-none focus:border-accent-500/50 transition-colors"
        />
      </div>

    </div>
  {/if}

</div>
