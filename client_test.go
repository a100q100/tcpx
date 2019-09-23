package tcpx

import (
	"fmt"
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
			for {
				n, e = c.ConnReader.Read(buf)
				if e != nil {
					fmt.Println(e.Error())
					return
				}
				// fmt.Println("receive:", string(buf[:n]))
				 _ = n
				c.ConnWriter.Write([]byte("hello,I am server."))
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

	for i := 0; i < 15; i++ {
		go func(){
			resp, e := Call([]byte("hello"), option)
			if e != nil {
				fmt.Println(e.Error())
				return
			}
			fmt.Println(resp)
			time.Sleep(1 * time.Second)
		}()
	}

	time.Sleep(15 * time.Second)
}
