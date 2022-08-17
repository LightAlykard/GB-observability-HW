package handler

import (
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type PanicHandler struct {
	logger *zap.Logger
	tracer opentracing.Tracer
}

func NewPanicHandler(l *zap.Logger, t opentracing.Tracer) PanicHandler {
	return PanicHandler{
		logger: l,
		tracer: t,
	}
}

func (h PanicHandler) Panic(c *gin.Context) {
	h.logger.Info("PanicHandler.Panic called", zap.Field{Key: "method", String: c.Request.Method, Type: zapcore.StringType})
	panic("Panic at path /panic")
}

func (h *PanicHandler) RecoveryHandler(c *gin.Context, err interface{}) {
	sentry.CaptureException(fmt.Errorf("panic: %v", err))
	span, _ := opentracing.StartSpanFromContextWithTracer(c, h.tracer,
		"PanicHandler.Panic")
	defer span.Finish()

	span.SetTag("method", c.Request.Method)
	span.SetTag("params", c.Params)
	h.logger.Warn("Recovered from panic",
		zap.Error(fmt.Errorf("error: %v", err)),
	)
	span.LogFields(
		log.Error(fmt.Errorf("panic: %v", err)),
	)
	c.HTML(500, "error.tpl", gin.H{
		"title": "Internal server error",
		"err":   err,
	})
}

func (h PanicHandler) Log(c *gin.Context) {
	sentry.CaptureMessage("log is sent")
	span, _ := opentracing.StartSpanFromContextWithTracer(c, h.tracer,
		"PanicHandler.Log")
	defer span.Finish()
	span.SetTag("method", c.Request.Method)
	span.SetTag("params", c.Params)
	h.logger.Info("log message",
		zap.String("info", "add log"),
	)
	span.LogFields(
		log.String("add log", "info"),
	)
	c.JSON(http.StatusOK, gin.H{"msg": "sent to sentry"})
}
