package api

import (
	"github.com/gin-gonic/gin"

	"github.com/lizw91/otel-demo/internal/library/response"
	"github.com/lizw91/otel-demo/internal/service"
)

var Healthy = healthyAPI{}

type healthyAPI struct{}

func (h *healthyAPI) Status(c *gin.Context) {
	status := service.Healthy.Status()

	response.SuccessResponse(c, status)
}
