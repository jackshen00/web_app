package controllers

import (
	"errors"
	"fmt"
	"web_app/dao/mysql"
	"web_app/logic"
	"web_app/models"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(c *gin.Context) {
	// 1.获取参数和参数校验
	p := new(models.ParamsSignUp)

	// 1. 获取参数和参数校验
	// ShouldBindJSON 只能进行简单的数据校验，格式对不对
	err := c.ShouldBindJSON(p)
	if err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			//c.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			ResponseError(c, CodeInvalidParam)
			return
		}
		//c.JSON(http.StatusOK, gin.H{
		//	"mag": removeTopStruct(errors.Translate(trans)),
		//})
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errors.Translate(trans)))
		return
	}

	// 2.业务处理
	err1 := logic.SignUp(p)
	if err1 != nil {
		zap.L().Error("logic.signup failed", zap.Error(err1))
		if errors.Is(err1, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
		//fmt.Println(err1)
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "注册失败",
		//	"err": fmt.Sprintf("%s", err1),
		//})

		return
	}
	// 3.返回响应
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "success",
	//})
	//return
	ResponseSuccess(c, nil)
}

func Loginhandler(c *gin.Context) {
	// 1.获取请求参数及参数校验
	p := new(models.ParamsLogin)
	err := c.ShouldBindJSON(p)
	if err != nil {
		zap.L().Error("login with invalid param", zap.Error(err))
		err, ok := err.(validator.ValidationErrors)
		if !ok {
			//c.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(err.Translate(trans)))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(err.Translate(trans)),
		//})
		return
	}
	//fmt.Printf("login p:%v \n", p)

	// 2.业务逻辑处理

	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("login failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "用户名或密码有误",
		//})
		return
	}

	//fmt.Printf("jwt token:%s;\n", token)
	// 3.返回响应
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "登录成功",
	//})
	//return
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserId),
		"user_name": user.Username,
		"token":     user.Token,
	})
}
