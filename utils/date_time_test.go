package utils

import (
	"testing"
	"time"
)

func TestCurrentTimestamp_UsesConfig(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Kuala_Lumpur")
	if err != nil {
		t.Skip("Asia/Kuala_Lumpur timezone not available")
	}

	ts := CurrentTimestamp(loc)
	if ts == "" {
		t.Error("expected non-empty timestamp")
	}

	tsUTC := CurrentTimestamp(time.UTC)
	if tsUTC == "" {
		t.Error("expected non-empty UTC timestamp")
	}
}
