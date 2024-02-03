package date

import (
	"testing"
	"time"
)

func TestTimeStampSlug(t *testing.T) {
	tests := []struct {
		input    time.Time
		expected string
	}{
		{time.Date(2022, time.January, 15, 10, 30, 45, 0, time.UTC), "2022-01-15_10-30-45"},
		{time.Date(2022, time.February, 25, 18, 45, 12, 0, time.UTC), "2022-02-25_18-45-12"},
		{time.Date(2023, time.March, 8, 5, 0, 0, 0, time.UTC), "2023-03-08_05-00-00"},
		{time.Date(2023, time.April, 1, 12, 15, 30, 0, time.UTC), "2023-04-01_12-15-30"},
		{time.Date(2024, time.May, 20, 23, 59, 59, 0, time.UTC), "2024-05-20_23-59-59"},
		{time.Date(2024, time.June, 10, 3, 5, 7, 0, time.UTC), "2024-06-10_03-05-07"},
		{time.Date(2025, time.July, 5, 15, 30, 0, 0, time.UTC), "2025-07-05_15-30-00"},
		{time.Date(2025, time.August, 12, 6, 45, 59, 0, time.UTC), "2025-08-12_06-45-59"},
		{time.Date(2026, time.September, 18, 1, 20, 10, 0, time.UTC), "2026-09-18_01-20-10"},
		{time.Date(2026, time.October, 30, 9, 12, 55, 0, time.UTC), "2026-10-30_09-12-55"},
		// False detected cases
		{time.Date(2022, time.January, 15, 10, 30, 45, 0, time.FixedZone("UTC-5", -5*60*60)), "2022-01-15_10-30-45"},
		{time.Date(2022, time.February, 25, 18, 45, 12, 0, time.FixedZone("UTC+3", 3*60*60)), "2022-02-25_18-45-12"},
		{time.Date(2023, time.March, 8, 5, 0, 0, 0, time.FixedZone("UTC-6", -6*60*60)), "2023-03-08_05-00-00"},
		{time.Date(2023, time.April, 1, 12, 15, 30, 0, time.FixedZone("UTC+2", 2*60*60)), "2023-04-01_12-15-30"},
		{time.Date(2024, time.May, 20, 23, 59, 59, 0, time.FixedZone("UTC-7", -7*60*60)), "2024-05-20_23-59-59"},
	}

	for _, test := range tests {
		result := TimeStampSlug(test.input)
		if result != test.expected {
			t.Errorf("TimeStampSlug(%v) returned %s, expected %s", test.input, result, test.expected)
		}
	}
}
