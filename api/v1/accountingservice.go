package v1

type GetUserNameReq struct {
	Id string `json:"id" binding:"required"`
}

type GetUserNameRes struct {
	UserName string `json:"user_name"`
}

type UserInfoRes struct {
	UserName string `json:"user_name"`
}
