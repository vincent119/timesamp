// Package converter 提供時間戳轉換功能的測試
package converter

import (
	"testing"
	"time"
)

func TestNewConverter(t *testing.T) {
	tests := []struct {
		name     string
		timezone string
		wantErr  bool
	}{
		{"empty timezone defaults to Local", "", false},
		{"Local timezone", "Local", false},
		{"UTC timezone", "UTC", false},
		{"Asia/Taipei timezone", "Asia/Taipei", false},
		{"America/New_York timezone", "America/New_York", false},
		{"invalid timezone", "Invalid/Zone", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conv, err := NewConverter(tt.timezone)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConverter(%q) error = %v, wantErr %v", tt.timezone, err, tt.wantErr)
				return
			}
			if !tt.wantErr && conv == nil {
				t.Errorf("NewConverter(%q) returned nil converter", tt.timezone)
			}
		})
	}
}

func TestDetectFormat(t *testing.T) {
	conv, _ := NewConverter("UTC")

	tests := []struct {
		name       string
		input      string
		wantFormat TimestampFormat
		wantErr    bool
	}{
		{"unix seconds 10 digits", "1642781234", UnixSeconds, false},
		{"unix milliseconds 13 digits", "1642781234000", UnixMilliseconds, false},
		{"unix microseconds 16 digits", "1642781234000000", UnixMicroseconds, false},
		{"unix nanoseconds 19 digits", "1642781234000000000", UnixNanoseconds, false},
		{"RFC3339 format", "2022-01-21T12:00:34Z", RFC3339, false},
		{"RFC3339 with timezone", "2022-01-21T12:00:34+08:00", RFC3339, false},
		{"RFC3339Nano format", "2022-01-21T12:00:34.123456789Z", RFC3339Nano, false},
		{"DateTime format", "2022-01-21 12:00:34", DateTime, false},
		{"DateOnly format", "2022-01-21", DateOnly, false},
		{"TimeOnly format", "12:00:34", TimeOnly, false},
		{"invalid format", "not-a-timestamp", UnixSeconds, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format, err := conv.DetectFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectFormat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && format != tt.wantFormat {
				t.Errorf("DetectFormat(%q) = %v, want %v", tt.input, format, tt.wantFormat)
			}
		})
	}
}

func TestParse(t *testing.T) {
	conv, _ := NewConverter("UTC")

	tests := []struct {
		name    string
		input   string
		format  TimestampFormat
		wantErr bool
	}{
		{"parse unix seconds", "1642781234", UnixSeconds, false},
		{"parse unix milliseconds", "1642781234000", UnixMilliseconds, false},
		{"parse unix microseconds", "1642781234000000", UnixMicroseconds, false},
		{"parse unix nanoseconds", "1642781234000000000", UnixNanoseconds, false},
		{"parse RFC3339", "2022-01-21T12:00:34Z", RFC3339, false},
		{"parse RFC3339Nano", "2022-01-21T12:00:34.123456789Z", RFC3339Nano, false},
		{"parse DateTime", "2022-01-21 12:00:34", DateTime, false},
		{"parse DateOnly", "2022-01-21", DateOnly, false},
		{"parse TimeOnly", "12:00:34", TimeOnly, false},
		{"invalid unix", "not-a-number", UnixSeconds, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := conv.Parse(tt.input, tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse(%q, %v) error = %v, wantErr %v", tt.input, tt.format, err, tt.wantErr)
			}
		})
	}
}

func TestParseUnixSeconds(t *testing.T) {
	conv, _ := NewConverter("UTC")

	// Test specific timestamp: 2022-01-21 16:07:14 UTC
	result, err := conv.Parse("1642781234", UnixSeconds)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expected := time.Date(2022, 1, 21, 16, 7, 14, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("Parse(1642781234) = %v, want %v", result, expected)
	}
}

func TestConvert(t *testing.T) {
	conv, _ := NewConverter("UTC")

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"convert unix seconds", "1642781234", false},
		{"convert RFC3339", "2022-01-21T12:00:34Z", false},
		{"convert DateTime", "2022-01-21 12:00:34", false},
		{"invalid input", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := conv.Convert(tt.input, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if result.Original != tt.input {
					t.Errorf("Convert(%q).Original = %q, want %q", tt.input, result.Original, tt.input)
				}
				if result.UnixSeconds == 0 {
					t.Errorf("Convert(%q).UnixSeconds = 0, want non-zero", tt.input)
				}
			}
		})
	}
}

func TestConvertResult(t *testing.T) {
	conv, _ := NewConverter("UTC")

	result, err := conv.Convert("1642781234", nil)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// Verify all fields are populated
	if result.UnixSeconds != 1642781234 {
		t.Errorf("UnixSeconds = %d, want 1642781234", result.UnixSeconds)
	}
	if result.UnixMillis != 1642781234000 {
		t.Errorf("UnixMillis = %d, want 1642781234000", result.UnixMillis)
	}
	if result.UnixMicros != 1642781234000000 {
		t.Errorf("UnixMicros = %d, want 1642781234000000", result.UnixMicros)
	}
	if result.UnixNanos != 1642781234000000000 {
		t.Errorf("UnixNanos = %d, want 1642781234000000000", result.UnixNanos)
	}
	if result.DateTime != "2022-01-21 16:07:14" {
		t.Errorf("DateTime = %q, want %q", result.DateTime, "2022-01-21 16:07:14")
	}
	if result.DateOnly != "2022-01-21" {
		t.Errorf("DateOnly = %q, want %q", result.DateOnly, "2022-01-21")
	}
	if result.TimeOnly != "16:07:14" {
		t.Errorf("TimeOnly = %q, want %q", result.TimeOnly, "16:07:14")
	}
}

func TestAddTimeOffset(t *testing.T) {
	conv, _ := NewConverter("UTC")
	baseTime := time.Date(2022, 1, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		offset   string
		expected time.Time
		wantErr  bool
	}{
		{"add 1 day", "+1d", time.Date(2022, 1, 16, 12, 0, 0, 0, time.UTC), false},
		{"subtract 1 day", "-1d", time.Date(2022, 1, 14, 12, 0, 0, 0, time.UTC), false},
		{"add 1 week", "+1w", time.Date(2022, 1, 22, 12, 0, 0, 0, time.UTC), false},
		{"subtract 1 week", "-1w", time.Date(2022, 1, 8, 12, 0, 0, 0, time.UTC), false},
		{"add 1 month", "+1M", time.Date(2022, 2, 15, 12, 0, 0, 0, time.UTC), false},
		{"subtract 1 month", "-1M", time.Date(2021, 12, 15, 12, 0, 0, 0, time.UTC), false},
		{"add 1 year", "+1y", time.Date(2023, 1, 15, 12, 0, 0, 0, time.UTC), false},
		{"subtract 1 year", "-1y", time.Date(2021, 1, 15, 12, 0, 0, 0, time.UTC), false},
		{"add 2 hours", "+2h", time.Date(2022, 1, 15, 14, 0, 0, 0, time.UTC), false},
		{"subtract 30 minutes", "-30m", time.Date(2022, 1, 15, 11, 30, 0, 0, time.UTC), false},
		{"add 30 seconds", "+30s", time.Date(2022, 1, 15, 12, 0, 30, 0, time.UTC), false},
		{"empty offset", "", baseTime, false},
		{"invalid offset", "invalid", baseTime, true},
		{"invalid unit", "+1x", baseTime, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := conv.AddTimeOffset(baseTime, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddTimeOffset(%v, %q) error = %v, wantErr %v", baseTime, tt.offset, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !result.Equal(tt.expected) {
				t.Errorf("AddTimeOffset(%v, %q) = %v, want %v", baseTime, tt.offset, result, tt.expected)
			}
		})
	}
}

func TestParseTimeOffset(t *testing.T) {
	tests := []struct {
		name     string
		offset   string
		expected time.Duration
		wantErr  bool
	}{
		{"empty offset", "", 0, false},
		{"1 second", "+1s", time.Second, false},
		{"negative 1 second", "-1s", -time.Second, false},
		{"5 minutes", "+5m", 5 * time.Minute, false},
		{"2 hours", "+2h", 2 * time.Hour, false},
		{"1 day", "+1d", 24 * time.Hour, false},
		{"1 week", "+1w", 7 * 24 * time.Hour, false},
		{"1 month (30 days)", "+1M", 30 * 24 * time.Hour, false},
		{"1 year (365 days)", "+1y", 365 * 24 * time.Hour, false},
		{"invalid format", "invalid", 0, true},
		{"invalid unit", "+1x", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTimeOffset(tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimeOffset(%q) error = %v, wantErr %v", tt.offset, err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("ParseTimeOffset(%q) = %v, want %v", tt.offset, result, tt.expected)
			}
		})
	}
}

func TestTimezoneConversion(t *testing.T) {
	// Test converting same timestamp in different timezones
	convUTC, _ := NewConverter("UTC")
	convTaipei, _ := NewConverter("Asia/Taipei")

	input := "1642781234" // 2022-01-21 16:07:14 UTC

	resultUTC, err := convUTC.Convert(input, nil)
	if err != nil {
		t.Fatalf("Convert in UTC failed: %v", err)
	}

	resultTaipei, err := convTaipei.Convert(input, nil)
	if err != nil {
		t.Fatalf("Convert in Asia/Taipei failed: %v", err)
	}

	// Unix timestamps should be the same
	if resultUTC.UnixSeconds != resultTaipei.UnixSeconds {
		t.Errorf("Unix timestamps differ: UTC=%d, Taipei=%d", resultUTC.UnixSeconds, resultTaipei.UnixSeconds)
	}

	// DateTime should differ by 8 hours (Taipei is UTC+8)
	if resultUTC.DateTime == resultTaipei.DateTime {
		t.Errorf("DateTime should differ between timezones: UTC=%s, Taipei=%s", resultUTC.DateTime, resultTaipei.DateTime)
	}

	// Verify Taipei time is 8 hours ahead
	expectedTaipeiDateTime := "2022-01-22 00:07:14" // UTC 16:07:14 + 8 hours = 00:07:14 next day
	if resultTaipei.DateTime != expectedTaipeiDateTime {
		t.Errorf("Taipei DateTime = %q, want %q", resultTaipei.DateTime, expectedTaipeiDateTime)
	}
}

func TestGetLocalTimezone(t *testing.T) {
	tz := GetLocalTimezone()
	if tz == "" {
		t.Error("GetLocalTimezone() returned empty string")
	}
}

// Benchmark tests
func BenchmarkDetectFormat(b *testing.B) {
	conv, _ := NewConverter("UTC")
	inputs := []string{
		"1642781234",
		"2022-01-21T12:00:34Z",
		"2022-01-21 12:00:34",
		"2022-01-21",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, input := range inputs {
			conv.DetectFormat(input)
		}
	}
}

func BenchmarkConvert(b *testing.B) {
	conv, _ := NewConverter("UTC")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conv.Convert("1642781234", nil)
	}
}
