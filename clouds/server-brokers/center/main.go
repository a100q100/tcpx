package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/fwhezfwhez/tcpx"
	"io/ioutil"
	"net"
	"net/http"
)

const (
	EVENT_SEND_TEXT = 101
	// varies events
)

func main() {
	fmt.Println(">> center server starts")
	srv := tcpx.NewTcpX(nil)
	srv.AddHandler(EVENT_SEND_TEXT, EventSendText)

	srv.ListenAndServe("tcp", ":8999")
}

func EventSendText(c *tcpx.Context) {
	type TextMessage struct {
		FromUser int64  `json:"from_user"`
		ToUser   int64  `json:"to_user"`
		Content  string `json:"content"`
	}
	var param TextMessage
	if _, e := c.Bind(&param); e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		c.JSON(400, tcpx.H{"message": errorx.Wrap(e).Error()})
		return
	}
	b, e := json.Marshal(param)
	if e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		c.JSON(500, tcpx.H{"message": errorx.Wrap(e).Error()})
		return
	}
	type Request struct {
		EventName string `json:"event_name"`
		EventType string `json:"event_type"`
		Stream    []byte `json:"stream"`
	}
	req := Request{
		EventName: "text-sending",
		EventType: "once",
		Stream:    b,
	}
	toUserInfo, e := RegisterGetUserInfo(param.ToUser)
	if e == UserOffLineErr {
		c.JSON(400, tcpx.H{"message": "dest user offline"})
		return
	}

	if e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		c.JSON(500, tcpx.H{"message": errorx.Wrap(e).Error()})
		return
	}
	host, e := RegisterGetUserPoolHost(toUserInfo.PoolId)
	if e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		c.JSON(500, tcpx.H{"message": errorx.Wrap(e).Error()})
		return
	}
	conn, e := net.Dial("tcp", host)
	if e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		c.JSON(500, tcpx.H{"message": errorx.Wrap(e).Error()})
		return
	}
	buf, e := tcpx.PackJSON.Pack(5, req)
	if e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		c.JSON(500, tcpx.H{"message": errorx.Wrap(e).Error()})
		return
	}
	if _, e := conn.Write(buf); e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		c.JSON(500, tcpx.H{"message": errorx.Wrap(e).Error()})
		return
	}
	c.JSON(200, tcpx.H{"message": "success"})

}

func RegisterGetUserPoolHost(poolId string) (string, error) {
	type Request struct {
		PoolId string `json:"pool_id"`
	}
	type Response struct {
		PoolId string `json:"pool_id"`
		Host   string `json:"host"`
	}
	var req = Request{
		PoolId: poolId,
	}
	buf, e := json.Marshal(req)
	if e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		return "", errorx.Wrap(e)
	}
	resp, e := http.Post("http://localhost:6666/register/get-pool-info/", "application/json", bytes.NewReader(buf))
	if resp != nil && resp.Body != nil {
		b, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			fmt.Println(errorx.Wrap(e).Error())
			return "", errorx.Wrap(e)
		}
		var rsp Response
		if e := json.Unmarshal(b, &rsp); e != nil {
			fmt.Println(errorx.Wrap(e).Error())
			return "", errorx.Wrap(e)
		}
		return rsp.Host, nil
	}
	return "", errorx.NewFromString("resp nil")
}

type UserInfo struct {
	UserId int64  `json:"user_id" binding:"required"` // 用户id，全局唯一
	PoolId string `json:"pool_id" binding:"required"` // 用户池id
}

var UserOffLineErr = fmt.Errorf("dest user offline")

func RegisterGetUserInfo(userId int64) (UserInfo, error) {
	type Request struct {
		UserId int64 `json:"user_id" binding:"required"` // 用户id，全局唯一
	}

	var req = Request{
		UserId: userId,
	}
	buf, e := json.Marshal(req)
	if e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		return UserInfo{}, errorx.Wrap(e)
	}
	resp, e := http.Post("http://localhost:6666/register/get-info/", "application/json", bytes.NewReader(buf))
	if resp != nil && resp.Body != nil {
		b, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			fmt.Println(errorx.Wrap(e).Error())
			return UserInfo{}, errorx.Wrap(e)
		}

		if resp.StatusCode == 400 {
			type result struct {
				TipId int `json:"tip_id"`
			}
			var rs result
			if e := json.Unmarshal(b, &rs); e != nil {
				fmt.Println(errorx.Wrap(e).Error())
				return UserInfo{}, errorx.Wrap(e)
			}
			if rs.TipId == 1 {
				return UserInfo{}, UserOffLineErr
			}
			return UserInfo{}, errorx.Wrap(e)
		}

		var userInfo UserInfo
		e = json.Unmarshal(b, &userInfo)
		if e != nil {
			fmt.Println(errorx.Wrap(e).Error())
			return UserInfo{}, errorx.Wrap(e)
		}
		fmt.Println(userInfo.UserId, userInfo.PoolId)
		return userInfo, nil
	}
	return UserInfo{}, errorx.NewFromString("resp nil")
}
