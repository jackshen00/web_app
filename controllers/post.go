package controllers

import (
	"strconv"
	"web_app/logic"
	"web_app/models"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func CreatePostHandler(c *gin.Context) {
	// 1、获取参数及校验
	p := new(models.Post)
	err := c.ShouldBindJSON(p)
	if err != nil {
		zap.L().Debug("1111", zap.Any("err", err))
		zap.L().Error("ShouldBindJSON err")
		ResponseError(c, CodeInvalidParam)
		return
	}

	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2、创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3、返回响应
	ResponseSuccess(c, nil)
}

// CreatePostDetailHandler 获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	// 1、获取参数（从url中获取帖子id）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("post detail invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2、根据id取出帖子数据-查数据库
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3、返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表的处理函数
func GetPostListHandler(c *gin.Context) {

	// 获取分页参数
	page, size := getPageInfo(c)

	// 1、获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 2、返回响应
	ResponseSuccess(c, data)

}

// GetPostListHandler2 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口(api分组展示使用的)
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer JWT"
// @Param object query models.ParamsPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	// GET 请求参数 /api/v1/posts2?page=1
	// 初始化结构体时指定初始参数
	p := &models.ParamsPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 1、获取数据
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 2、返回响应
	ResponseSuccess(c, data)

}

// GetCommunityPostListHandler 根据社区查询帖子列表
//func GetCommunityPostListHandler(c *gin.Context) {
//	p := &models.ParamsCommunityPostList{
//		Page:  1,
//		Size:  10,
//		Order: models.OrderTime,
//	}
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("GetCommunityPostListHandler with invalid params", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//	// 1、获取数据
//	data, err := logic.GetCommunityPostList(p)
//	if err != nil {
//		zap.L().Error("logic.GetPostList failed", zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//		return
//	}
//	// 2、返回响应
//	ResponseSuccess(c, data)
//}
