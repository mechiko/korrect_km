package zaplog

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

var Logger *zap.Logger

type ZapLog struct {
	logs map[LogName]*zap.Logger
}

func (z *ZapLog) GetLogger(name string) (*zap.Logger, error) {
	if isValidLogName(name) {
		if log, ok := z.logs[LogName(name)]; ok {
			return log, nil
		} else {
			return nil, fmt.Errorf("%s is not created", name)
		}
	} else {
		return nil, fmt.Errorf("%s not valid", name)
	}
}

func (z *ZapLog) Shutdown() {
	fmt.Println("zap log shutdown")
	for _, log := range z.logs {
		log.Sync()
	}
}

//	var logsOutConfig = map[string]zaplog.LogConfig{
//		"logger": {
//			ErrorOutputPaths: []string{"stdout", filepath.Join(cfg.LogPath(), config.Name)},
//			Debug:            debug,
//			Console:          true,
//		},
//		"echo": {
//			ErrorOutputPaths: []string{"stdout", filepath.Join(cfg.LogPath(), "echo")},
//			Debug:            debug,
//			Console:          false,
//		},
//		"reductor": {
//			ErrorOutputPaths: []string{"stdout", filepath.Join(cfg.LogPath(), "reductor")},
//			Debug:            debug,
//			Console:          true,
//		},
//	}
func New(outConfig map[string]LogConfig) (*ZapLog, error) {
	// проверяем мапу настройки логов
	for key := range outConfig {
		if !isValidLogName(key) {
			return nil, fmt.Errorf("wrong name %s", key)
		}
	}
	z := &ZapLog{
		logs: make(map[LogName]*zap.Logger),
	}
	err := z.init(outConfig)
	if err != nil {
		return nil, fmt.Errorf("init zap logger %v", err)
	}
	return z, nil
}

func (z *ZapLog) Run(ctx context.Context) error {
	// ожидаем завершения контекста
	<-ctx.Done()
	fmt.Println("zaplog receive ctx shutdown")
	z.Shutdown()
	return nil
}

func (z *ZapLog) init(outConfig map[string]LogConfig) (err error) {
	for key, output := range outConfig {
		if isValidLogName(key) {
			lg, err := createLogger(output.ErrorOutputPaths, output.Debug, output.Console)
			if err != nil {
				return fmt.Errorf("name %s %w", key, err)
			}
			z.logs[LogName(key)] = lg
		} else {
			return fmt.Errorf("wrong name for logger %s", key)
		}
	}
	// для совместимости основной логер пропишем в глобальную переменную
	if l, ok := z.logs[LogNameLogger]; ok {
		Logger = l
	}
	return nil
}
