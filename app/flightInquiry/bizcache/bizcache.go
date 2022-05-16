package bizcache

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/hash"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
	"time"
	"uranus/common/stringandbytes"
	"uranus/commonModel"
)

const BizFLICachePrefix = `biz:fli:cache:%s`
const BizSpaceCachePrefix = `biz:space:cache:%s`
const BizTransferCachePrefix = `biz:transfer:cache:%s`

type FLI commonModel.FlightInfos

// AddID 提供id存储
func AddID(r redis.Redis, id int64, zset, prefix string) error {
	bizFLICacheKey := fmt.Sprintf(prefix, zset)
	_, err := r.Zadd(bizFLICacheKey, time.Now().UnixNano()/1e6, strconv.Itoa(int(id)))
	if err != nil {
		return err
	}

	// 过期时间为 12 小时
	err = r.Expire(bizFLICacheKey, 43200)
	return err
}

func AddTransfer(r redis.Redis, transfer commonModel.Transfer, zset, prefix string) error {
	bizTransferCacheKey := fmt.Sprintf(prefix, zset)
	content, err := jsonx.Marshal(transfer)
	if err != nil {
		return err
	}

	_, err = r.Zadd(bizTransferCacheKey, time.Now().UnixNano()/1e6, string(content))
	if err != nil {
		return err
	}

	// 过期时间为 6 小时
	err = r.Expire(bizTransferCacheKey, 21600)
	return err
}

// DelID 提供id删除
func DelID(r redis.Redis, id int64, zset, prefix string) error {
	bizFLICacheKey := fmt.Sprintf(prefix, zset)
	_, err := r.Zrem(bizFLICacheKey, strconv.Itoa(int(id)))
	return err
}

// 航班内容压缩
func compress(f *FLI) (string, error) {
	data, err := f.Marshal()
	if err != nil {
		return "", err
	}

	return stringandbytes.Bytes2String(data), nil
}

// 航班内容解压
func unCompress(v string) (*FLI, error) {
	data := stringandbytes.String2Bytes(v)
	var ret *FLI
	ret, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// ListByRangeTime 提供根据时间段进行数据查询
func ListByRangeTime(r redis.Redis, zset, prefix string, start, end time.Time) (map[int64]struct{}, error) {
	idList := make(map[int64]struct{})
	bizFLICacheKey := fmt.Sprintf(prefix, zset)
	kvs, err := r.ZrangebyscoreWithScores(bizFLICacheKey, start.UnixNano()/1e6, end.UnixNano()/1e6)
	if err != nil {
		return nil, err
	}

	for _, kv := range kvs {
		id, _ := strconv.ParseInt(kv.Key, 10, 64)
		idList[id] = struct{}{}
	}

	return idList, nil
}

func ListAll(r redis.Redis, zset, prefix string) (idList map[int64]struct{}, err error) {
	idList = make(map[int64]struct{})
	bizFLICacheKey := fmt.Sprintf(prefix, zset)
	stringList, err := r.Zrange(bizFLICacheKey, 0, -1)
	if err != nil {
		return nil, err
	}

	for _, s := range stringList {
		id, _ := strconv.ParseInt(s, 10, 64)
		idList[id] = struct{}{}
	}

	return idList, nil
}

func ListAllTransfers(r redis.Redis, zset, prefix string) (transfers []*commonModel.Transfer, err error) {
	bizTransferCacheKey := fmt.Sprintf(prefix, zset)
	stringContent, err := r.Zrange(bizTransferCacheKey, 0, -1)
	if err != nil {
		return nil, err
	}

	transfers = make([]*commonModel.Transfer, 0)
	for _, s := range stringContent {
		transfer := &commonModel.Transfer{}
		err := jsonx.UnmarshalFromString(s, transfer)
		if err != nil {
			return nil, err
		}

		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

func ListByZrangeStartStop(r redis.Redis, zset, prefix string, start, stop int64) (idList map[int64]struct{}, err error) {
	idList = make(map[int64]struct{})
	bizFLICacheKey := fmt.Sprintf(prefix, zset)
	stringList, err := r.Zrange(bizFLICacheKey, start, stop)
	if err != nil {
		return nil, err
	}

	for _, s := range stringList {
		id, _ := strconv.ParseInt(s, 10, 64)
		idList[id] = struct{}{}
	}

	return idList, nil

}

func (f *FLI) Marshal() ([]byte, error) {
	return jsonx.Marshal(f)
}

func Unmarshal(data []byte) (*FLI, error) {
	var fli FLI
	err := jsonx.Unmarshal(data, &fli)
	if err != nil {
		return nil, err
	}

	return &fli, nil
}

// String2Score 字符串转换成 score
func String2Score(s string) int64 {
	data := stringandbytes.String2Bytes(s)
	return int64(hash.Hash(data) / 2)
}
