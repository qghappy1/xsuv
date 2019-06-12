package redis

import (
	"fmt"
	"strconv"
	"strings"
)

type RedisDB struct {
	*Client
}

// connect:14.17.77.164:6381|0|myddz003@yz.com
func NewRedisDB(redisConfig string) (*RedisDB, error) {
	strAry := strings.Split(redisConfig, "|")
	if len(strAry) != 3 {
		return nil, fmt.Errorf("redis config error", redisConfig)
	}
	dbIndex, err := strconv.Atoi(strAry[1])
	if err != nil {
		return nil, err
	}
	db, err := newRedisDB(strAry[0], dbIndex, strAry[2])
	if err != nil {
		return nil, err
	}
	return db, nil
}

func newRedisDB(addr string, db int, pwd string) (*RedisDB, error) {
	inst := new(RedisDB)
	inst.Client = &Client{Addr: addr, Db: db, Password: pwd, MaxPoolSize: 50}
	if err := inst.Client.Auth(pwd); err != nil {
		return inst, err
	}
	return inst, nil
}
