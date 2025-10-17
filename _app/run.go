package app

import (
	"context"
)

// пинг на апи /ready идет из корневого index.html (embeded)
// по таймеру проверяем время последнего пинга и если больше 3 сек закрываем приложение
func (a *app) Run(ctx context.Context, cancel context.CancelFunc) error {
	// a.readyPingLastTime = time.Now()
	// go func() {
	// 	a.reductor.RegisterGui(a.ReductorUpdaterHttp)
	// 	a.reductor.Run(ctx)
	// }()
	// go func() {
	// 	a.effects.Run(ctx)
	// }()
	// для wails убираем проверку
	// go func() {
	// 	ticker := time.NewTicker(time.Second)
	// 	defer ticker.Stop()

	// 	for range ticker.C {
	// 		dur := time.Since(a.readyPingLastTime).Seconds()
	// 		if dur > durationTimePingOut {
	// 			// a.Logger().Debugf("ping timeout duration %v mode %s", dur, domain.Mode)
	// 			if domain.Mode != "development" {
	// 				a.Logger().Errorf("ping not present app shutdown!")
	// 				cancel()
	// 			}
	// 		}
	// 	}
	// }()
	// убираем запуск приложения через браузер в продакшн
	// if domain.Mode == "development" {
	// 	port := a.Configuration().HostPort
	// 	host := a.Configuration().Hostname
	// 	proto := `http://`
	// 	url := fmt.Sprintf("%s%s:%s", proto, host, port)
	// 	a.Open(url)
	// }
	<-ctx.Done()
	return nil
}
