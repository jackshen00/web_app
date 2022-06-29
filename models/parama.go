package models

// 定义请求的参数结构体

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// ParamsSignUp 注册请求参数
type ParamsSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required, eqfield=Password"`
}

// ParamsLogin 登录请求参数
type ParamsLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 投票
type ParamsVoteData struct {
	// UserUD 从请求中获取当前用户
	PostID    string `json:"post_id" binding:"required"`              // 帖子ID
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成(1)反对(-1)取消(0)
}

// ParamsPostList 获取帖子列表query string参数
type ParamsPostList struct {
	Page  int64  `json:"page" form:"page"`
	Size  int64  `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
}

// ParamsCommunityPostList 获取社区帖子列表query string参数
type ParamsCommunityPostList struct {
	ParamsPostList
	CommunityID int64 `json:"community_id" form:"community_id"`
}
