package controllers

import (
	"web_app/logic"
	"web_app/models"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

func PostVoteController(c *gin.Context) {
	// 参数校验
	p := new(models.ParamsVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errors.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}
	// 获取用户id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 投票业务逻辑
	logic.VoteForPost(userID, p)
	ResponseSuccess(c, nil)
}
