package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	apiv1 "github.com/lizw91/otel-demo/api/v1"
	"github.com/lizw91/otel-demo/internal/library/response"
	"github.com/lizw91/otel-demo/internal/model"
	"github.com/lizw91/otel-demo/internal/service"
)

// 对外公开的实例化对象名
var AccountintService = hAccountintService{}

// 数据结构名, 命名规则: h + 业务邻域名
type hAccountintService struct{}

// GetUserName
// @Summary 获取用户名称
// @Description 获取用户名称
// @Tags User
// @Produce json
// @Success 200 {object} apiv1.UserInfoRes
// @Failure 400 {object} response.ErrResponse
// @Router /user/:id [get]
func (h *hAccountintService) GetUserName(c *gin.Context) {
	req := &apiv1.GetUserNameReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.ParamErrResponse(c, err)
		return
	}

	input := &model.GetUserNameInput{}
	_ = copier.Copy(input, req)
	output, err := service.AccountingService.GetUserName(c, input)
	if err != nil {
		response.SvcErrResponse(c, err)
		return
	}

	res := &apiv1.GetUserNameRes{}
	_ = copier.Copy(res, output)
	response.SuccessResponse(c, res)
}
