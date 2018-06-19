package redisdao

import (
	"fmt"
	"go_libs/service/logger"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/cast"
)

/*
RawKey in redis will be Prefix/Root/Id
Default Prefix = var PREFIX
*/

var PREFIX string = "CACHE"

type redisCache struct {
	root   string
	prefix string
	ins    *RedisDao
}

func NewRedisCache(root string) (this *redisCache) {
	this = new(redisCache)
	this.root = root
	this.prefix = PREFIX
	this.ins = GetInstance("cache")
	return
}

func (this *redisCache) Prefix(prefix string) {
	this.prefix = prefix
}

func (this *redisCache) set(key string, value interface{}, option string, expire int64) (ok bool) {
	conn := this.ins.Get()
	defer conn.Close()
	var reply interface{}
	var err error
	params := []interface{}{key, value}
	if len(option) > 0 {
		params = append(params, option)
	}
	if expire > 0 {
		params = append(params, "EX")
		params = append(params, expire)
	}
	reply, err = conn.Do("set", params...)
	if err != nil {
		logger.W("RedisCache", "set %s %v %s expire %d Failed %v", key, value, option, expire, err)
		return false
	}
	if reply == nil {
		return false
	}
	if r, e := redis.String(reply, err); r == "OK" && e == nil {
		return true
	} else {
		logger.W("RedisCache", "set %s %v %s expire %d;reply:%v %v", key, value, option, expire, reply, e)
	}
	return false
}

func (this *redisCache) hset(key string, field string, value interface{}) (ok bool) {
	conn := this.ins.Get()
	defer conn.Close()
	var reply interface{}
	var err error
	reply, err = conn.Do("hset", key, field, value)
	if err != nil {
		logger.W("RedisCache", "hset %s %s %v Failed %v", key, field, value, err)
		return false
	}
	if reply == nil {
		return false
	}
	return true
}

func (this *redisCache) hincrby(key string, field string, value int64) (reply interface{}, ok bool) {
	conn := this.ins.Get()
	defer conn.Close()
	var err error
	reply, err = conn.Do("hincrby", key, field, value)
	if err != nil {
		logger.W("RedisCache", "hincrby %s %s %v Failed %v", key, field, value, err)
		ok = false
	}
	if reply == nil {
		ok = false
	}
	ok = true
	return
}
func (this *redisCache) ttl(key string) (ret interface{}, ok bool) {
	conn := this.ins.Get()
	defer conn.Close()
	reply, err := conn.Do("ttl", key)
	if err != nil {
		logger.W("RedisCache", "ttl %v Failed %v", key, err)
		ok = false
	} else {
		ret = reply
		ok = true
	}
	return
}

func (this *redisCache) get(key string) (ret interface{}, ok bool) {
	conn := this.ins.Get()
	defer conn.Close()
	reply, err := conn.Do("get", key)
	if err != nil {
		logger.W("RedisCache", "get %v Failed %v", key, err)
		ok = false
	} else {
		ret = reply
		ok = true
	}
	return
}

func (this *redisCache) hget(key string, field string) (ret interface{}, ok bool) {
	conn := this.ins.Get()
	defer conn.Close()
	reply, err := conn.Do("hget", key, field)
	if err != nil {
		logger.W("RedisCache", "hget %s %s Failed %v", key, field, err)
		ok = false
	} else {
		ret = reply
		ok = true
	}
	return
}

func (this *redisCache) hgetall(key string) (ret interface{}, ok bool) {
	conn := this.ins.Get()
	defer conn.Close()
	reply, err := conn.Do("hgetall", key)
	if err != nil {
		logger.W("RedisCache", "hgetall %s Failed %v", key, err)
		ok = false
	} else {
		ret = reply
		ok = true
	}
	return
}

func (this *redisCache) hdel(key string, field string) (ok bool) {
	conn := this.ins.Get()
	defer conn.Close()
	reply, err := conn.Do("hdel", key, field)
	if err != nil {
		logger.W("RedisCache", "hdel %s %s Failed %v", key, field, err)
		return false
	}
	if reply == nil {
		logger.W("RedisCache", "hdel %s %s Failed %v", key, field, reply)
		return false
	}
	if r, _ := redis.Int(reply, err); r > 0 {
		return true
	}
	return false
}

func (this *redisCache) zincrby(key, member string, score float64) (ret float64, ok bool) {
	conn := this.ins.Get()
	defer conn.Close()
	reply, err := conn.Do("zincrby", key, score, member)
	if err != nil {
		logger.W("RedisCache", "zincrby %s %f %s Failed %v", key, score, member, err)
		ok = false
	} else {
		var er2 error
		ret, er2 = redis.Float64(reply, err)
		if er2 != nil {
			logger.W("RedisCache", er2)
		}
		ok = true
	}
	return
}

func (this *redisCache) zadd(key, member string, score float64) (ret interface{}, ok bool) {
	conn := this.ins.Get()
	defer conn.Close()
	reply, err := conn.Do("zadd", key, score, member)
	if err != nil {
		logger.W("RedisCache", "zadd %s %f %s Failed %v", key, score, member, err)
		ok = false
	} else {
		ret = reply
		ok = true
	}
	return
}

func (this *redisCache) zcard(key string) (ret interface{}, ok bool) {
	conn := this.ins.Get()
	defer conn.Close()
	reply, err := conn.Do("zcard", key)
	if err != nil {
		logger.W("RedisCache", "zcard %s Failed %v", key, err)
		ok = false
	} else {
		ret = reply
		ok = true
	}
	return
}

func (this *redisCache) zrem(key string, field string) (ok bool) {
	conn := this.ins.Get()
	defer conn.Close()
	reply, err := conn.Do("zrem", key, field)
	if err != nil {
		logger.W("RedisCache", "zrem %s %s Failed %v", key, field, err)
		ok = false
	} else {
		if reply == nil {
			return false
		}
		if r, _ := redis.Int(reply, err); r > 0 {
			return true
		}
		ok = false
	}
	return
}

func (this *redisCache) zscore(key, member string) (ret interface{}, ok bool) {
	conn := this.ins.Get()
	defer conn.Close()
	reply, err := conn.Do("zscore", key, member)
	if err != nil {
		logger.W("RedisCache", "zscore %s %s Failed %v", key, member, err)
		ok = false
	} else {
		ret = reply
		ok = true
	}
	return
}

func (this *redisCache) zrevrangebyscore(key string, offset, count int64, option string) (ret interface{}, ok bool) {
	if count > 1000 || count <= 0 {
		count = 1000
	}
	conn := this.ins.Get()
	defer conn.Close()
	var reply interface{}
	var err error
	params := []interface{}{key, "+inf", "-inf", "limit", offset, count}
	if len(option) > 0 {
		params = append(params, option)
	}
	reply, err = conn.Do("zrevrangebyscore", params...)
	if err != nil {
		logger.W("RedisCache", "zrevrangebyscore %s +inf -inf limit %d %d Failed %v", key, offset, count, err)
		ok = false
	} else {
		ret = reply
		ok = true
	}
	return
}

func (this *redisCache) exists(key string) (ret bool, ok bool) {
	conn := this.ins.Get()
	defer conn.Close()
	reply, err := conn.Do("exists", key)
	if err != nil {
		logger.W("RedisCache", "exists %v Failed %v", key, err)
		ok = false
	} else {
		if r, e := redis.Int(reply, err); e == nil {
			if r == 1 {
				ret = true
			} else {
				ret = false
			}
			ok = true
		} else {
			ok = false
		}
	}
	return
}

func (this *redisCache) del(key string) (ok bool) {
	conn := this.ins.Get()
	defer conn.Close()
	reply, err := conn.Do("del", key)
	if err != nil {
		logger.W("RedisCache", "del %v Failed %v", key, err)
		return false
	}
	if reply == nil {
		return false
	}
	if r, _ := redis.Int(reply, err); r > 0 {
		return true
	}
	return false
}

func (this *redisCache) incrby(key string, value int64) (ret int64, ok bool) {
	conn := this.ins.Get()
	defer conn.Close()
	var reply interface{}
	var err error
	if value == 1 {
		reply, err = conn.Do("incr", key)
	} else {
		reply, err = conn.Do("incrby", key, value)
	}
	if err != nil {
		logger.W("RedisCache", "incrby %v %v Failed %v", key, value, err)
		ok = false
	} else {
		ret = int64(cast.ToInt(reply))
		ok = true
	}
	return
}

func (this *redisCache) Id(id interface{}) (ret *redisCommandObject) {
	ret = new(redisCommandObject)
	ret.rpage = *this
	ret.key = fmt.Sprintf("%s/%s/%v", ret.rpage.prefix, ret.rpage.root, id)
	ret.value = nil
	return
}

type redisCommandObject struct {
	rpage  redisCache
	key    string
	value  interface{}
	expire int64
	option string
	offset int64
	count  int64
	field  string
}

func (this *redisCommandObject) reproduce() (ret *redisCommandObject) {
	ret = new(redisCommandObject)
	ret.rpage = this.rpage
	ret.key = this.key
	ret.expire = this.expire
	ret.value = this.value
	ret.option = this.option
	ret.field = this.field
	ret.offset = this.offset
	ret.count = this.count
	return
}

func (this *redisCommandObject) Limit(offset int64, count int64) (ret *redisCommandObject) {
	ret = this.reproduce()
	ret.count = count
	ret.offset = offset
	return
}

func (this *redisCommandObject) Option(option string) (ret *redisCommandObject) {
	ret = this.reproduce()
	ret.option = option
	return
}

func (this *redisCommandObject) Field(field string) (ret *redisCommandObject) {
	ret = this.reproduce()
	ret.field = field
	return
}

func (this *redisCommandObject) Value(value interface{}) (ret *redisCommandObject) {
	ret = this.reproduce()
	ret.value = value
	return
}

func (this *redisCommandObject) Expire(ex int64) (ret *redisCommandObject) {
	ret = this.reproduce()
	ret.expire = ex
	return
}

func (this *redisCommandObject) Exists() (ret bool, ok bool) {
	ret, ok = this.rpage.exists(this.key)
	return
}

func (this *redisCommandObject) TTL() (ret *redisResult) {
	reply, ok := this.rpage.ttl(this.key)
	ret = new(redisResult)
	ret.Reply = reply
	ret.Ok = ok
	return
}

func (this *redisCommandObject) Get() (ret *redisResult) {
	reply, ok := this.rpage.get(this.key)
	ret = new(redisResult)
	ret.Reply = reply
	ret.Ok = ok
	return
}

func (this *redisCommandObject) HGet() (ret *redisResult) {
	ret = new(redisResult)
	if len(this.field) == 0 {
		ret.Ok = false
		return
	}
	reply, ok := this.rpage.hget(this.key, this.field)
	ret.Reply = reply
	ret.Ok = ok
	return
}

func (this *redisCommandObject) HGetAll() (ret *redisResult) {
	reply, ok := this.rpage.hgetall(this.key)
	ret = new(redisResult)
	ret.Reply = reply
	ret.Ok = ok
	return
}

func (this *redisCommandObject) HSet() (ok bool) {
	if this.value != nil && len(this.field) > 0 {
		ok = this.rpage.hset(this.key, this.field, this.value)
	} else {
		logger.W("RedisCache", "missing required parameters")
		ok = false
	}
	return
}

func (this *redisCommandObject) HDel() (ok bool) {
	if len(this.field) > 0 {
		ok = this.rpage.hdel(this.key, this.field)
	} else {
		ok = false
	}
	return
}

func (this *redisCommandObject) HIncrBy(value int64) (ret *redisResult) {
	ret = new(redisResult)
	if len(this.field) > 0 {
		reply, ok := this.rpage.hincrby(this.key, this.field, value)
		ret.Reply = reply
		ret.Ok = ok
	} else {
		ret.Ok = false
	}
	return
}

func (this *redisCommandObject) ZAdd() (ok bool) {
	if len(this.field) > 0 {
		switch value := this.value.(type) {
		case float64:
			_, ok = this.rpage.zadd(this.key, this.field, value)
		case int:
			_, ok = this.rpage.zadd(this.key, this.field, float64(value))
		case int64:
			_, ok = this.rpage.zadd(this.key, this.field, float64(value))
		default:
			ok = false
		}
	} else {
		ok = false
	}
	return
}

func (this *redisCommandObject) ZCard() (ret *redisResult) {
	reply, ok := this.rpage.zcard(this.key)
	ret = new(redisResult)
	ret.Reply = reply
	ret.Ok = ok
	return
}

func (this *redisCommandObject) ZScore() (ret *redisResult) {
	reply, ok := this.rpage.zscore(this.key, this.field)
	ret = new(redisResult)
	ret.Reply = reply
	ret.Ok = ok
	return
}

func (this *redisCommandObject) ZRem() (ok bool) {
	if len(this.field) > 0 {
		ok = this.rpage.zrem(this.key, this.field)
	} else {
		ok = false
	}
	return
}

func (this *redisCommandObject) ZIncrBy(score float64) (ret float64, ok bool) {
	if len(this.field) > 0 {
		ret, ok = this.rpage.zincrby(this.key, this.field, score)
	} else {
		ok = false
	}
	return
}

func (this *redisCommandObject) ZRevRangeByScore() (ret *redisResult) {
	reply, ok := this.rpage.zrevrangebyscore(this.key, this.offset, this.count, this.option)
	ret = new(redisResult)
	ret.Reply = reply
	ret.Ok = ok
	return
}

func (this *redisCommandObject) Set() (ok bool) {
	if this.value != nil {
		ok = this.rpage.set(this.key, this.value, this.option, this.expire)
	} else {
		ok = false
	}
	return
}

func (this *redisCommandObject) Del() (ok bool) {
	ok = this.rpage.del(this.key)
	return
}

func (this *redisCommandObject) Incr() (ret int64, ok bool) {
	ret, ok = this.rpage.incrby(this.key, 1)
	return
}

func (this *redisCommandObject) Incrby(value int64) (ret int64, ok bool) {
	ret, ok = this.rpage.incrby(this.key, value)
	return
}

type redisResult struct {
	Reply interface{}
	Ok    bool
}

func (this *redisResult) Interface() (ret interface{}, ok bool) {
	ret, ok = this.Reply, true
	return
}
func (this *redisResult) Int() (ret int, ok bool) {
	if v, e := redis.Int(this.Reply, nil); e == nil && this.Ok {
		ret, ok = v, true
		return
	} else {
		ok = false
	}
	return
}
func (this *redisResult) Int64() (ret int64, ok bool) {
	if v, e := redis.Int64(this.Reply, nil); e == nil && this.Ok {
		ret, ok = v, true
		return
	} else {
		ok = false
	}
	return
}
func (this *redisResult) Bool() (ret bool, ok bool) {
	if v, e := redis.Bool(this.Reply, nil); e == nil && this.Ok {
		ret, ok = v, true
		return
	} else {
		ok = false
	}
	return
}
func (this *redisResult) Float64() (ret float64, ok bool) {
	if v, e := redis.Float64(this.Reply, nil); e == nil && this.Ok {
		ret, ok = v, true
		return
	} else {
		ok = false
	}
	return
}
func (this *redisResult) String() (ret string, ok bool) {
	if v, e := redis.String(this.Reply, nil); e == nil && this.Ok {
		ret, ok = v, true
		return
	} else {
		ok = false
	}
	return
}
func (this *redisResult) Ints() (ret []int, ok bool) {
	if v, e := redis.Ints(this.Reply, nil); e == nil && this.Ok {
		ret, ok = v, true
		return
	} else {
		ok = false
	}
	return
}
func (this *redisResult) Bytes() (ret []byte, ok bool) {
	if v, e := redis.Bytes(this.Reply, nil); e == nil && this.Ok {
		ret, ok = v, true
		return
	} else {
		ok = false
	}
	return
}
func (this *redisResult) Strings() (ret []string, ok bool) {
	if v, e := redis.Strings(this.Reply, nil); e == nil && this.Ok {
		ret, ok = v, true
		return
	} else {
		ok = false
	}
	return
}
func (this *redisResult) IntMap() (ret map[string]int, ok bool) {
	if v, e := redis.IntMap(this.Reply, nil); e == nil && this.Ok {
		ret, ok = v, true
		return
	} else {
		ok = false
	}
	return
}
func (this *redisResult) StringMap() (ret map[string]string, ok bool) {
	if v, e := redis.StringMap(this.Reply, nil); e == nil && this.Ok {
		ret, ok = v, true
		return
	} else {
		ok = false
	}
	return
}
