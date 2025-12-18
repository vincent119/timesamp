package i18n

import (
	"testing"
)

func TestInit(t *testing.T) {
	// Init should not panic
	Init()
}

func TestSetLanguage(t *testing.T) {
	tests := []struct {
		name string
		lang string
	}{
		{"English", "en"},
		{"Traditional Chinese", "zh-TW"},
		{"Simplified Chinese", "zh-CN"},
		{"Japanese", "ja"},
		{"Unknown language falls back", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// SetLanguage should not panic
			SetLanguage(tt.lang)
		})
	}
}

func TestT(t *testing.T) {
	// Initialize i18n first
	Init()
	SetLanguage("en")

	tests := []struct {
		name   string
		msgID  string
		notNil bool
	}{
		{"existing message", "cmd.root.short", true},
		{"non-existing message returns ID", "non.existing.key", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := T(tt.msgID)
			if tt.notNil && result == "" {
				t.Errorf("T(%q) returned empty string", tt.msgID)
			}
		})
	}
}

func TestLanguageSwitching(t *testing.T) {
	Init()

	// Test that switching languages works
	SetLanguage("en")
	enResult := T("cmd.root.short")

	SetLanguage("zh-TW")
	zhTWResult := T("cmd.root.short")

	SetLanguage("ja")
	jaResult := T("cmd.root.short")

	// Results should all be non-empty
	if enResult == "" {
		t.Error("English translation is empty")
	}
	if zhTWResult == "" {
		t.Error("Traditional Chinese translation is empty")
	}
	if jaResult == "" {
		t.Error("Japanese translation is empty")
	}

	// Log the results for debugging
	t.Logf("en: %s", enResult)
	t.Logf("zh-TW: %s", zhTWResult)
	t.Logf("ja: %s", jaResult)
}
