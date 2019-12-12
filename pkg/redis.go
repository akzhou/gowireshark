/**
 * @Author: Administrator
 * @Description:
 * @File:  redis
 * @Version: 1.0.0
 * @Date: 2019/12/10 19:40
 */

package pkg

//var pool *redis.Pool
//
////TODO:初始化redis pool
//func init() {
//	pool = &redis.Pool{
//		MaxIdle:   wireSharkCfg.Redis[0].Idle,
//		MaxActive: wireSharkCfg.Redis[0].Active,
//		Dial: func() (redis.Conn, error) {
//			c, err := redis.Dial("tcp", wireSharkCfg.Redis[0].Addr, redis.DialPassword(wireSharkCfg.Redis[0].Password))
//			if err != nil {
//				return nil, err
//			}
//			return c, nil
//		},
//	}
//}
//
////TODO:增量更新
//func IncrBy(key string, step int) {
//	c := pool.Get()
//	defer c.Close()
//	if _, err := c.Do("INCRBY", key, step); err != nil {
//		log.Error(err)
//	}
//	if _, err := c.Do("EXPIRE", key, 5*60); err != nil {
//		log.Error(err)
//	}
//}
