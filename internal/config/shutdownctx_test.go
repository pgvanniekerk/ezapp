package config

import (
	"os"
	"testing"
	"time"
)

func TestShutdownCtx(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		envValue      string
		expectedError bool
		expectedTime  int // in seconds
	}{
		{
			name:          "default value",
			envValue:      "",
			expectedError: false,
			expectedTime:  15,
		},
		{
			name:          "custom value",
			envValue:      "30",
			expectedError: false,
			expectedTime:  30,
		},
		{
			name:          "invalid value",
			envValue:      "not-a-number",
			expectedError: true,
			expectedTime:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variable
			if tc.envValue != "" {
				os.Setenv("EZAPP_SHUTDOWN_TIMEOUT", tc.envValue)
				defer os.Unsetenv("EZAPP_SHUTDOWN_TIMEOUT")
			} else {
				os.Unsetenv("EZAPP_SHUTDOWN_TIMEOUT")
			}

			// Call the function
			ctx, err := ShutdownCtx()

			// Check error
			if tc.expectedError && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}

			// If no error, check deadline
			if !tc.expectedError {
				deadline, ok := ctx.Deadline()
				if !ok {
					t.Errorf("context should have a deadline")
				}

				// Calculate expected deadline (with some tolerance for test execution time)
				expectedDeadline := time.Now().Add(time.Duration(tc.expectedTime) * time.Second)
				tolerance := 100 * time.Millisecond

				diff := deadline.Sub(expectedDeadline)
				if diff < -tolerance || diff > tolerance {
					t.Errorf("expected deadline around %v but got %v (diff: %v)",
						expectedDeadline, deadline, diff)
				}
			}
		})
	}
}
