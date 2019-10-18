package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestOnline(t *testing.T) {
    rs,e:= RegisterGetUserInfo(17492019)
    fmt.Println(rs,e)
}

type UserInfo struct {
	UserId int64  `json:"user_id" binding:"required"` // 用户id，全局唯一
	PoolId string `json:"pool_id" binding:"required"` // 用户池id
}

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
		fmt.Println(resp.Status)
		b, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			fmt.Println(errorx.Wrap(e).Error())
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
