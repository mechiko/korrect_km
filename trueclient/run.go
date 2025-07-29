package trueclient

import (
	"context"
)

// errgroup
func (t *trueClient) Run(ctx context.Context) error {
	t.Logger().Info("trueclient service starting")
	defer t.Logger().Info("trueclient service stopped")

	// Add actual service logic here
	// e.g., periodic tasks, event processing, etc.

	<-ctx.Done()

	// Add cleanup logic here if needed
	t.Logger().Info("trueclient service shutting down gracefully")
	return nil
}
