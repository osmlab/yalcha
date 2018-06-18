package osm

import (
	"testing"
	"time"
)

func TestTimeString(t *testing.T) {
	newYork, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("invalid timezone: %v", err)
	}

	cases := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "iso 8601 format",
			time:     time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "2012-01-01T00:00:00Z",
		},
		{
			name:     "always UTC",
			time:     time.Date(2012, 1, 1, 0, 0, 0, 0, newYork),
			expected: "2012-01-01T05:00:00Z",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tosm := Time(tc.time)
			if v := tosm.String(); v != tc.expected {
				t.Errorf("incorrect format: %v", v)
			}
		})
	}
}
