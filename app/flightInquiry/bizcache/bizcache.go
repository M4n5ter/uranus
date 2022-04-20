package bizcache

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"
	"uranus/commonModel"
)

type FLI commonModel.FlightInfos

const bizContentCacheKey = `biz#fli#cache`

// AddFLI 提供航班内容存储
func AddFLI(r redis.Redis, c *FLI) error {
	v := compress(c)
	_, err := r.Zadd(bizContentCacheKey, c.CreateTime.UnixNano()/1e6, v)
	return err
}

// DelFLI 提供内容删除
func DelFLI(r redis.Redis, c *FLI) error {
	v := compress(c)
	_, err := r.Zrem(bizContentCacheKey, v)

	return err
}

// 内容压缩
func compress(c *FLI) string {
	// todo
	var ret string
	return ret
}

// 内容解压
func unCompress(v string) *FLI {
	// todo
	var ret FLI
	return &ret
}

// ListByRangeTime 提供根据时间段进行数据查询
func ListByRangeTime(r redis.Redis, start, end time.Time) ([]*FLI, error) {
	kvs, err := r.ZrangebyscoreWithScores(bizContentCacheKey, start.UnixNano()/1e6, end.UnixNano()/1e6)
	if err != nil {
		return nil, err
	}

	var list []*FLI
	for _, kv := range kvs {
		data := unCompress(kv.Key)
		list = append(list, data)
	}

	return list, nil
}
