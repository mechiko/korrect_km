package app

import "context"

// startup is called at application startup WAILS
func (a *app) Startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

func (a *app) BeforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *app) Shutdown(ctx context.Context) {
	// Perform your teardown here
}
