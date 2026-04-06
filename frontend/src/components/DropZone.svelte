<script>
  import { createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()

  // Визуальная подсветка при drag (Wails шлёт CSS-класс на body)
  let dragOver = false

  function onDragOver(e) {
    e.preventDefault()
    dragOver = true
  }
  function onDragLeave(e) {
    if (!e.currentTarget.contains(e.relatedTarget)) dragOver = false
  }
  function onDrop(e) {
    e.preventDefault()
    dragOver = false
    // Реальные пути приходят через OnFileDrop в App.svelte
    // Здесь только сбрасываем визуал
  }
</script>

<div
  class="w-full max-w-md"
  on:dragover={onDragOver}
  on:dragleave={onDragLeave}
  on:drop={onDrop}
  role="region"
>
  <div class="
    relative rounded-2xl border-2 border-dashed transition-all duration-200 p-12
    flex flex-col items-center gap-4
    {dragOver
      ? 'border-accent-500 bg-accent-500/5'
      : 'border-white/10 hover:border-white/20 bg-white/[0.02]'}
  ">

    <div class="
      w-16 h-16 rounded-2xl flex items-center justify-center text-3xl transition-all
      {dragOver ? 'bg-accent-500/20 scale-110' : 'bg-white/5'}
    ">
      {dragOver ? '📂' : '📁'}
    </div>

    <div class="text-center">
      <p class="text-white/70 font-medium mb-1">
        {dragOver ? 'Отпускай!' : 'Перетащи файлы или папку'}
      </p>
      <p class="text-white/30 text-xs">JPG, PNG, GIF, BMP, TIFF → WebP</p>
      <p class="text-white/30 text-xs">MP4, MOV, AVI, MKV, WMV → WebM</p>
    </div>

    <div class="flex items-center gap-3 w-full max-w-[200px]">
      <div class="flex-1 h-px bg-white/10"></div>
      <span class="text-white/20 text-xs">или</span>
      <div class="flex-1 h-px bg-white/10"></div>
    </div>

    <div class="flex gap-2">
      <button
        on:click={() => dispatch('pickFiles')}
        class="px-4 py-2 rounded-lg bg-white/5 hover:bg-white/10 text-white/60 hover:text-white/90
               text-sm transition-all border border-white/5 hover:border-white/10"
      >
        Выбрать файлы
      </button>
      <button
        on:click={() => dispatch('pickFolder')}
        class="px-4 py-2 rounded-lg bg-white/5 hover:bg-white/10 text-white/60 hover:text-white/90
               text-sm transition-all border border-white/5 hover:border-white/10"
      >
        Папку
      </button>
    </div>
  </div>
</div>
