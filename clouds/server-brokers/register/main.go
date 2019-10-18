package main

import (
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"time"
)

var RedisPool *redis.Pool

func init() {
	RedisPool = GetRedis("redis://localhost:6379")

}

// 获取redis池
func GetRedis(url string) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 200,
		//MaxActive:   0,
		IdleTimeout: 10 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(url)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
func main() {
	fmt.Println(">> register server starts")
	r := gin.Default()

	// for user-pool to storage online and offline info
	r.POST("/register/online-info/", onlineInfo)
	r.POST("/register/offline-info/", offlineInfo)

	// get pool info by pool_id
	r.POST("/register/add-pool-info/", addPoolInfo)
	r.POST("/register/get-pool-info/", getPoolInfo)
	// for center to get user-info, build pool bridge
	r.POST("/register/get-info/", getInfo)

	r.Run(":6666")
}

// storage online user-info
func onlineInfo(c *gin.Context) {
	type Param struct {
		UserId int64  `json:"user_id" binding:"required"` // 用户id，全局唯一
		PoolId string `json:"pool_id" binding:"required"` // 用户池id
	}
	var param Param
	if e := c.Bind(&param); e != nil {
		c.JSON(400, gin.H{"message": errorx.Wrap(e)})
		return
	}
	conn := RedisPool.Get()
	defer conn.Close()
	buf, _ := json.Marshal(param)
	var key = fmt.Sprintf("user_flash_info:%d", param.UserId)
	conn.Do("SETEX", key, 60*60, buf)
}

// remove online info
func offlineInfo(c *gin.Context) {
	type Param struct {
		UserId int64 `json:"user_id" binding:"required"` // 用户id，全局唯一
	}
	var param Param
	if e := c.Bind(&param); e != nil {
		c.JSON(400, gin.H{"message": errorx.Wrap(e)})
		return
	}
	conn := RedisPool.Get()
	defer conn.Close()
	var key = fmt.Sprintf("user_flash_info:%d", param.UserId)
	conn.Do("DEL", key)
}

// get online info
func getInfo(c *gin.Context) {
	type Param struct {
		UserId int64 `json:"user_id"`
	}
	var param Param
	if e := c.Bind(&param); e != nil {
		c.JSON(400, gin.H{"message": errorx.Wrap(e)})
		return
	}

	// tips scopes
	var tip1 = fmt.Sprintf("user_id flash info '%d' not found", param.UserId)

	conn := RedisPool.Get()
	defer conn.Close()
	var key = fmt.Sprintf("user_flash_info:%d", param.UserId)
	buf, e := redis.Bytes(conn.Do("GET", key))
	if e == redis.ErrNil || len(buf) == 0 {
		c.JSON(400, gin.H{
			"message": "user not found",
			"tip":     tip1,
			"tip_id":  1,
		})
		return
	}
	if e != nil {
		c.JSON(500, gin.H{"message": errorx.Wrap(e)})
		return
	}
	c.Status(200)
	c.Writer.Write(buf)
}

// add pool info
func addPoolInfo(c *gin.Context) {
	type Param struct {
		PoolId string `json:"pool_id"`
		Host   string `json:"host"`
	}
	var param Param
	if e := c.Bind(&param); e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		c.JSON(400, gin.H{"message": errorx.Wrap(e).Error()})
		return
	}

	conn := RedisPool.Get()
	defer conn.Close()
	var key = fmt.Sprintf("pool_flash_info:%s", param.PoolId)

	buf, _ := json.Marshal(param)

	conn.Do("SETEX", key, 60*60, buf)
	c.JSON(200, gin.H{"message": "success"})
}

func getPoolInfo(c *gin.Context) {
	type Param struct {
		PoolId string `json:"pool_id" binding:"required"`
	}
	type Result struct {
		PoolId string `json:"pool_id"`
		Host   string `json:"host"`
	}

	var param Param
	if e := c.Bind(&param); e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		c.JSON(400, gin.H{"message": errorx.Wrap(e).Error()})
		return
	}

	var key = fmt.Sprintf("pool_flash_info:%s", param.PoolId)
	conn := RedisPool.Get()
	defer conn.Close()
	buf, e := redis.Bytes(conn.Do("GET", key))
	if e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
		return
	}

	c.Status(200)
	c.Writer.Write(buf)
}
