package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/fwhezfwhez/tcpx"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var PoolID string
var Host string
var Port string

func main() {
	fmt.Println(">> user-pool server starts")

	// specific pool id
	flag.StringVar(&PoolID, "pool_id", "pool_1", "-pool_id 1")
	flag.StringVar(&Host, "host", "localhost", "-host localhost")
	flag.StringVar(&Port, "port", ":6661", "-port :6661")

	flag.Parse()
	if PoolID == "" {
		panic("pool_id is empty")
	}

	// register pool net host
	go RegisterPoolInfo()

	go pool()
	select {}

}

const (
	ONLINE      = 1
	OFFLINE     = 3
	RECV_BRIDGE = 5
)

func pool() {
	srv := tcpx.NewTcpX(nil)
	srv.WithBuiltInPool(true)

	srv.AddHandler(ONLINE, online)
	srv.AddHandler(OFFLINE, offline)

	srv.AddHandler(RECV_BRIDGE, recvBridge)

	srv.ListenAndServe("tcp", Port)
}

func online(c *tcpx.Context) {
	type UserLogin struct {
		UserId   int64  `json:"user_id"`
		Password string `json:"password"`
	}

	var userLogin UserLogin

	if _, e := c.Bind(&userLogin); e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		c.JSON(400, tcpx.H{"message": errorx.Wrap(e).Error()})
		return
	}
	c.Online(strconv.FormatInt(userLogin.UserId, 10))

	RegisterInfo(userLogin.UserId)

	c.JSON(200, tcpx.H{"message": "success"})
}

func offline(c *tcpx.Context) {
	type UserOffline struct {
		UserId int64 `json:"user_id"`
	}
	var userOffline UserOffline
	if _, e := c.Bind(&userOffline); e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		c.JSON(400, tcpx.H{"message": errorx.Wrap(e).Error()})
		return
	}
	RegisterRemoveInfo(userOffline.UserId)
	c.Offline()
	c.JSON(200, tcpx.H{"message": "success"})
}

func recvBridge(c *tcpx.Context) {
	type Param struct {
		EventName string `json:"event_name"`
		EventType string `json:"event_type"`
		Stream    []byte `json:"stream"`
	}
	var param Param
	if _, e := c.Bind(&param); e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		c.JSON(400, gin.H{"message": errorx.Wrap(e).Error()})
		return
	}

	if param.EventType != "keep_alive" {
		defer c.CloseConn()
	}
	switch param.EventName {
	case "text-sending":
		type TextMessage struct {
			FromUser int64  `json:"from_user"`
			ToUser   int64  `json:"to_user"`
			Content  string `json:"content"`
		}
		var message TextMessage
		e := json.Unmarshal(param.Stream, &message)
		if e != nil {
			fmt.Println(errorx.Wrap(e).Error())
			return
		}
		c.SendToUsername(strconv.FormatInt(message.ToUser, 10), 10, message)
	}
}

func RegisterInfo(userId int64) {
	type Request struct {
		UserId int64  `json:"user_id"`
		PoolId string `json:"pool_id"`
	}

	var req = Request{
		UserId: userId,
		PoolId: PoolID,
	}
	buf, e := json.Marshal(req)
	if e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		return
	}
	resp, e := http.Post("http://localhost:6666/register/online-info/", "application/json", bytes.NewReader(buf))
	if resp != nil && resp.Body != nil {
		fmt.Println(resp.Status)
		b, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			fmt.Println(errorx.Wrap(e).Error())
			return
		}
		fmt.Println(string(b))
		return
	}
	fmt.Println(errorx.NewFromString("resp nil"))
}
func RegisterRemoveInfo(userId int64) {
	type Request struct {
		UserId int64 `json:"user_id"`
	}

	var req = Request{
		UserId: userId,
	}
	buf, e := json.Marshal(req)
	if e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		return
	}
	resp, e := http.Post("http://localhost:6666/register/offline-info/", "application/json", bytes.NewReader(buf))
	if resp != nil && resp.Body != nil {
		fmt.Println(resp.Status)
		b, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			fmt.Println(errorx.Wrap(e).Error())
			return
		}
		fmt.Println(string(b))
		return
	}
	fmt.Println(errorx.NewFromString("resp nil"))
}
func RegisterPoolInfo() {
	fmt.Println(">> register pool host")
	type Request struct {
		PoolId string `json:"pool_id"`
		Host   string `json:"host"`
	}

	var req = Request{
		PoolId: PoolID,
		Host:   fmt.Sprintf("%s%s", Host, Port),
	}
	buf, e := json.Marshal(req)
	if e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		return
	}
	resp, e := http.Post("http://localhost:6666/register/add-pool-info/", "application/json", bytes.NewReader(buf))
	if resp != nil && resp.Body != nil {
		fmt.Println(resp.Status)
		b, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			fmt.Println(errorx.Wrap(e).Error())
			return
		}
		fmt.Println(string(b))
		return
	}
	fmt.Println(errorx.NewFromString("resp nil"))
}
func randomString() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var r *rand.Rand

	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	var result string
	for i := 0; i < 32; i++ {
		result += string(str[r.Intn(len(str))])
	}
	return result
}
