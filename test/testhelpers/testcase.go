package testhelpers

import "testing"

// TestCase represents a single test case for table-driven tests
type TestCase struct {
	Name        string
	Input       interface{}
	Expected    interface{}
	ShouldError bool
	ErrorMsg    string
	Setup       func() error
	Cleanup     func() error
}

// RunTests executes a series of test cases
func RunTests(t *testing.T, cases []TestCase, testFunc func(*testing.T, TestCase)) {
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.Setup != nil {
				if err := tc.Setup(); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}
			defer func() {
				if tc.Cleanup != nil {
					if err := tc.Cleanup(); err != nil {
						t.Errorf("Cleanup failed: %v", err)
					}
				}
			}()

			testFunc(t, tc)
		})
	}
}
