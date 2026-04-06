# FFWrap — Установка и запуск

## 1. Установить Go
https://go.dev/dl/ → скачать Windows installer (.msi) → установить

## 2. Установить Wails
```
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

## 3. Запустить в dev-режиме (живая перезагрузка)
```
cd ffwrap
wails dev
```

## 4. Собрать .exe для раздачи коллегам
```
wails build
```
Готовый `FFWrap.exe` появится в папке `build/bin/`

---

## ffmpeg для коллег
Положи `ffmpeg.exe` рядом с `FFWrap.exe`.
Или пусть установят через winget:
```
winget install ffmpeg
```
