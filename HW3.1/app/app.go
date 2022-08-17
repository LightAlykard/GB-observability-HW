package app

import (
	"io"

	"github.com/LightAlykard/GB-observability-HW/HW3/handler"
	"github.com/LightAlykard/GB-observability-HW/HW3/l"

	//"github.com/LightAlykard/GB-observability-HW/HW3/s"
	"github.com/LightAlykard/GB-observability-HW/HW3/store"
	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type App struct {
	logger *zap.Logger
	tracer opentracing.Tracer
}

func (a *App) Init() (io.Closer, error) {
	//ctx := context.Background()
	// Предустановленный конфиг. Можно выбрать
	// NewProduction/NewDevelopment/NewExample или создать свой
	// Production - уровень логгирования InfoLevel, формат вывода: json
	// Development - уровень логгирования DebugLevel, формат вывода: console
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	defer func() { _ = logger.Sync() }()
	// Трейсер
	// Можно "захардкодить" при инициализации
	//tracer, closer := l.InitJaeger("App", "jaeger:6831", logger)
	// Или использовать переменные окружения
	tracer, closer := l.InitJaeger(logger)

	a.logger = logger
	a.tracer = tracer

	return closer, nil
}

func (a *App) Serve() error {
	//Sentry error handler
	//s.NewSentryLogger()

	//Initialize Stores
	articleStore, err := store.NewArticleStore(a.logger, a.tracer)
	parseErr(err)

	//Initialize Handlers
	articleHandler := handler.NewArticleHandler(articleStore, a.logger, a.tracer)
	panicHandler := handler.NewPanicHandler(a.logger, a.tracer)

	//Initialize Router and add Middleware
	router := gin.Default()
	router.Use(nice.Recovery(panicHandler.RecoveryHandler))
	router.LoadHTMLFiles("template/error.tpl")

	//Routes
	router.GET("/article/id/:id", articleHandler.Id)
	router.POST("/article/add", articleHandler.Add)
	router.POST("/article/search", articleHandler.Search)
	router.GET("/panic", panicHandler.Panic)
	router.POST("/log/add", panicHandler.Log)

	// Start serving the application
	return router.Run()
}

func parseErr(err error) {
	if err != nil {
		l.F(err)
	}
	l.L("Serve start")
}
