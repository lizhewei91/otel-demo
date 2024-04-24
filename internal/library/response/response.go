package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	successCode    = "SUCCESS"
	serverErrCode  = "SERVER_ERR"
	authErrCode    = "AUTH_ERR"
	paramErrCode   = "PARAM_ERR"
	notFindErrCode = "NOT_FIND_ERR"
)

// I18nLang 多语言 struct
type I18nLang struct {
	ZhHans string `json:"zh_hans" bson:"zh_hans"`
	En     string `json:"en" bson:"en"`
}

// CommonResponse 通用返回结构
type CommonResponse struct {
	Code string `json:"code"` // 代码
}

// SuccessResponse 成功调用 Response
func SuccessResponse(c *gin.Context, data interface{}) {
	if data != nil {
		c.JSON(http.StatusOK, data)
	} else {
		res := CommonResponse{
			Code: successCode,
		}
		c.JSON(http.StatusOK, res)
	}
}

// CreatedResponse 创建成功 Response, POST或者PUT请求创建一个新的资源完成
func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// NoContentResponse 无返回内容 Response
func NoContentResponse(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// SvcErrResponse 业务错误 Response
func SvcErrResponse(c *gin.Context, err error) {
	var (
		code    = serverErrCode
		message string
		detail  interface{}
	)

	if errData, ok := err.(*ErrResponse); ok {
		code = errData.Code
		message = errData.Message
		detail = errData.Data
	} else {
		message = err.Error()
	}

	res := &ErrResponse{
		Code:    code,
		Message: message,
		Data:    detail,
	}

	c.JSON(http.StatusBadRequest, res)
}

// AuthErrResponse 认证错误 Response
func AuthErrResponse(c *gin.Context, err error) {
	res := &ErrResponse{
		Code:    authErrCode,
		Message: err.Error(),
	}

	c.JSON(http.StatusUnauthorized, res)
}

// NotFindResponse 404错误 Response
func NotFindResponse(c *gin.Context) {
	res := &ErrResponse{
		Code:    notFindErrCode,
		Message: "Entity Not find",
		Data:    nil,
	}

	c.JSON(http.StatusNotFound, res)
}

// ParamErrResponse 参数检验错误 Response
func ParamErrResponse(c *gin.Context, err error) {
	res := &ErrResponse{
		Code:    paramErrCode,
		Message: err.Error(),
	}

	c.JSON(http.StatusUnprocessableEntity, res)
}
