package builder

import (
	"context"
	"fmt"
	"runtime"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// ProcessPackagesParallel processes packages concurrently with bounded parallelism.
// Concurrency is limited to GOMAXPROCS * 2. Processing stops on first error,
// but other goroutines may complete before the error is returned.
//
// The processor function is called once per package. If any processor returns an error,
// ProcessPackagesParallel returns that error after all goroutines complete.
//
// Context cancellation will stop new processors from starting and cause already-running
// processors to be interrupted if they check ctx.Err().
func ProcessPackagesParallel(ctx context.Context, packages []string, processor func(context.Context, string) error) error {
	if len(packages) == 0 {
		return nil
	}

	// Set concurrency limit based on CPU count
	maxConcurrency := runtime.GOMAXPROCS(0) * 2
	if maxConcurrency < 1 {
		maxConcurrency = 1
	}

	// Use errgroup for coordinated error handling
	var g errgroup.Group

	// Create semaphore to limit concurrency
	sem := semaphore.NewWeighted(int64(maxConcurrency))

	for _, pkg := range packages {
		pkg := pkg // Capture loop variable

		g.Go(func() error {
			// Acquire semaphore with context for cancellation
			if err := sem.Acquire(ctx, 1); err != nil {
				return err // context.Canceled or context.DeadlineExceeded
			}
			defer sem.Release(1)

			// Check if context was cancelled before processing
			if ctx.Err() != nil {
				return ctx.Err()
			}

			// Process package with context
			if err := processor(ctx, pkg); err != nil {
				return fmt.Errorf("process package %s: %w", pkg, err)
			}

			return nil
		})
	}

	// Wait for all goroutines and collect errors
	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
