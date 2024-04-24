package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		span := trace.SpanFromContext(ctx)

		c.Next()
		statusCode := c.Writer.Status()
		logrus.WithFields(logrus.Fields{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"trace-id": span.SpanContext().TraceID(),
			"code":     statusCode,
			"latency":  time.Since(start).String(),
			"sampled":  span.SpanContext().IsSampled(),
		}).Info(http.StatusText(statusCode))
	}
}
