package cron

import (
	"fmt"
	"log/slog"
	"time"
)

func SafeGo(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				slog.Error(fmt.Sprintf("SafeGo panic: %s", err))
			}
		}()
		fn()
	}()
}

func RunRepeat(fn func(), dur time.Duration) {
	ticker := time.NewTicker(dur)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			SafeGo(fn)
		}
	}
}

func RunDelayRepeat(fn func(), delay, dur time.Duration) {
	if delay > 0 {
		time.Sleep(delay)
	}
	SafeGo(fn)
	RunRepeat(fn, dur)
}

func RunDay(fn func(), day int) {
	RunSpecify(fn, day, 0, 0, 0)
}

func RunDayTime(fn func(), day, hour, min int) {
	RunSpecify(fn, day, 0, hour, min)
}

func RunWeek(fn func(), week int) {
	RunSpecify(fn, 0, week, 0, 0)
}

func RunWeekTime(fn func(), week, hour, min int) {
	RunSpecify(fn, 0, week, hour, min)
}

func RunTime(fn func(), hour, min int) {
	RunSpecify(fn, 0, 0, hour, min)
}

func RunSpecify(fn func(), day, week, hour, min int) {
	now := time.Now()
	target := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
	if target.Before(now) {
		target.AddDate(0, 0, 1)
	}
	RunDelayRepeat(func() {
		execTime := time.Now()
		if day > 0 && execTime.Day() != day {
			return
		}
		if week > 0 && int(execTime.Weekday()) != week {
			return
		}
		fn()
	}, target.Sub(now), time.Hour*24)
}
