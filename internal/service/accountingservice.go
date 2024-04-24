package service

import (
	"context"

	"github.com/lizw91/otel-demo/internal/model"
)

// 对外公开的实例化对象名
var AccountingService = sAccountingServiceSvc{}

// 数据结构名, 命名规则: s + 业务邻域名
type sAccountingServiceSvc struct{}

func (s *sAccountingServiceSvc) GetUserName(ctx context.Context, input *model.GetUserNameInput) (*model.GetUserNameOutput, error) {

	return nil, nil
}
