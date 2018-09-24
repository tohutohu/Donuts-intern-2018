package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

type C struct {
	client *redis.Client
}

func New() *C {
	c := &C{}
	c.client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return c
}

func (c *C) Get(key string) *redis.StringCmd {
	return c.client.Get(key)
}

func (c *C) Set(key, value string, duration time.Duration) *redis.StatusCmd {
	return c.client.Set(key, value, duration)
}

func (c *C) HGet(key, field string) *redis.StringCmd {
	return c.client.HGet(key, field)
}

func (c *C) HSet(key string, field string, value interface{}) *redis.BoolCmd {
	return c.client.HSet(key, field, value)
}

func (c *C) LPush(key, value string) *redis.IntCmd {
	return c.client.LPush(key, value)
}

func (c *C) RPush(key, value string) *redis.IntCmd {
	return c.client.RPush(key, value)
}

func (c *C) LRange(key string, start, stop int64) *redis.StringSliceCmd {
	return c.client.LRange(key, start, stop)
}

func (c *C) LLen(key string) *redis.IntCmd {
	return c.client.LLen(key)
}

func (c *C) ZAdd(key string, value redis.Z) *redis.IntCmd {
	return c.client.ZAdd(key, value)
}

func (c *C) ZRange(key string, start, stop int64) *redis.StringSliceCmd {
	return c.client.ZRange(key, start, stop)
}

func (c *C) ZRemRangeByRank(key string, start, stop int64) *redis.IntCmd {
	return c.client.ZRemRangeByRank(key, start, stop)
}

func (c *C) ZRevRange(key string, start, stop int64) *redis.StringSliceCmd {
	return c.client.ZRevRange(key, start, stop)
}

func (c *C) ZIncrBy(key string, score float64, value string) *redis.FloatCmd {
	return c.client.ZIncrBy(key, score, value)
}

func main() {
	var method string
	var key string
	var field string
	var value string
	var start int64
	var stop int64
	var score float64

	flag.StringVar(&method, "m", "", "")
	flag.StringVar(&key, "k", "", "")
	flag.StringVar(&field, "f", "", "")
	flag.StringVar(&value, "v", "", "")
	flag.Int64Var(&start, "start", 0, "")
	flag.Int64Var(&stop, "stop", 0, "")
	flag.Float64Var(&score, "score", 0.0, "")
	flag.Parse()
	if flag.NFlag() < 2 {
		flag.PrintDefaults()
		return
	}
	c := New()

	switch method {
	case "Get":
		res := c.Get(key)
		errCh(res)
		fmt.Println(res.Val())

	case "Set":
		res := c.Set(key, value, 1000*time.Hour)
		errCh(res)
		fmt.Println(res.Val())

	case "HGet":
		res := c.HGet(key, field)
		errCh(res)
		fmt.Println(res.Val())

	case "HSet":
		res := c.HSet(key, field, value)
		errCh(res)
		fmt.Println(res.Val())

	case "LPush":
		res := c.LPush(key, value)
		errCh(res)
		fmt.Println(res.Val())

	case "RPush":
		res := c.RPush(key, value)
		errCh(res)
		fmt.Println(res.Val())

	case "LRange":
		res := c.LRange(key, start, stop)
		errCh(res)
		fmt.Println(res.Val())

	case "LLen":
		res := c.LLen(key)
		errCh(res)
		fmt.Println(res.Val())

	case "ZAdd":
		res := c.ZAdd(key, redis.Z{Score: score, Member: value})
		errCh(res)
		fmt.Println(res.Val())

	case "ZRemRangeByRank":
		res := c.ZRemRangeByRank(key, start, stop)
		errCh(res)
		fmt.Println(res.Val())
	case "ZRevRange":
		res := c.ZRevRange(key, start, stop)
		errCh(res)
		fmt.Println(res.Val())
	case "ZRange":
		res := c.ZRange(key, start, stop)
		errCh(res)
		fmt.Println(res.Val())

	case "ZIncrBy":
		res := c.ZAdd(key, redis.Z{Score: score, Member: value})
		errCh(res)
		fmt.Println(res.Val())
	}
}

func errCh(cmd redis.Cmder) {
	if cmd.Err() != nil {
		fmt.Println(cmd.Err())
		os.Exit(1)
	}
}
