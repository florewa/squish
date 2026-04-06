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

Если ffmpeg не установлен — покажет инструкцию как поставить.

## Стек

- [Wails v2](https://wails.io) — Go backend + WebView frontend
- [Svelte](https://svelte.dev) + [Tailwind CSS](https://tailwindcss.com)
- [ffmpeg](https://ffmpeg.org) — для самой конвертации

## Установка и запуск

**Зависимости:**
- [Go 1.21+](https://go.dev/dl/)
- [Node.js](https://nodejs.org/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)
- [ffmpeg](https://www.gyan.dev/ffmpeg/builds/) — либо в PATH, либо рядом с exe

```bash
# Установить Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Клонировать и запустить
git clone <repo>
cd squish
wails dev
```

## Сборка

```bash
wails build
```

Готовый `Squish.exe` появится в `build/bin/`.

## ffmpeg

Приложение ищет ffmpeg в двух местах:
1. Рядом с exe (`ffmpeg.exe` в той же папке)
2. В системном PATH

Быстрая установка через winget:
```
winget install ffmpeg
```

## Форматы

| Вход | Выход |
|------|-------|
| JPG, PNG, GIF, BMP, TIFF | WebP |
| MP4, MOV, AVI, MKV, WMV, FLV | WebM (VP9 + Opus) |
