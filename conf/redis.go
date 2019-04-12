package conf

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/rs/zerolog"
)

func OpenRedis(addr string, size int, logger zerolog.Logger) (*RedisClient, error) {
	p, err := pool.New("tcp", addr, size)
	if err != nil {
		return nil, err
	}

	c := new(RedisClient)
	c.Pool = p
	c.logger = logger
	return c, nil
}

type RedisClient struct {
	*pool.Pool
	logger zerolog.Logger
}

func (c *RedisClient) Cmd(cmd string, args ...interface{}) *redis.Resp {
	resp := c.Pool.Cmd(cmd, args...)
	c.logger.Debug().Str("cmd", cmd).Interface("args", args).Str("resp", resp.String()).Msg("Redis command.")
	return resp
}

type RedisResponse struct {
	*redis.Resp
}

func (r *RedisResponse) IsNil() bool {
	if r.Err != nil {
		panic(r.Err)
	}
	return r.IsType(redis.Nil)
}
