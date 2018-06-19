package redisdao

import (
	"fmt"
	"math/rand"
	"time"

	"go_libs/service/logger"
	"go_libs/utils/confutil"

	"github.com/garyburd/redigo/redis"
)

type RedisDao struct {
	pool *redis.Pool
	Locker
}

func ExecuteCommand(cluster string, commandName string, args ...interface{}) (reply interface{}, err error) {
	instance := GetInstance(cluster)
	if instance == nil {
		err = fmt.Errorf("No redis cluster %s!", cluster)
		logger.E("redis", err)
		return
	}
	conn := instance.Get()
	defer conn.Close()
	reply, err = conn.Do(commandName, args...)
	if err != nil {
		logger.E("redis", err)
	}
	return
}

func Del(cluster string, keys ...interface{}) (err error) {
	_, err = ExecuteCommand(cluster, "del", keys...)
	return
}

func NewRedisDao(server []string) *RedisDao {
	this := new(RedisDao)
	rand.Seed(time.Now().UnixNano())
	this.pool = &redis.Pool{
		MaxIdle:     100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			x := rand.Intn(len(server))
			//ToDo:针对集群增加超时
			c, err := redis.Dial("tcp", server[x])
			if err != nil {
				logger.E("RedisConnect", "connect to redis:%s err=[%+v]", server[x], err)
				return nil, err
			}
			logger.I("RedisConnect", "connect to redis:%s\n", server[x])
			return c, err
		},
	}
	return this
}

func (this *RedisDao) Close() {
	this.pool.Close()
}

func (this *RedisDao) Get() redis.Conn {
	return this.pool.Get()
}

// Factory Method
var (
	redis_instances     map[string]*RedisDao
	default_lock_expire = 6
)

func init() {
	redis_instances = make(map[string]*RedisDao, 0)
	for k, v := range confutil.GetConfArrayMap("Redis") {
		redis_instances[k] = NewRedisDao(v)
	}
}

func GetInstance(key string) *RedisDao {
	if instance, ok := redis_instances[key]; ok {
		return instance
	} else {
		return nil
	}
}

type Locker struct {
	Key   string
	Value string
	redis.Conn
	Timeout int
}

func (this *Locker) Lock() error {
	_, err := redis.String(this.Do("SET", this.Key, this.Value, "EX", this.Timeout, "NX"))
	if err != nil {
		return err
	}
	return nil
}

func (this *Locker) Unlock() (err error) {
	_, err = this.Do("del", this.Key)
	return
}
func (this *RedisDao) Lock(conn redis.Conn, key string, val string, timeout ...int) (*Locker, error) {
	lock := &Locker{Key: key, Value: val, Conn: conn}
	if len(timeout) > 0 {
		lock.Timeout = timeout[0]
	} else {
		lock.Timeout = default_lock_expire
	}
	err := lock.Lock()
	if err != nil {
		return nil, err
	}
	return lock, nil
}
