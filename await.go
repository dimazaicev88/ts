package ts

import (
	"context"
	"fmt"
	"time"
)

func Until(ctx context.Context, timeout time.Duration, interval time.Duration, msg string, condition func() (completed bool, err error)) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	result, err := condition()
	if err != nil {
		return err
	}

	if result {
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%s: %w", msg, ctx.Err())
		case <-ticker.C:
			result, err = condition()
			if err != nil {
				return err
			}

			if result {
				return nil
			}
		}
	}
}
