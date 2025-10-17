package zaplog

import (
	"korrectkm/config"
	"os"
	"path"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var onlyOnce sync.Once

var Logger *zap.SugaredLogger
var Reductor *zap.Logger
var EchoSugar *zap.Logger
var TrueSugar *zap.SugaredLogger

var LogPath = ""
var LogName = "log"

func onRun(path string, name string) {
	LogPath = path
	if name != "" {
		LogName = name
	}
	onlyOnce.Do(func() {
		initLogger()
	})
}

func OnShutdown() {
	Logger.Sync()
	Reductor.Sync()
	EchoSugar.Sync()
	TrueSugar.Sync()
}

// возврат только после прерывания контекста
func Run(path string, name string) error {
	onRun(path, name)
	return nil
}

func initLogger() {
	var defaultLogLevel zapcore.Level

	logname := config.Name + ".log"
	logpath := path.Join(LogPath, logname)
	configLogger := zap.NewProductionEncoderConfig()
	// configLogger.EncodeLevel = zapcore.CapitalColorLevelEncoder
	configLogger.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	fileEncoderLogger := zapcore.NewConsoleEncoder(configLogger)
	consoleEncoder := zapcore.NewConsoleEncoder(configLogger)
	// writer := zapcore.AddSync(wLoggerRotate)
	fileLogger, _ := os.OpenFile(logpath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(fileLogger)
	defaultLogLevel = zapcore.DebugLevel

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoderLogger, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	Logger = zap.New(core, zap.AddCaller()).Sugar()

	configReductor := zap.NewProductionEncoderConfig()
	configReductor.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	logpathReductor := path.Join(LogPath, "reductor.log")
	reductorLogFile, _ := os.OpenFile(logpathReductor, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	fileReductorEncoder := zapcore.NewConsoleEncoder(configReductor)
	writerReductor := zapcore.AddSync(reductorLogFile)
	core2 := zapcore.NewTee(
		zapcore.NewCore(fileReductorEncoder, writerReductor, defaultLogLevel),
	)
	Reductor = zap.New(core2, zap.AddCaller())

	configEcho := zap.NewProductionEncoderConfig()
	configEcho.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	logpathEcho := path.Join(LogPath, "echo.log")
	echoLogFile, _ := os.OpenFile(logpathEcho, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	fileEchoEncoder := zapcore.NewConsoleEncoder(configEcho)
	writerEcho := zapcore.AddSync(echoLogFile)
	core5 := zapcore.NewTee(
		zapcore.NewCore(fileEchoEncoder, writerEcho, defaultLogLevel),
	)
	EchoSugar = zap.New(core5, zap.AddCaller())

	configTrue := zap.NewProductionEncoderConfig()
	configTrue.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	logpathTrue := path.Join(LogPath, "true.log")
	trueLogFile, _ := os.OpenFile(logpathTrue, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	fileTrueEncoder := zapcore.NewConsoleEncoder(configTrue)
	writerTrue := zapcore.AddSync(trueLogFile)
	core6 := zapcore.NewTee(
		zapcore.NewCore(fileTrueEncoder, writerTrue, defaultLogLevel),
	)
	TrueSugar = zap.New(core6, zap.AddCaller()).Sugar()

}
