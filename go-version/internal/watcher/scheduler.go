package watcher

import (
	"fmt"
	"time"
)

// ScheduleConfig holds scheduling configuration
type ScheduleConfig struct {
	DaytimeInterval int     // Minutes between checks during daytime
	NightChecks     []int   // Hours to check during night (e.g., [2, 6])
	NightStart      float64 // Night start hour (e.g., 0 for midnight)
	NightEnd        float64 // Night end hour (e.g., 8.5 for 8:30 AM)
}

// DefaultScheduleConfig returns the default schedule configuration
func DefaultScheduleConfig() *ScheduleConfig {
	return &ScheduleConfig{
		DaytimeInterval: 15,
		NightChecks:     []int{2, 6},
		NightStart:      0,
		NightEnd:        8.5,
	}
}

// Scheduler handles watch daemon scheduling
type Scheduler struct {
	config *ScheduleConfig
}

// NewScheduler creates a new Scheduler instance
func NewScheduler(config *ScheduleConfig) *Scheduler {
	if config == nil {
		config = DefaultScheduleConfig()
	}
	return &Scheduler{
		config: config,
	}
}

// CalculateNextCheckTime calculates the next check time based on current time
func (s *Scheduler) CalculateNextCheckTime(now time.Time) time.Time {
	currentHour := float64(now.Hour()) + float64(now.Minute())/60.0

	// Check if we're in night mode
	if currentHour >= s.config.NightStart && currentHour < s.config.NightEnd {
		return s.calculateNextNightCheck(now)
	}

	// Daytime mode: add interval minutes
	return now.Add(time.Duration(s.config.DaytimeInterval) * time.Minute)
}

// calculateNextNightCheck calculates the next night check time
func (s *Scheduler) calculateNextNightCheck(now time.Time) time.Time {
	currentHour := now.Hour()

	// Find next night check hour
	for _, hour := range s.config.NightChecks {
		if hour > currentHour {
			// Next check is today at this hour
			nextCheck := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())
			return nextCheck
		}
	}

	// No more night checks today, schedule first daytime check
	nextDayStart := now.Add(24 * time.Hour)
	daytimeStartHour := int(s.config.NightEnd)
	daytimeStartMinute := int((s.config.NightEnd - float64(daytimeStartHour)) * 60)
	
	nextCheck := time.Date(nextDayStart.Year(), nextDayStart.Month(), nextDayStart.Day(), 
		daytimeStartHour, daytimeStartMinute, 0, 0, now.Location())
	
	return nextCheck
}

// SleepUntilNextCheck sleeps until the next check time and returns a channel
func (s *Scheduler) SleepUntilNextCheck(now time.Time) <-chan time.Time {
	nextCheck := s.CalculateNextCheckTime(now)
	duration := nextCheck.Sub(now)
	
	if duration < 0 {
		duration = 0
	}
	
	return time.After(duration)
}

// GetCurrentMode returns the current mode (daytime or night)
func (s *Scheduler) GetCurrentMode(now time.Time) string {
	currentHour := float64(now.Hour()) + float64(now.Minute())/60.0

	if currentHour >= s.config.NightStart && currentHour < s.config.NightEnd {
		return "night"
	}
	return "daytime"
}

// FormatNextCheckTime formats the next check time as a human-readable string
func (s *Scheduler) FormatNextCheckTime(nextCheck time.Time) string {
	now := time.Now()
	duration := nextCheck.Sub(now)

	if duration < 0 {
		return "now"
	}

	if duration < time.Minute {
		return "in less than a minute"
	}

	if duration < time.Hour {
		minutes := int(duration.Minutes())
		return fmt.Sprintf("in %d minute(s)", minutes)
	}

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	if minutes == 0 {
		return fmt.Sprintf("in %d hour(s)", hours)
	}
	return fmt.Sprintf("in %d hour(s) %d minute(s)", hours, minutes)
}

