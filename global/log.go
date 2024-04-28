package global

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	logging "github.com/lizw91/otel-demo/pkg/log"
)

var LoggerF logging.Factory

var once sync.Once

func InitLoggerFactory(serviceName string) {
	once.Do(func() {
		zapOptions := []zap.Option{
			zap.AddStacktrace(zapcore.FatalLevel),
			zap.AddCallerSkip(1),
		}

		logger, _ := zap.NewDevelopment(zapOptions...)
		zapLogger := logger.With(zap.String("service", serviceName))
		LoggerF = logging.NewFactory(zapLogger)
	})
}
