package tcpx

import (
	"fmt"
	"sync/atomic"

	//	"sync/atomic"
	"testing"
	"time"
)

func TestAutoConnect(t *testing.T) {
	var serverStart = make(chan int, 0)
	// Start a serve
	go func() {
		srv := NewTcpX(nil)
		srv.OnConnect = func(c *Context) {
			fmt.Println("receive a new connection")
		}
		srv.HandleRaw = func(c *Context) {
			go func() {
				time.Sleep(10 * time.Second)
				fmt.Println("server call conn close")
				c.Conn.Close()
			}()

			var times int32
			for {
				stream, e := FirstBlockOf(c.ConnReader)
				if e != nil {
					fmt.Println(e.Error())
					return
				}
				// fmt.Println("receive:", string(buf[:n]))
				requestID, e := RequestIDOf(stream)
				if e != nil {
					return
				}
				buf, _ := PackWithMarshallerAndBody(Message{
					MessageID: 1,
					Header: map[string]interface{}{
						"tcpx-request-id": requestID,
					},
				}, []byte("hello"))
				fmt.Println("server times:", atomic.AddInt32(&times, 1))
				c.ConnWriter.Write(buf)
			}
		}
		go func() {
			time.Sleep(2 * time.Second)
			serverStart <- 1

			time.Sleep(30 * time.Second)
			srv.Stop(true)

			time.Sleep(30 * time.Second)
			srv.Start()
		}()
		srv.ListenAndServeRaw("tcp", ":6634")
	}()

	<-serverStart

	option := OptionTODO("tcp", "localhost:6634").
		AutoConnect(true, 5*time.Second)

	if e := option.InitConnection(); e != nil {
		fmt.Println(e.Error())
		return
	}

	time.Sleep(2 * time.Minute)
}

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

}
