package tcpx

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"math/rand"
	"strings"
	"sync"
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

func RandomString() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var r *rand.Rand
	var once = sync.Once{}
	once.Do(func() {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	})
	var result string
	for i := 0; i < 32; i++ {
		result += string(str[r.Intn(len(str))])
	}
	return result
}

func RetryHandler(n int, f func() (bool, error)) error {
	ok, er := f()
	if ok && er == nil {
		return nil
	}
	if n-1 > 0 {
		return RetryHandler(n-1, f)
	}
	return er
}

// Execute f() n times, each time on fail will keep specific interval in order of `secondIntervals`
// If n < 0, it will forever execute f(), interval will used in order of `secondIntervals` and keep using the last interval element at last.
// Work in serial.
func RetryHandlerWithInterval(n int, f func() (bool, error), secondIntervals ... int) error {
	var offset = 0
	var ok bool
	var e error
	return retryHandlerWithIntervalOffset(&ok, &e, n, f, &offset, secondIntervals...)
}

// sub function of RetryHandlerWithInterval
func retryHandlerWithIntervalOffset(ok *bool, e *error, n int, f func() (bool, error), offset *int, secondIntervals ... int) error {
	if n == 0 {
		return nil
	}
	*ok, *e = f()
	if *ok && *e == nil {
		return nil
	}

	if n-1 > 0 && !(n < 0) {
		if *offset < len(secondIntervals) {
			time.Sleep(time.Duration(secondIntervals[*offset]) * time.Second)
			if *offset < len(secondIntervals)-1 {
				*offset ++
			}
		}
		return retryHandlerWithIntervalOffset(ok, e, n-1, f, offset, secondIntervals...)
	} else {
		if n != 1 {
			if *offset < len(secondIntervals) {
				time.Sleep(time.Duration(secondIntervals[*offset]) * time.Second)
				if *offset < len(secondIntervals)-1 {
					*offset ++
				}
			}
			return retryHandlerWithIntervalOffset(ok, e, n, f, offset, secondIntervals...)
		}
	}
	return *e
}
