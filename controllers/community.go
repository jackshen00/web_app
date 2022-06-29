package controllers

import (
	"strconv"
	"web_app/logic"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ---- 社区相关 ----

func CommunityHandler(c *gin.Context) {
	// 查询到所有社区 (community_id, community_name) 以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	detail, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail(id) failed", zap.Error(err))
		return
	}
	ResponseSuccess(c, detail)
}
