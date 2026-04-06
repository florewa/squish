package main

import (
	_ "embed"
	"os"
	"path/filepath"
	"sync"
)

//go:embed ffmpeg.exe
var ffmpegBinary []byte

var (
	extractedFFmpegPath string
	extractOnce         sync.Once
)

// getFFmpegPath возвращает путь к ffmpeg — сначала извлекает из бинаря если нужно
func getFFmpegPath() (string, error) {
	var extractErr error
	extractOnce.Do(func() {
		appData, err := os.UserCacheDir()
		if err != nil {
			appData = os.TempDir()
		}
		dir := filepath.Join(appData, "Squish")
		if err := os.MkdirAll(dir, 0755); err != nil {
			extractErr = err
			return
		}

		dest := filepath.Join(dir, "ffmpeg.exe")

		// Перезаписываем только если файл изменился (сравниваем размер)
		if info, err := os.Stat(dest); err != nil || info.Size() != int64(len(ffmpegBinary)) {
			if err := os.WriteFile(dest, ffmpegBinary, 0755); err != nil {
				extractErr = err
				return
			}
		}

		extractedFFmpegPath = dest
	})
	return extractedFFmpegPath, extractErr
}
