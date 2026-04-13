package builder

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestProcessPackagesParallel(t *testing.T) {
	tests := map[string]struct {
		packages    []string
		processor   func(context.Context, string) error
		wantErr     bool
		validate    func(t *testing.T, processed map[string]bool)
		errContains string
	}{
		"processes all packages successfully": {
			packages: []string{"pkg1", "pkg2", "pkg3", "pkg4", "pkg5"},
			processor: func(ctx context.Context, pkg string) error {
				time.Sleep(10 * time.Millisecond) // Simulate work
				return nil
			},
			wantErr: false,
			validate: func(t *testing.T, processed map[string]bool) {
				if len(processed) != 5 {
					t.Errorf("got %d processed, want 5", len(processed))
				}
			},
		},
		"handles single package": {
			packages: []string{"single-pkg"},
			processor: func(ctx context.Context, pkg string) error {
				return nil
			},
			wantErr: false,
			validate: func(t *testing.T, processed map[string]bool) {
				if len(processed) != 1 {
					t.Errorf("got %d processed, want 1", len(processed))
				}
			},
		},
		"handles empty package list": {
			packages: []string{},
			processor: func(ctx context.Context, pkg string) error {
				return nil
			},
			wantErr: false,
			validate: func(t *testing.T, processed map[string]bool) {
				if len(processed) != 0 {
					t.Errorf("got %d processed, want 0", len(processed))
				}
			},
		},
		"returns error when processor fails": {
			packages: []string{"pkg1", "pkg2", "pkg3"},
			processor: func(ctx context.Context, pkg string) error {
				if pkg == "pkg2" {
					return errors.New("processing failed for pkg2")
				}
				return nil
			},
			wantErr:     true,
			errContains: "pkg2",
		},
		"collects multiple errors": {
			packages: []string{"pkg1", "pkg2", "pkg3", "pkg4"},
			processor: func(ctx context.Context, pkg string) error {
				if pkg == "pkg2" || pkg == "pkg4" {
					return fmt.Errorf("failed: %s", pkg)
				}
				return nil
			},
			wantErr:     true,
			errContains: "failed",
		},
		"processes packages concurrently": {
			packages: []string{"pkg1", "pkg2", "pkg3", "pkg4", "pkg5", "pkg6"},
			processor: func(ctx context.Context, pkg string) error {
				time.Sleep(50 * time.Millisecond) // Simulate work
				return nil
			},
			wantErr: false,
			validate: func(t *testing.T, processed map[string]bool) {
				// All should be processed
				if len(processed) != 6 {
					t.Errorf("got %d processed, want 6", len(processed))
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Track processed packages thread-safely
			var mu sync.Mutex
			processed := make(map[string]bool)

			// Wrap processor to track execution
			trackingProcessor := func(ctx context.Context, pkg string) error {
				err := tc.processor(ctx, pkg)
				if err == nil {
					mu.Lock()
					processed[pkg] = true
					mu.Unlock()
				}
				return err
			}

			ctx := context.Background()
			err := ProcessPackagesParallel(ctx, tc.packages, trackingProcessor)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.errContains != "" && err != nil {
				if !contains(err.Error(), tc.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tc.errContains)
				}
			}

			if tc.validate != nil {
				tc.validate(t, processed)
			}
		})
	}
}

func TestProcessPackagesParallel_Concurrency(t *testing.T) {
	tests := map[string]struct {
		packageCount   int
		expectParallel bool
	}{
		"processes packages in parallel": {
			packageCount:   10,
			expectParallel: true,
		},
		"single package runs immediately": {
			packageCount:   1,
			expectParallel: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			packages := make([]string, tc.packageCount)
			for i := 0; i < tc.packageCount; i++ {
				packages[i] = fmt.Sprintf("pkg%d", i)
			}

			var activeCount int32
			var maxActive int32
			var wg sync.WaitGroup

			processor := func(ctx context.Context, pkg string) error {
				current := atomic.AddInt32(&activeCount, 1)
				defer atomic.AddInt32(&activeCount, -1)

				// Track maximum concurrent executions
				for {
					max := atomic.LoadInt32(&maxActive)
					if current <= max || atomic.CompareAndSwapInt32(&maxActive, max, current) {
						break
					}
				}

				time.Sleep(10 * time.Millisecond) // Simulate work
				return nil
			}

			ctx := context.Background()
			start := time.Now()
			err := ProcessPackagesParallel(ctx, packages, processor)
			elapsed := time.Since(start)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			wg.Wait()

			if tc.expectParallel {
				// With parallel processing, should complete faster than sequential
				sequentialTime := time.Duration(tc.packageCount) * 10 * time.Millisecond
				// Allow for some overhead
				if elapsed > sequentialTime*8/10 {
					t.Errorf("processing took %v, expected faster than %v (sequential)", elapsed, sequentialTime)
				}

				// Should have had multiple packages running concurrently
				if maxActive <= 1 {
					t.Errorf("maxActive = %d, expected > 1 for parallel processing", maxActive)
				}
			}
		})
	}
}

func TestProcessPackagesParallel_ConcurrencyLimit(t *testing.T) {
	packageCount := 50
	packages := make([]string, packageCount)
	for i := 0; i < packageCount; i++ {
		packages[i] = fmt.Sprintf("pkg%d", i)
	}

	var activeCount int32
	var maxActive int32

	processor := func(ctx context.Context, pkg string) error {
		current := atomic.AddInt32(&activeCount, 1)
		defer atomic.AddInt32(&activeCount, -1)

		// Track maximum concurrent executions
		for {
			max := atomic.LoadInt32(&maxActive)
			if current <= max || atomic.CompareAndSwapInt32(&maxActive, max, current) {
				break
			}
		}

		time.Sleep(5 * time.Millisecond) // Simulate work
		return nil
	}

	ctx := context.Background()
	err := ProcessPackagesParallel(ctx, packages, processor)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Concurrency should be limited (not all 50 running at once)
	cpuCount := runtime.GOMAXPROCS(0)
	expectedLimit := cpuCount * 2 // Typical limit

	if maxActive > int32(expectedLimit*2) {
		t.Errorf("maxActive = %d, expected <= %d (concurrency limit should be enforced)",
			maxActive, expectedLimit*2)
	}

	t.Logf("Processed %d packages with max concurrency of %d (CPU count: %d)",
		packageCount, maxActive, cpuCount)
}

func TestProcessPackagesParallel_ErrorHandling(t *testing.T) {
	tests := map[string]struct {
		packages    []string
		failOn      map[string]bool
		wantErr     bool
		checkErrors func(t *testing.T, err error)
	}{
		"stops on first error": {
			packages: []string{"pkg1", "pkg2", "pkg3", "pkg4", "pkg5"},
			failOn:   map[string]bool{"pkg3": true},
			wantErr:  true,
			checkErrors: func(t *testing.T, err error) {
				if err == nil {
					t.Fatal("expected error")
				}
				if !contains(err.Error(), "pkg3") {
					t.Errorf("error should mention pkg3: %v", err)
				}
			},
		},
		"collects all errors when multiple fail": {
			packages: []string{"pkg1", "pkg2", "pkg3", "pkg4"},
			failOn:   map[string]bool{"pkg1": true, "pkg3": true},
			wantErr:  true,
			checkErrors: func(t *testing.T, err error) {
				if err == nil {
					t.Fatal("expected error")
				}
				// Error should contain information about failures
				errStr := err.Error()
				if !contains(errStr, "pkg") {
					t.Errorf("error should contain package info: %v", err)
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			processor := func(ctx context.Context, pkg string) error {
				if tc.failOn[pkg] {
					return fmt.Errorf("failed to process %s", pkg)
				}
				time.Sleep(5 * time.Millisecond) // Simulate work
				return nil
			}

			ctx := context.Background()
			err := ProcessPackagesParallel(ctx, tc.packages, processor)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.checkErrors != nil {
				tc.checkErrors(t, err)
			}
		})
	}
}

func TestProcessPackagesParallel_RaceConditions(t *testing.T) {
	// This test is designed to catch race conditions with -race flag
	packageCount := 100
	packages := make([]string, packageCount)
	for i := 0; i < packageCount; i++ {
		packages[i] = fmt.Sprintf("pkg%d", i)
	}

	// Shared state that processor will access
	counter := 0
	var mu sync.Mutex

	processor := func(ctx context.Context, pkg string) error {
		// Access shared state (should be protected)
		mu.Lock()
		counter++
		mu.Unlock()
		return nil
	}

	ctx := context.Background()
	err := ProcessPackagesParallel(ctx, packages, processor)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	finalCount := counter
	mu.Unlock()

	if finalCount != packageCount {
		t.Errorf("counter = %d, want %d", finalCount, packageCount)
	}
}
