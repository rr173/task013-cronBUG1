package cron

import (
	"testing"
	"time"
)

// TestHourAdvanceTruncatesMinute verifies that when the hour field does not
// match, the Next algorithm resets minute to 0 before advancing to the next
// hour. Without this reset, a non-zero minute carries over and causes the
// scheduler to skip the correct trigger time.
func TestHourAdvanceTruncatesMinute(t *testing.T) {
	// "0 */2 * * *" fires at minute=0 of every even hour.
	// from = 2026-01-15T11:01:00Z => initial t = 11:02.
	// hour=11 is odd => advance hour.
	// Expected: reset to 11:00 + 1h = 12:00. hour=12 even, minute=0 => match.
	s, _ := Parse("0 */2 * * *")
	from := time.Date(2026, 1, 15, 11, 1, 0, 0, time.UTC)
	got, err := s.Next(from)
	if err != nil {
		t.Fatalf("Next returned error: %v", err)
	}
	want := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Next = %v, want %v", got, want)
	}
}

// TestHourAdvanceDoesNotCarryMinute is a second case proving the minute
// truncation matters when hour changes span midnight.
func TestHourAdvanceDoesNotCarryMinute(t *testing.T) {
	// "30 6 * * *" fires at 06:30 daily.
	// from = 2026-01-15T06:31:00Z => t = 06:32.
	// hour=6 matches but minute=32 does not match 30.
	// Minute advances: 32->33->...->59->07:00.
	// hour=7 does not match 6 => advance hour. Should reset minute to 0.
	// 07:00+1h = 08:00... eventually wraps to next day 06:00, minute=0 != 30,
	// then 06:01...06:30 => match at 2026-01-16T06:30.
	s, _ := Parse("30 6 * * *")
	from := time.Date(2026, 1, 15, 6, 31, 0, 0, time.UTC)
	got, err := s.Next(from)
	if err != nil {
		t.Fatalf("Next returned error: %v", err)
	}
	want := time.Date(2026, 1, 16, 6, 30, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Next = %v, want %v", got, want)
	}
}
