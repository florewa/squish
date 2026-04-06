<script>
  import { onMount } from 'svelte'
  import { EventsOn, OnFileDrop } from '../wailsjs/runtime/runtime.js'
  import {
    CheckFFmpeg, AddFiles, AddFolder, AddPaths, GetQueue,
    ClearQueue, RemoveItem, StartConversion,
    OpenFileDialog, OpenFolderDialog, SelectOutputDir,
    OpenFFmpegDownload, RevealInExplorer
  } from '../wailsjs/go/main/App.js'

  import DropZone from './components/DropZone.svelte'
  import FileQueue from './components/FileQueue.svelte'
  import Settings from './components/Settings.svelte'
  import NoFFmpeg from './components/NoFFmpeg.svelte'

  let ffmpegStatus = null
  let queue = []
  let converting = false
  let outputDir = ''
  let draggingOver = false
  let dragCounter = 0
  let loading = false

  let settings = {
    quality: 85,
    preset: 'balanced',
    maxWidth: 0,
    maxHeight: 0,
    outputDir: '',
    deleteOriginal: false,
    // расширенные
    lossless: false,
    stripMetadata: false,
    noAudio: false,
    twoPass: false,
    audioBitrate: 128,
    threads: 1,
    suffix: '',
  }

  onMount(async () => {
    ffmpegStatus = await CheckFFmpeg()
    queue = (await GetQueue()) || []

    // Показываем оверлей при любом drag над окном
    document.addEventListener('dragenter', e => { e.preventDefault(); dragCounter++; draggingOver = true })
    document.addEventListener('dragleave', () => { dragCounter--; if (dragCounter <= 0) { dragCounter = 0; draggingOver = false } })
    document.addEventListener('dragover', e => { e.preventDefault() })
    document.addEventListener('drop', e => { e.preventDefault(); dragCounter = 0; draggingOver = false })

    EventsOn('queue:update', (newQueue) => {
      queue = newQueue || []
      converting = queue.some(i => i.status === 'converting')
    })

    OnFileDrop(async (x, y, paths) => {
      loading = true
      await AddPaths(paths)
      queue = (await GetQueue()) || []
      loading = false
    }, true)
  })

  async function handleDrop(paths) {
    const added = await AddFiles(paths)
    queue = await GetQueue()
  }

  async function handleFolderDrop(path) {
    await AddFolder(path)
    queue = await GetQueue()
  }

  async function handlePickFiles() {
    const files = await OpenFileDialog()
    if (files && files.length) {
      loading = true
      await AddFiles(files)
      queue = (await GetQueue()) || []
      loading = false
    }
  }

  async function handlePickFolder() {
    const dir = await OpenFolderDialog()
    if (dir) {
      loading = true
      await AddFolder(dir)
      queue = (await GetQueue()) || []
      loading = false
    }
  }

  async function handleRemove(id) {
    await RemoveItem(id)
    queue = (await GetQueue()) || []
  }

  async function handleClear() {
    await ClearQueue()
    queue = (await GetQueue()) || []
    dragCounter = 0
    draggingOver = false
    loading = false
    converting = false
  }

  async function handleStart() {
    converting = true
    await StartConversion(settings)
  }

  async function handleSelectOutputDir() {
    const dir = await SelectOutputDir()
    if (dir) {
      settings = { ...settings, outputDir: dir }
      outputDir = dir
    }
  }

  $: pendingCount = queue.filter(i => i.status === 'pending').length
  $: doneCount = queue.filter(i => i.status === 'done').length
  $: hasQueue = queue.length > 0
  $: canStart = pendingCount > 0 && !converting && ffmpegStatus?.found
</script>

<div class="flex flex-col h-screen bg-[#0d0d0d] text-white relative" style="--wails-drop-target: drop">

  <!-- Header -->
  <div class="flex items-center justify-between px-5 py-4 border-b border-white/5">
    <div class="flex items-center gap-2.5">
      <svg width="28" height="28" viewBox="0 0 256 256" fill="none" xmlns="http://www.w3.org/2000/svg">
        <defs>
          <linearGradient id="hbg" x1="0" y1="0" x2="256" y2="256" gradientUnits="userSpaceOnUse">
            <stop stop-color="#a78bfa"/>
            <stop offset="1" stop-color="#5b21b6"/>
          </linearGradient>
        </defs>
        <rect width="256" height="256" rx="60" fill="url(#hbg)"/>
        <path d="M118,128 L70,90 L70,114 L44,114 L44,142 L70,142 L70,166 Z" fill="white" fill-opacity="0.95"/>
        <path d="M138,128 L186,90 L186,114 L212,114 L212,142 L186,142 L186,166 Z" fill="white" fill-opacity="0.95"/>
        <rect x="121" y="116" width="14" height="24" rx="4" fill="white"/>
      </svg>
      <span class="font-semibold text-[15px] tracking-tight">Squish</span>
    </div>
  </div>

  <!-- Лоадер при добавлении файлов -->
  {#if loading}
    <div class="absolute inset-0 z-50 flex items-center justify-center bg-[#0d0d0d]/60 backdrop-blur-sm pointer-events-none">
      <div class="flex flex-col items-center gap-3">
        <div class="w-8 h-8 rounded-full border-2 border-white/10 border-t-accent-500 animate-spin"></div>
        <span class="text-xs text-white/40">Добавляю файлы...</span>
      </div>
    </div>
  {/if}

  <!-- Drag-оверлей поверх очереди -->
  {#if draggingOver && hasQueue}
    <div class="absolute inset-0 z-50 pointer-events-none flex items-center justify-center
                bg-[#0d0d0d]/80 backdrop-blur-sm">
      <div class="flex flex-col items-center gap-4 rounded-2xl border-2 border-dashed
                  border-accent-500 bg-accent-500/5 px-16 py-12">
        <div class="text-4xl">📂</div>
        <p class="text-white/80 font-medium">Отпускай!</p>
        <p class="text-white/30 text-xs">Файлы и папки добавятся в очередь</p>
      </div>
    </div>
  {/if}

  {#if ffmpegStatus === null}
    <!-- Загрузка -->
    <div class="flex-1 flex items-center justify-center text-white/30 text-sm">Проверка ffmpeg...</div>

  {:else if !ffmpegStatus.found}
    <!-- ffmpeg не найден -->
    <NoFFmpeg on:download={OpenFFmpegDownload} />

  {:else}
    <!-- Основной интерфейс -->
    <div class="flex flex-col flex-1 overflow-hidden">

      <!-- Drop zone (только если очередь пустая) -->
      {#if !hasQueue}
        <div class="flex-1 flex flex-col items-center justify-center p-6">
          <DropZone
            on:drop={e => handleDrop(e.detail)}
            on:folderDrop={e => handleFolderDrop(e.detail)}
            on:pickFiles={handlePickFiles}
            on:pickFolder={handlePickFolder}
          />
        </div>
      {:else}
        <!-- Очередь -->
        <div class="flex-1 overflow-y-auto">
          <FileQueue
            {queue}
            on:remove={e => handleRemove(e.detail)}
            on:reveal={e => RevealInExplorer(e.detail)}
          />
        </div>

        <!-- Кнопка добавить ещё -->
        <div class="px-5 py-2 flex items-center gap-2">
          <button
            on:click={handlePickFiles}
            class="text-xs text-white/40 hover:text-white/70 transition-colors"
          >
            + Файлы
          </button>
          <span class="text-white/20">·</span>
          <button
            on:click={handlePickFolder}
            class="text-xs text-white/40 hover:text-white/70 transition-colors"
          >
            + Папку
          </button>
          {#if !converting}
            <span class="text-white/20 ml-auto">·</span>
            <button
              on:click={handleClear}
              class="text-xs text-white/25 hover:text-red-400 transition-colors"
            >
              Очистить всё
            </button>
          {/if}
        </div>
      {/if}

      <!-- Settings + Action -->
      {#if hasQueue}
        <div class="border-t border-white/5 bg-[#111]">
          <Settings
            bind:settings
            {outputDir}
            on:selectOutputDir={handleSelectOutputDir}
          />

          <div class="px-5 pb-5">
            <button
              on:click={handleStart}
              disabled={!canStart}
              class="w-full py-3 rounded-xl font-semibold text-sm transition-all
                {canStart
                  ? 'bg-accent-500 hover:bg-accent-400 text-white cursor-pointer shadow-lg shadow-accent-500/20'
                  : 'bg-white/5 text-white/20 cursor-not-allowed'}"
            >
              {#if converting}
                Конвертирую...
              {:else if pendingCount === 0 && doneCount > 0}
                Готово ✓
              {:else}
                Конвертировать {pendingCount > 0 ? `(${pendingCount})` : ''}
              {/if}
            </button>
          </div>
        </div>
      {/if}

    </div>
  {/if}
</div>
