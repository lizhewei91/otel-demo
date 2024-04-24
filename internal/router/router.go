package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/lizw91/otel-demo/api"
	handlerv1 "github.com/lizw91/otel-demo/api/v1/handler"
)

func Routers() http.Handler {
	engine := gin.New()
	engine.Use(
		otelgin.Middleware("accounting-service", otelgin.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/probe"
		})),
	)
	engine.GET("/user/:id", handlerv1.AccountintService.GetUserName)
	engine.GET("/probe", api.Healthy.Status)
	return engine
}
