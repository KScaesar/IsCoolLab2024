package pkg

import (
	"time"
)

type TimeFunc interface {
	Now() time.Time
	Sleep(d time.Duration)
}

// NewMockTimeFunc
// This can be useful for testing scenarios that involve time-sensitive operations without
// actually manipulating the system clock.
//
// Parameters:
// RFC3339 (string): The time value in RFC3339 format that will be used for simulation.
//
// Example usage:
//
//	NewMockTimeFunc("2023-08-19T12:00:00Z")
//	NewMockTimeFunc("2023-08-19T20:00:00+08:00")
func NewMockTimeFunc(RFC3339 string) MockTimeFunc {
	t, err := time.Parse(time.RFC3339, RFC3339)
	if err != nil {
		panic(err)
	}

	return MockTimeFunc{
		now: t,
	}
}

type MockTimeFunc struct {
	now time.Time
}

func (f MockTimeFunc) Now() time.Time {
	return f.now
}

func (f *MockTimeFunc) Sleep(d time.Duration) {
	f.now = f.now.Add(d)
}
