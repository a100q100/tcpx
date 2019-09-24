package tcpx

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"math/rand"
	"strings"
	"time"
)

type H map[string]interface{}

func Debug(src interface{}) string {
	buf, e := json.MarshalIndent(src, "  ", "  ")
	if e != nil {
		fmt.Println(errorx.Wrap(e).Error())
	}
	return string(buf)
}

// Whether s in arr
// Support %%
func In(s string, arr []string) bool {
	for _, v := range arr {
		if strings.Contains(v, "%") {
			if strings.HasPrefix(v, "%") && strings.HasSuffix(v, "%") {
				if strings.Contains(s, string(v[1:len(v)-1])) {
					return true
				}
			} else if strings.HasPrefix(v, "%") {
				if strings.HasSuffix(s, string(v[1:])) {
					return true
				}
			} else if strings.HasSuffix(v, "%") {
				if strings.HasPrefix(s, string(v[:len(v)-1])) {
					return true
				}
			}
		} else {
			if v == s {
				return true
			}
		}
	}
	return false
}

// Defer eliminates all panic cases and handle panic reason by handlePanicError
func Defer(f func(), handlePanicError ...func(interface{})) {
	defer func() {
		if e := recover(); e != nil {
			for _, handler := range handlePanicError {
				handler(e)
			}
		}
	}()
	f()
}

// CloseChanel(func(){close(chan)})
func CloseChanel(f func()) {
	defer func() {
		if e := recover(); e != nil {
			// when close(chan) panic from 'close of closed chan' do nothing
		}
	}()
	f()
}
func MD5(rawMsg string) string {
	data := []byte(rawMsg)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has)
	return strings.ToUpper(md5str1)
}

func RandomString(length int) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var result string
	for i := 0; i < length; i++ {
		result += string(str[r.Intn(len(str))])
	}
	return result
}
