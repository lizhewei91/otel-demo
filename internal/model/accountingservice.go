package model

/*
 * model 目录下为业务输入/输出数据结构定义。
 * 定义数据结构时以业务模块名开头，输入和输出分别以 Input 和 Output 结尾。
 */

type GetUserNameInput struct {
	Id string `json:"id" binding:"required"`
}

type GetUserNameOutput struct {
	UserName string `json:"user_name"`
}
