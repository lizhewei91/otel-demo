package service

import "github.com/lizw91/otel-demo/internal/model"

var Healthy = sHealthy{}

type sHealthy struct{}

// Status system status
func (h *sHealthy) Status() *model.HealthySrvStatusRes {
	status := &model.HealthySrvStatusRes{
		Status: "OK",
	}

	return status
}
