package internal

import (
	"testing"
	"time"

	"github.com/akyrey/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 04, 29, 22, 9, 0, 0, time.UTC),
			want: "29 Apr 2024 at 22:09",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 04, 29, 22, 9, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "29 Apr 2024 at 21:09",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := humanDate(tt.tm)
			assert.Equal(t, got, tt.want)
		})
	}
}
