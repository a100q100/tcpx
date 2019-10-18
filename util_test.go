package tcpx

import (
	"fmt"
	"testing"
	"time"
)

func TestIn(t *testing.T) {
	// true
	fmt.Println(In("ipskzk", []string{"ip%"}))
	fmt.Println(In("ipskzk", []string{"%ip%"}))
	fmt.Println(In("kjlk;lip", []string{"%ip"}))
	fmt.Println(In("kjlk;lip", []string{"%ip%"}))
}

func TestDebug(t *testing.T) {
	fmt.Println(Debug("hello"))
}

func TestDefer(t *testing.T) {
	f := func() {
		fmt.Println(1)
		panic(1)
	}
	Defer(f, func(v interface{}) {
		fmt.Println(v)
	})
}

func TestRetryHandlerWithInterval(t *testing.T) {
	t1 := time.Now()
	var times int
	RetryHandlerWithInterval(-1, func() (bool, error) {
		times ++
		if times == 3 {
			return true, nil
		}
		fmt.Println("exec")
		return false, fmt.Errorf("nil return")
	}, )
	t2 := time.Now()
	fmt.Println(t2.Sub(t1).Seconds())
}
