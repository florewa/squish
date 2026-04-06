# Squish

Конвертирует картинки в WebP и видео в WebM. Drag & drop, пакетная обработка, никаких лишних кнопок.

![Squish UI](build/appicon.png)

## Возможности

- Перетащи файлы или папку прямо в окно
- Пакетная обработка с очередью и прогрессом
- Настройка качества, пресетов скорости и размера
- Показывает сколько места сэкономил (или нет)
- Открывает папку с результатом одним кликом
- Расширенные настройки: lossless, убрать EXIF, two-pass для видео, параллельные потоки, суффикс к имени

## Использование

Скачай `Squish.exe` из [Releases](../../releases) и запусти. ffmpeg уже внутри, ничего дополнительно устанавливать не нужно.

## Стек

- [Wails v2](https://wails.io) — Go backend + WebView frontend
- [Svelte](https://svelte.dev) + [Tailwind CSS](https://tailwindcss.com)
- [ffmpeg](https://ffmpeg.org) — для самой конвертации

## Сборка из исходников

**Зависимости:**
- [Go 1.21+](https://go.dev/dl/)
- [Node.js](https://nodejs.org/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)
- `ffmpeg.exe` в корне проекта (для вшивания в бинарь)

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
git clone https://github.com/florewa/squish.git
cd squish
wails build
```

## Форматы

| Вход | Выход |
|------|-------|
| JPG, PNG, GIF, BMP, TIFF | WebP |
| MP4, MOV, AVI, MKV, WMV, FLV | WebM (VP9 + Opus) |
