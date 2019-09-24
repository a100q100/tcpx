package tcpx

import (
	"fmt"
	"sync/atomic"

	//	"sync/atomic"
	"testing"
	"time"
)

func TestCall(t *testing.T) {
	var serverStart = make(chan int, 0)
	// Start a serve
	go func() {
		srv := NewTcpX(nil)
		srv.OnConnect = func(c *Context) {
			fmt.Println("receive a new connection")
		}
		srv.HandleRaw = func(c *Context) {
			var buf = make([]byte, 500)
			var n int
			var e error
			var times int32
			for {
				n, e = c.ConnReader.Read(buf)
				if e != nil {
					fmt.Println(e.Error())
					return
				}
				// fmt.Println("receive:", string(buf[:n]))
				_ = n
				buf, _ := PackWithMarshallerAndBody(Message{
					MessageID: 1,
					Header:    nil,
				}, []byte("hello"))
				fmt.Println("server times:", atomic.AddInt32(&times, 1))
				c.ConnWriter.Write(buf)
			}
		}
		go func() {
			time.Sleep(2 * time.Second)
			serverStart <- 1
		}()
		srv.ListenAndServeRaw("tcp", ":6633")
	}()

	<-serverStart

	// Init an option and send 15 times request by Call(). Each request start in 1 second
	var option = OptionTODO().
		SetNetworkHost("tcp", "localhost:6633").
		SetTimeout(15 * time.Second).
		SetKeepAlive(true, 20*time.Second)

	var times int32
	for i := 0; i < 15; i++ {
		go func() {
			buf, _ := PackWithMarshallerAndBody(Message{
				MessageID: 1,
				Header:    nil,
			}, []byte("hello"))
			resp, e := Call(buf, option)
			if e != nil {
				fmt.Println(e.Error())
				return
			}
			fmt.Println(resp)
			fmt.Println("send times:", atomic.AddInt32(&times, 1))
			time.Sleep(1 * time.Second)
		}()
	}

	time.Sleep(20 * time.Second)
}
