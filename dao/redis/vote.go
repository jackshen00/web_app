package redis

import (
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432
)

var (
	ErrVoteTimeExpire = errors.New("超出投票时间")
	ErrVoteRepested   = errors.New("不允许重复投票")
)

// CreatePost
func CreatePost(postID int64) error {
	// 事务操作
	pipeline := client.TxPipeline()

	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	_, err := pipeline.Exec()
	return err
}

// VoteForPost
func VoteForPost(userID, postID string, value float64) error {
	// 1、判断投票限制
	// 去redis取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	// 2 3 需要放到一个事务中操作 pipeline

	// 2、更新帖子分数
	// 先查询之前的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	if value == ov {
		return ErrVoteRepested
	}

	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}

	diff := math.Abs(ov - value)

	pipeline := client.TxPipeline()

	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)
	//if ErrVoteTimeExpire != nil {
	//	return err
	//}
	// 3、记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), postID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value, // 赞成票还是反对票
			Member: nil,
		})
	}
	_, err := pipeline.Exec()
	return err
}
