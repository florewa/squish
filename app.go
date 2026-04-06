package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"strconv"
	"strings"
	"sync"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx       context.Context
	queue     []QueueItem
	mu        sync.Mutex
	converting bool
}

type QueueItem struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Type     string `json:"type"` // "image" | "video"
	Status   string `json:"status"` // "pending" | "converting" | "done" | "error"
	Progress  int    `json:"progress"`
	PassLabel string `json:"passLabel"`
	Output    string `json:"output"`
	Error     string `json:"error"`
	OldSize  int64  `json:"oldSize"`
	NewSize  int64  `json:"newSize"`
}

type ConvertSettings struct {
	// Основные
	Quality        int    `json:"quality"`
	Preset         string `json:"preset"`        // "fast" | "balanced" | "quality"
	MaxWidth       int    `json:"maxWidth"`
	MaxHeight      int    `json:"maxHeight"`
	OutputDir      string `json:"outputDir"`
	DeleteOriginal bool   `json:"deleteOriginal"`
	// Расширенные — изображения
	Lossless      bool `json:"lossless"`
	StripMetadata bool `json:"stripMetadata"`
	// Расширенные — видео
	NoAudio      bool   `json:"noAudio"`
	TwoPass      bool   `json:"twoPass"`
	AudioBitrate int    `json:"audioBitrate"` // 64 | 128 | 192 | 320
	// Расширенные — общие
	Threads int    `json:"threads"` // 1-4
	Suffix  string `json:"suffix"`  // суффикс к имени файла
}

type FFmpegStatus struct {
	Found   bool   `json:"found"`
	Path    string `json:"path"`
	Version string `json:"version"`
}

func NewApp() *App {
	return &App{queue: make([]QueueItem, 0)}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Распаковываем встроенный ffmpeg заранее, пока грузится UI
	go getFFmpegPath()
}

// CheckFFmpeg проверяет наличие ffmpeg
func (a *App) CheckFFmpeg() FFmpegStatus {
	// 1. Встроенный ffmpeg (распакован в кэш)
	if path, err := getFFmpegPath(); err == nil && path != "" {
		version := getFFmpegVersion(path)
		return FFmpegStatus{Found: true, Path: path, Version: version}
	}

	// 2. Рядом с exe
	exe, _ := os.Executable()
	localPath := filepath.Join(filepath.Dir(exe), "ffmpeg")
	if goruntime.GOOS == "windows" {
		localPath += ".exe"
	}
	if _, err := os.Stat(localPath); err == nil {
		version := getFFmpegVersion(localPath)
		return FFmpegStatus{Found: true, Path: localPath, Version: version}
	}

	// 3. В PATH
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		version := getFFmpegVersion(path)
		return FFmpegStatus{Found: true, Path: path, Version: version}
	}

	return FFmpegStatus{Found: false}
}

func getFFmpegVersion(path string) string {
	out, err := exec.Command(path, "-version").Output()
	if err != nil {
		return "unknown"
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 3 {
			return parts[2]
		}
	}
	return "unknown"
}

// OpenFFmpegDownload открывает страницу загрузки ffmpeg
func (a *App) OpenFFmpegDownload() {
	url := "https://www.gyan.dev/ffmpeg/builds/"
	switch goruntime.GOOS {
	case "darwin":
		exec.Command("open", url).Start()
	case "linux":
		exec.Command("xdg-open", url).Start()
	default:
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	}
}

// AddFiles добавляет файлы в очередь
func (a *App) AddFiles(paths []string) []QueueItem {
	a.mu.Lock()
	defer a.mu.Unlock()

	var added []QueueItem
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		fileType := detectType(path)
		if fileType == "" {
			continue
		}

		item := QueueItem{
			ID:      generateID(),
			Path:    path,
			Name:    filepath.Base(path),
			Size:    info.Size(),
			OldSize: info.Size(),
			Type:    fileType,
			Status:  "pending",
		}
		a.queue = append(a.queue, item)
		added = append(added, item)
	}
	return added
}

// AddFolder добавляет все поддерживаемые файлы из папки
func (a *App) AddFolder(dir string) []QueueItem {
	var paths []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if detectType(path) != "" {
			paths = append(paths, path)
		}
		return nil
	})
	return a.AddFiles(paths)
}

// AddPaths принимает и файлы и папки (для drag & drop)
func (a *App) AddPaths(paths []string) []QueueItem {
	var files []string
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			continue
		}
		if info.IsDir() {
			filepath.Walk(p, func(wp string, wi os.FileInfo, werr error) error {
				if werr != nil || wi.IsDir() {
					return nil
				}
				if detectType(wp) != "" {
					files = append(files, wp)
				}
				return nil
			})
		} else {
			files = append(files, p)
		}
	}
	return a.AddFiles(files)
}

// GetQueue возвращает текущую очередь
func (a *App) GetQueue() []QueueItem {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.queue
}

// ClearQueue очищает очередь (только не конвертируемые)
func (a *App) ClearQueue() {
	a.mu.Lock()
	defer a.mu.Unlock()
	var active []QueueItem
	for _, item := range a.queue {
		if item.Status == "converting" {
			active = append(active, item)
		}
	}
	if active == nil {
		active = make([]QueueItem, 0)
	}
	a.queue = active
}

// RemoveItem удаляет элемент из очереди
func (a *App) RemoveItem(id string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	var newQueue []QueueItem
	for _, item := range a.queue {
		if item.ID != id {
			newQueue = append(newQueue, item)
		}
	}
	if newQueue == nil {
		newQueue = make([]QueueItem, 0)
	}
	a.queue = newQueue
}

// StartConversion запускает конвертацию всей очереди
func (a *App) StartConversion(settings ConvertSettings) {
	if a.converting {
		return
	}
	a.converting = true

	threads := settings.Threads
	if threads < 1 {
		threads = 1
	}
	if threads > 8 {
		threads = 8
	}

	go func() {
		defer func() { a.converting = false }()

		ffmpeg := a.CheckFFmpeg()
		if !ffmpeg.Found {
			return
		}

		a.mu.Lock()
		pending := make([]QueueItem, 0)
		for _, item := range a.queue {
			if item.Status == "pending" {
				pending = append(pending, item)
			}
		}
		a.mu.Unlock()

		sem := make(chan struct{}, threads)
		var wg sync.WaitGroup
		for _, item := range pending {
			wg.Add(1)
			sem <- struct{}{}
			go func(it QueueItem) {
				defer wg.Done()
				defer func() { <-sem }()
				a.convertItem(it, settings, ffmpeg.Path)
			}(item)
		}
		wg.Wait()
	}()
}

func (a *App) convertItem(item QueueItem, settings ConvertSettings, ffmpegPath string) {
	a.updateItem(item.ID, func(i *QueueItem) {
		i.Status = "converting"
		i.Progress = 0
	})
	wailsRuntime.EventsEmit(a.ctx, "queue:update", a.GetQueue())

	outPath := buildOutputPath(item.Path, item.Type, settings.OutputDir, settings.Suffix)

	var err error
	if item.Type == "video" && settings.TwoPass {
		pass1, pass2 := buildTwoPassArgs(item, settings, outPath)
		cmd1 := exec.Command(ffmpegPath, pass1...)
		cmd1.Stderr = &progressParser{app: a, itemID: item.ID, ctx: a.ctx, passLabel: "1/2"}
		if err = cmd1.Run(); err == nil {
			cmd2 := exec.Command(ffmpegPath, pass2...)
			cmd2.Stderr = &progressParser{app: a, itemID: item.ID, ctx: a.ctx, passLabel: "2/2"}
			err = cmd2.Run()
		}
	} else {
		args := buildFFmpegArgs(item, settings, outPath)
		cmd := exec.Command(ffmpegPath, args...)
		cmd.Stderr = &progressParser{app: a, itemID: item.ID, ctx: a.ctx}
		err = cmd.Run()
	}

	a.updateItem(item.ID, func(i *QueueItem) {
		if err != nil {
			i.Status = "error"
			i.Error = err.Error()
		} else {
			i.Status = "done"
			i.Progress = 100
			i.Output = outPath
			if info, e := os.Stat(outPath); e == nil {
				i.NewSize = info.Size()
			}
			if settings.DeleteOriginal {
				os.Remove(item.Path)
			}
		}
	})
	wailsRuntime.EventsEmit(a.ctx, "queue:update", a.GetQueue())
}

func buildFFmpegArgs(item QueueItem, s ConvertSettings, outPath string) []string {
	args := []string{"-y", "-i", item.Path}

	if s.StripMetadata {
		args = append(args, "-map_metadata", "-1")
	}

	scaleFilter := ""
	if s.MaxWidth > 0 || s.MaxHeight > 0 {
		w, h := s.MaxWidth, s.MaxHeight
		if w == 0 { w = -1 }
		if h == 0 { h = -1 }
		scaleFilter = fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", w, h)
	}

	if item.Type == "image" {
		if s.Lossless {
			args = append(args, "-lossless", "1")
		} else {
			quality := s.Quality
			if quality == 0 { quality = 85 }
			args = append(args, "-quality", strconv.Itoa(quality))
		}
		if scaleFilter != "" {
			args = append(args, "-vf", scaleFilter)
		}
	} else {
		crf := 63 - (s.Quality * 63 / 100)
		if crf < 4 { crf = 4 }

		speed := "2"
		switch s.Preset {
		case "fast":    speed = "4"
		case "balanced": speed = "2"
		case "quality": speed = "0"
		}

		args = append(args, "-c:v", "libvpx-vp9", "-crf", strconv.Itoa(crf), "-b:v", "0", "-speed", speed)

		if scaleFilter != "" {
			args = append(args, "-vf", scaleFilter)
		}

		if s.NoAudio {
			args = append(args, "-an")
		} else {
			audioBitrate := s.AudioBitrate
			if audioBitrate == 0 { audioBitrate = 128 }
			args = append(args, "-c:a", "libopus", "-b:a", fmt.Sprintf("%dk", audioBitrate))
		}
	}

	args = append(args, outPath)
	return args
}

func buildTwoPassArgs(item QueueItem, s ConvertSettings, outPath string) ([]string, []string) {
	crf := 63 - (s.Quality * 63 / 100)
	if crf < 4 { crf = 4 }
	speed := "2"
	switch s.Preset {
	case "fast":     speed = "4"
	case "balanced": speed = "2"
	case "quality":  speed = "0"
	}

	base := []string{"-y", "-i", item.Path}
	if s.StripMetadata {
		base = append(base, "-map_metadata", "-1")
	}
	base = append(base, "-c:v", "libvpx-vp9", "-crf", strconv.Itoa(crf), "-b:v", "0", "-speed", speed)

	// Копируем base чтобы избежать shared slice между pass1 и pass2
	pass1 := make([]string, len(base), len(base)+6)
	copy(pass1, base)
	nullOut := "/dev/null"
	if isWindows() { nullOut = "NUL" }
	pass1 = append(pass1, "-pass", "1", "-an", "-f", "null", nullOut)

	pass2 := make([]string, len(base), len(base)+6)
	copy(pass2, base)
	pass2 = append(pass2, "-pass", "2")
	if s.NoAudio {
		pass2 = append(pass2, "-an")
	} else {
		audioBitrate := s.AudioBitrate
		if audioBitrate == 0 { audioBitrate = 128 }
		pass2 = append(pass2, "-c:a", "libopus", "-b:a", fmt.Sprintf("%dk", audioBitrate))
	}
	pass2 = append(pass2, outPath)

	return pass1, pass2
}

func isWindows() bool {
	return goruntime.GOOS == "windows"
}

func buildOutputPath(inputPath, fileType, outputDir, suffix string) string {
	dir := filepath.Dir(inputPath)
	if outputDir != "" {
		dir = outputDir
	}
	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	ext := ".webp"
	if fileType == "video" {
		ext = ".webm"
	}
	return filepath.Join(dir, base+suffix+ext)
}

func detectType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	imageExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true, ".tiff": true, ".tif": true, ".webp": false}
	videoExts := map[string]bool{".mp4": true, ".mov": true, ".avi": true, ".mkv": true, ".wmv": true, ".flv": true, ".m4v": true, ".webm": false}

	if imageExts[ext] {
		return "image"
	}
	if videoExts[ext] {
		return "video"
	}
	return ""
}

func (a *App) updateItem(id string, fn func(*QueueItem)) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i := range a.queue {
		if a.queue[i].ID == id {
			fn(&a.queue[i])
			break
		}
	}
}

// OpenFileDialog открывает диалог выбора файлов
func (a *App) OpenFileDialog() []string {
	files, err := wailsRuntime.OpenMultipleFilesDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Выбери файлы",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "Изображения", Pattern: "*.jpg;*.jpeg;*.png;*.gif;*.bmp;*.tiff"},
			{DisplayName: "Видео", Pattern: "*.mp4;*.mov;*.avi;*.mkv;*.wmv;*.flv;*.m4v"},
			{DisplayName: "Все поддерживаемые", Pattern: "*.jpg;*.jpeg;*.png;*.gif;*.bmp;*.tiff;*.mp4;*.mov;*.avi;*.mkv;*.wmv"},
		},
	})
	if err != nil {
		return nil
	}
	return files
}

// OpenFolderDialog открывает диалог выбора папки
func (a *App) OpenFolderDialog() string {
	dir, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Выбери папку",
	})
	if err != nil {
		return ""
	}
	return dir
}

// RevealInExplorer открывает папку с файлом в проводнике
func (a *App) RevealInExplorer(filePath string) {
	switch goruntime.GOOS {
	case "darwin":
		exec.Command("open", "-R", filePath).Start()
	case "linux":
		exec.Command("xdg-open", filepath.Dir(filePath)).Start()
	default:
		exec.Command("explorer", "/select,", filePath).Start()
	}
}

// SelectOutputDir выбор папки для сохранения
func (a *App) SelectOutputDir() string {
	dir, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "Куда сохранять файлы",
	})
	if err != nil {
		return ""
	}
	return dir
}

var idCounter int
var idMu sync.Mutex

func generateID() string {
	idMu.Lock()
	defer idMu.Unlock()
	idCounter++
	return fmt.Sprintf("item_%d", idCounter)
}
