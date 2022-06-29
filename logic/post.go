package logic

import (
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 1、生成post id
	p.ID = int64(snowflake.GenID())
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID)
	return err
}

// GetPostById 根据id返回帖子内容
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合我们接口想要的数据

	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed ", zap.Int64("pid", pid), zap.Error(err))
		return
	}
	// 根据作者id查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById failed ", zap.Int64("AuthorID", post.AuthorID), zap.Error(err))
		return
	}
	// 根据社区id查询社区信息
	community, err := mysql.GetCommunityDetailByID(int64(post.CommunityID))
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID failed ", zap.Int32("CommunityID", post.CommunityID), zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed ", zap.Error(err))
		return
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))

	for _, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed ", zap.Int64("AuthorID", post.AuthorID), zap.Error(err))
			continue
		}
		// 根据社区id查询社区信息
		community, err := mysql.GetCommunityDetailByID(int64(post.CommunityID))
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed ", zap.Int32("CommunityID", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostList2 帖子列表 升级版
func GetPostList2(p *models.ParamsPostList) (data []*models.ApiPostDetail, err error) {
	// 2、去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	zap.L().Debug("GetPostIDsInOrder", zap.Any("ids", ids))
	if err != nil {
		return
	}
	// 3、根据id去数据库查询帖子详细信息
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder return 0")
		return
	}
	posts, err := mysql.GetPostListByIDs(ids)
	zap.L().Debug("GetPostListByIDs", zap.Any("posts", posts))
	if err != nil {
		return
	}

	// 提前查询好每篇帖子的投票数据
	voteDdata, err := redis.GetPosyVoteDate(ids)
	if err != nil {
		return
	}

	// 将帖子的作者以及分区信息查询出来填充到帖子里
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed ", zap.Int64("AuthorID", post.AuthorID), zap.Error(err))
			continue
		}
		// 根据社区id查询社区信息
		community, err := mysql.GetCommunityDetailByID(int64(post.CommunityID))
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed ", zap.Int32("CommunityID", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteDdata[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetCommunityPostList 获取社区帖子列表
//func GetCommunityPostList(p *models.ParamsCommunityPostList) (data []*models.ApiPostDetail, err error) {
//	// 2、去redis查询id列表
//	ids, err := redis.GetCommunityPostIDsInOrder(p)
//	zap.L().Debug("GetPostIDsInOrder", zap.Any("ids", ids))
//	if err != nil {
//		return
//	}
//	// 3、根据id去数据库查询帖子详细信息
//	if len(ids) == 0 {
//		zap.L().Warn("redis.GetPostIDsInOrder return 0")
//		return
//	}
//	posts, err := mysql.GetPostListByIDs(ids)
//	zap.L().Debug("GetPostListByIDs", zap.Any("posts", posts))
//	if err != nil {
//		return
//	}
//
//	// 提前查询好每篇帖子的投票数据
//	voteDdata, err := redis.GetPosyVoteDate(ids)
//	if err != nil {
//		return
//	}
//
//	// 将帖子的作者以及分区信息查询出来填充到帖子里
//	for idx, post := range posts {
//		// 根据作者id查询作者信息
//		user, err := mysql.GetUserById(post.AuthorID)
//		if err != nil {
//			zap.L().Error("mysql.GetUserById failed ", zap.Int64("AuthorID", post.AuthorID), zap.Error(err))
//			continue
//		}
//		// 根据社区id查询社区信息
//		community, err := mysql.GetCommunityDetailByID(int64(post.CommunityID))
//		if err != nil {
//			zap.L().Error("mysql.GetCommunityDetailByID failed ", zap.Int32("CommunityID", post.CommunityID), zap.Error(err))
//			continue
//		}
//		postDetail := &models.ApiPostDetail{
//			AuthorName:      user.Username,
//			VoteNum:         voteDdata[idx],
//			Post:            post,
//			CommunityDetail: community,
//		}
//		data = append(data, postDetail)
//	}
//	return
//}
