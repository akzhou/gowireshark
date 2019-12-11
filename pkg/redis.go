/**
 * @Author: Administrator
 * @Description:
 * @File:  redis
 * @Version: 1.0.0
 * @Date: 2019/12/10 19:40
 */

package pkg

import (
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"runtime"
	"time"
)

var (
	pool      *redis.Pool
	redisConf *RedisConfig
	confPath  string
)

type RedisConfig struct {
	Redis []struct {
		Name         string
		Addr         string
		Active       int
		Idle         int
		DialTimeout  time.Duration
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
		DBNum        string
		Password     string
	}
}

//TODO:初始化redis pool
func init() {
	var tomlPath string
	if runtime.GOOS == `windows` {
		tomlPath = "e:/xinxinserver/config/gowireshark.toml"
	} else {
		tomlPath = "/config/gowireshark.toml"
	}
	flag.StringVar(&confPath, "conf", tomlPath, "config path")

	viper.SetConfigName("gowireshark")
	viper.SetConfigType("toml")
	viper.AddConfigPath(confPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&redisConf); err != nil {
		panic(err)
	}

	if len(redisConf.Redis) == 0 {
		panic(fmt.Errorf("未配置Redis"))
	}

	pool = &redis.Pool{
		MaxIdle:   redisConf.Redis[0].Idle,
		MaxActive: redisConf.Redis[0].Active,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisConf.Redis[0].Addr, redis.DialPassword(redisConf.Redis[0].Password))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}

//TODO:增量更新
func IncrBy(key string, step int) {
	c := pool.Get()
	defer c.Close()
	if _, err := c.Do("INCRBY", key, step); err != nil {
		log.Error(err)
	}
	if _, err := c.Do("EXPIRE", key, 5*60*60); err != nil {
		log.Error(err)
	}
}
