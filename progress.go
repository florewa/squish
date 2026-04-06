package main

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// progressParser парсит stderr ffmpeg и эмитит прогресс
type progressParser struct {
	app       *App
	itemID    string
	ctx       context.Context
	duration  float64
	lastEmit  time.Time
	passLabel string // "1/2" | "2/2" | ""
}

var durationRe = regexp.MustCompile(`Duration:\s+(\d+):(\d+):(\d+)\.(\d+)`)
var timeRe = regexp.MustCompile(`time=(\d+):(\d+):(\d+)\.(\d+)`)

func (p *progressParser) Write(b []byte) (int, error) {
	s := string(b)

	if p.duration == 0 {
		if m := durationRe.FindStringSubmatch(s); m != nil {
			p.duration = parseTimeToSeconds(m[1], m[2], m[3], m[4])
		}
	}

	if m := timeRe.FindStringSubmatch(s); m != nil {
		current := parseTimeToSeconds(m[1], m[2], m[3], m[4])
		if p.duration > 0 {
			progress := int((current / p.duration) * 100)
			if progress > 99 {
				progress = 99
			}

			// Эмитим не чаще раза в 100ms
			if time.Since(p.lastEmit) > 100*time.Millisecond {
				p.lastEmit = time.Now()
				label := p.passLabel
				p.app.updateItem(p.itemID, func(i *QueueItem) {
					i.Progress = progress
					if label != "" {
						i.PassLabel = label
					}
				})
				wailsRuntime.EventsEmit(p.ctx, "queue:update", p.app.GetQueue())
			}
		}
	}

	// Для изображений прогресс не парсится — просто ставим 50% пока идёт
	if strings.Contains(s, "Output #0") {
		p.app.updateItem(p.itemID, func(i *QueueItem) {
			if i.Progress < 50 {
				i.Progress = 50
			}
		})
	}

	return len(b), nil
}

func parseTimeToSeconds(h, m, s, ms string) float64 {
	hours, _ := strconv.ParseFloat(h, 64)
	minutes, _ := strconv.ParseFloat(m, 64)
	seconds, _ := strconv.ParseFloat(s, 64)
	return hours*3600 + minutes*60 + seconds
}
