package redis

import (
	"web_app/models"

	"github.com/go-redis/redis"
)

// GetPostIDsInOrder
func GetPostIDsInOrder(p *models.ParamsPostList) ([]string, error) {
	// 从redis获取id
	// 根据 order 确定取数方式
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 确定查询的索引起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1

	// zrevrange 查询
	return client.ZRevRange(key, start, end).Result()
}

// GetPosyVoteDate 根据ids 查询帖子的票数据
func GetPosyVoteDate(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPF + id)
	//	// 查找key中分数是1的元素数量
	//	v := client.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}

	// 使用pip 减少rtt
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	for _, cmder := range cmders {
		val := cmder.(*redis.IntCmd).Val()
		data = append(data, val)
	}

	return
}

// GetCommunityPostIDsInOrder 按社区根据ids 查询帖子的票数据
//func GetCommunityPostIDsInOrder() ([]string, error) {
//
//	// 从redis获取id
//	// 根据 order 确定取数方式
//	key := getRedisKey(KeyPostTimeZSet)
//	if p.Order == models.OrderScore {
//		key = getRedisKey(KeyPostScoreZSet)
//	}
//	// 确定查询的索引起始点
//	start := (p.Page - 1) * p.Size
//	end := start + p.Size - 1
//
//	// zrevrange 查询
//	return client.ZRevRange(key, start, end).Result()
//
//}
