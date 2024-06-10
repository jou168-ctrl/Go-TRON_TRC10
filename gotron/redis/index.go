package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisUtil struct {
	ctx context.Context
	db  *redis.Client
}

var Rd RedisUtil

func init() {
	Rd = RedisUtil{}
	url := "127.0.0.1:6379"
	password := ""
	Rd.Connect(url, password)
}

func (redisUtil *RedisUtil) Connect(url string, password string) {
	redisUtil.ctx = context.Background()
	redisUtil.db = redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       0,
	})
}

func (redisUtil *RedisUtil) Set(key string, value any, expiration time.Duration) error {
	err := redisUtil.db.Set(redisUtil.ctx, key, value, expiration).Err()
	if err != nil {
		panic(err)
	}
	return err
}

// 通过key获取值
// err == redis.Nil, "key不存在"
// err != nil,  "错误"
// val == "", "值是空字符串"
func (redisUtil *RedisUtil) Get(key string) (string, error) {
	value, err := redisUtil.db.Get(redisUtil.ctx, key).Result()
	if err != nil && value != "" {
		panic(err)
	}
	return value, err
}

// 删除key
func (redisUtil *RedisUtil) Del(keys ...string) error {
	err := redisUtil.db.Del(redisUtil.ctx, keys...).Err()
	if err != nil {
		panic(err)
	}
	return err
}

// 指定key的值+1
func (redisUtil *RedisUtil) Incr(key string) error {
	err := redisUtil.db.Incr(redisUtil.ctx, key).Err()
	if err != nil {
		panic(err)
	}
	return err
}

// 指定key的值+1
func (redisUtil *RedisUtil) Decr(key string) error {
	err := redisUtil.db.Decr(redisUtil.ctx, key).Err()
	if err != nil {
		panic(err)
	}
	return err
}

// 批量设置，没有过期时间
func (redisUtil *RedisUtil) MSet(values ...interface{}) error {
	err := redisUtil.db.MSet(redisUtil.ctx, values...).Err()
	return err
}

// 批量设置取数据
// 示例：values, err := MGet(key1, key2)
// for i, _ := range values {
// fmt.Println(values[i])
// }
func (redisUtil *RedisUtil) MGet(keys ...string) ([]interface{}, error) {
	values, err := redisUtil.db.MGet(redisUtil.ctx, keys...).Result()
	return values, err
}

// 执行命令
// 返回结果
// s, err := cmd.Text()
// flag, err := cmd.Bool()
// num, err := cmd.Int()
// num, err := cmd.Int64()
// num, err := cmd.Uint64()
// num, err := cmd.Float32()
// num, err := cmd.Float64()
// ss, err := cmd.StringSlice()
// ns, err := cmd.Int64Slice()
// ns, err := cmd.Uint64Slice()
// fs, err := cmd.Float32Slice()
// fs, err := cmd.Float64Slice()
// bs, err := cmd.BoolSlice()
func (redisUtil *RedisUtil) Do(args ...interface{}) *redis.Cmd {
	cmd := redisUtil.db.Do(redisUtil.ctx, args)
	return cmd
}

// 清空缓存
func (redisUtil *RedisUtil) FlushDB() error {
	err := redisUtil.db.FlushDB(redisUtil.ctx).Err()
	return err
}

// 发布
// 示例Publish("mychannel1", "payload").Err()
func (redisUtil *RedisUtil) Publish(channel string, msg string) error {
	err := redisUtil.db.Publish(redisUtil.ctx, channel, msg).Err()
	return err
}

// 订阅
func (redisUtil *RedisUtil) Subscribe(channel string, subscribe func(msg *redis.Message, err error)) {
	pubsub := redisUtil.db.Subscribe(redisUtil.ctx, channel)
	// 使用完毕，记得关闭
	defer pubsub.Close()
	for {
		msg, err := pubsub.ReceiveMessage(redisUtil.ctx)
		subscribe(msg, err)
	}
}

// 列表的头部（左边）,尾部（右边）
// 列表左边插入
func (redisUtil *RedisUtil) LPust(channel string, values ...interface{}) error {
	return redisUtil.db.LPush(redisUtil.ctx, channel, values...).Err()
}

// 列表从左边开始取出start至stop位置的数据
func (redisUtil *RedisUtil) LRange(key string, start, stop int64) error {
	return redisUtil.db.LRange(redisUtil.ctx, key, start, stop).Err()
}

// 列表左边取出
func (redisUtil *RedisUtil) LPop(key string) *redis.StringCmd {
	return redisUtil.db.LPop(redisUtil.ctx, key)
}

// 列表右边插入
func (redisUtil *RedisUtil) RPust(channel string, values ...interface{}) error {
	return redisUtil.db.RPush(redisUtil.ctx, channel, values...).Err()
}

// 列表右边取出
func (redisUtil *RedisUtil) RPop(key string) error {
	return redisUtil.db.RPop(redisUtil.ctx, key).Err()
}

// 列表哈希插入
func (redisUtil *RedisUtil) HSet(key string, values ...interface{}) error {
	return redisUtil.db.HSet(redisUtil.ctx, key, values...).Err()
}

// 列表哈希取出
func (redisUtil *RedisUtil) HGet(key, field string) *redis.StringCmd {
	return redisUtil.db.HGet(redisUtil.ctx, key, field)
}

// 列表哈希批量插入
func (redisUtil *RedisUtil) HMSet(key string, values ...interface{}) error {
	return redisUtil.db.HMSet(redisUtil.ctx, key, values...).Err()
}

// 列表哈希批量取出
func (redisUtil *RedisUtil) HMGet(key string, fields ...string) []interface{} {
	return redisUtil.db.HMGet(redisUtil.ctx, key, fields...).Val()
}

// 列表无序集合插入
func (redisUtil *RedisUtil) SAdd(key string, members ...interface{}) error {
	return redisUtil.db.SAdd(redisUtil.ctx, key, members...).Err()
}

// 列表无序集合，返回所有元素
func (redisUtil *RedisUtil) SMembers(key string) []string {
	return redisUtil.db.SMembers(redisUtil.ctx, key).Val()
}

// 列表无序集合，检查元素是否存在
func (redisUtil *RedisUtil) SIsMember(key string, member interface{}) bool {
	b, err := redisUtil.db.SIsMember(redisUtil.ctx, key, member).Result()
	if err != nil {
		panic(err)
	}
	return b
}
