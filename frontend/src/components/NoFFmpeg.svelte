<script>
  import { createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()
</script>

<div class="flex-1 flex flex-col items-center justify-center px-8 text-center gap-6">

  <div class="w-20 h-20 rounded-2xl bg-red-500/10 flex items-center justify-center text-4xl">
    ⚠️
  </div>

  <div class="space-y-2">
    <h2 class="text-lg font-semibold text-white/90">FFmpeg не найден</h2>
    <p class="text-sm text-white/40 leading-relaxed max-w-xs">
      FFWrap использует ffmpeg для конвертации.
      Скачай его и положи <code class="text-accent-400 bg-white/5 px-1 rounded">ffmpeg.exe</code>
      рядом с этой программой — или установи в систему.
    </p>
  </div>

  <div class="flex flex-col gap-2 w-full max-w-xs">
    <button
      on:click={() => dispatch('download')}
      class="w-full py-3 rounded-xl bg-accent-500 hover:bg-accent-400
             text-white font-semibold text-sm transition-all shadow-lg shadow-accent-500/20"
    >
      Скачать ffmpeg →
    </button>

    <p class="text-xs text-white/20">
      gyan.dev/ffmpeg/builds — Windows builds
    </p>
  </div>

  <div class="text-left bg-white/[0.03] rounded-xl p-4 w-full max-w-xs border border-white/5">
    <p class="text-xs text-white/40 font-medium mb-2">Быстрая установка через winget:</p>
    <code class="text-xs text-accent-400 font-mono">winget install ffmpeg</code>
  </div>

</div>
