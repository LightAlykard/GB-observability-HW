package main

import (
	"fmt"
	"io"
	"log"

	"github.com/LightAlykard/GB-observability-HW/HW3/app"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

type zapWrapper struct {
	logger *zap.Logger
}

// Error logs a message at error priority
func (w *zapWrapper) Error(msg string) {
	w.logger.Error(msg)
}

// Infof logs a message at info priority
func (w *zapWrapper) Infof(msg string, args ...interface{}) {
	w.logger.Sugar().Infof(msg, args...)
}
func initJaeger(service string, logger *zap.Logger) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(&zapWrapper{logger: logger}))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func main() {
	// Предустановленный конфиг. Можно выбрать
	// NewProduction/NewDevelopment/NewExample или создать свой
	// Production - уровень логгирования InfoLevel, формат вывода: json
	// Development - уровень логгирования DebugLevel, формат вывода: console
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = logger.Sync() }()
	// Трейсер
	tracer, closer := initJaeger("example", logger)
	defer closer.Close()

	a := app.App{}
	if err := a.Init(logger, tracer); err != nil {
		log.Fatal(err)
	}

	//a := app.App{}
	// if closer, err := a.Init(); err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	defer closer.Close()
	// }

	if err := a.Serve(); err != nil {
		log.Fatal(err)
	}
}
