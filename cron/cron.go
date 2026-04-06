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
	if dur <= 0 {
		return
	}
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
	if dur <= 0 {
		return
	}
	if delay > 0 {
		time.Sleep(delay)
	}
	SafeGo(fn)
	RunRepeat(fn, dur)
}

func RunDay(fn func(), day int) {
	RunSpecify(fn, day, -1, 0, 0)
}

func RunDayTime(fn func(), day, hour, min int) {
	RunSpecify(fn, day, -1, hour, min)
}

func RunWeek(fn func(), week time.Weekday) {
	RunSpecify(fn, -1, week, 0, 0)
}

func RunWeekTime(fn func(), week time.Weekday, hour, min int) {
	RunSpecify(fn, -1, week, hour, min)
}

func RunTime(fn func(), hour, min int) {
	RunSpecify(fn, -1, -1, hour, min)
}

func RunSpecify(fn func(), day int, week time.Weekday, hour, min int) {
	if err := validateSpecifyArgs(day, week, hour, min); err != nil {
		slog.Error("RunSpecify invalid args", "day", day, "week", week, "hour", hour, "minute", min, "error", err)
		return
	}
	for {
		target := nextSpecifyTime(time.Now(), day, week, hour, min)
		<-time.After(time.Until(target))
		SafeGo(fn)
	}
}

func nextSpecifyTime(now time.Time, day int, week time.Weekday, hour, min int) time.Time {
	target := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
	for target.Before(now) || !matchSpecifyDate(target, day, week) {
		target = target.AddDate(0, 0, 1)
		target = time.Date(target.Year(), target.Month(), target.Day(), hour, min, 0, 0, target.Location())
	}
	return target
}

func matchSpecifyDate(t time.Time, day int, week time.Weekday) bool {
	if day > 0 && t.Day() != day {
		return false
	}
	if week >= 0 && t.Weekday() != week {
		return false
	}
	return true
}

func validateSpecifyArgs(day int, week time.Weekday, hour, min int) error {
	if day > 31 {
		return fmt.Errorf("cron day must be <= 31: %d", day)
	}
	if week > time.Saturday {
		return fmt.Errorf("cron week must be less than or equal to %d: %d", time.Saturday, week)
	}
	if hour < 0 || hour > 23 {
		return fmt.Errorf("cron hour must be between 0 and 23: %d", hour)
	}
	if min < 0 || min > 59 {
		return fmt.Errorf("cron minute must be between 0 and 59: %d", min)
	}
	return nil
}
