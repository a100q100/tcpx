package tcpx

import (
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"io"
	"net"
	"testing"
	"time"
)

func TestBinderJSON(t *testing.T) {
	var serverStart = make(chan int, 1)
	var testResult = make(chan error, 1)
	// middlewareOrder suggest the execute order of three kinds middleware [1,2,3]
	var middlewareOrder = make([]int, 0, 10)

	go func() {
		time.Sleep(30 * time.Second)
		testResult <- nil
	}()
	// client
	go func() {
		<-serverStart

		conn, err := net.Dial("tcp", "localhost:8001")
		if err != nil {
			testResult <- errorx.Wrap(err)
			fmt.Println(errorx.Wrap(err).Error())
			return
		}
		// recv
		go func() {
			binder := NewClientBinder()
			var rs string
			binder.AddHandler(10086, &rs, func(modelPtr interface{}) {
				tmp := modelPtr.(*string)
				fmt.Println(*tmp)
			})
			for {
				stream, e := FirstBlockOf(conn)
				if e == io.EOF {
					return
				}
				if e != nil {
					fmt.Println(errorx.Wrap(e).Error())
					return
				}
				binder.Handle(stream)
			}
		}()

		buf, e := PackJSON.Pack(1, "hello, I'm client")

		if e != nil {
			testResult <- errorx.Wrap(e)
			fmt.Println(errorx.Wrap(e).Error())
			return
		}
		conn.Write(buf)
	}()

	// server
	go func() {
		srv := NewTcpX(JsonMarshaller{})
		srv.OnMessage = nil
		srv.BeforeExit(func() {
			fmt.Println("exit")
		})
		// router middleware
		srv.AddHandler(1, func(c *Context) {
			middlewareOrder = append(middlewareOrder, 3)
		}, func(c *Context) {
			c.Reply(10086, "hello, I'm server")
		})

		go func() {
			time.Sleep(time.Second * 10)
			serverStart <- 1
		}()
		e := srv.ListenAndServeTCP("tcp", ":8001")
		if e != nil {
			testResult <- errorx.Wrap(e)
			fmt.Println(e.Error())
			return
		}
	}()

	e := <-testResult
	if e != nil {
		fmt.Println(e.Error())
		t.Fail()
	}
}

func TestBinderProtobuf(t *testing.T) {
	var serverStart = make(chan int, 1)
	var testResult = make(chan error, 1)
	// middlewareOrder suggest the execute order of three kinds middleware [1,2,3]
	var middlewareOrder = make([]int, 0, 10)

	go func() {
		time.Sleep(30 * time.Second)
		testResult <- nil
	}()
	// client
	go func() {
		<-serverStart

		conn, err := net.Dial("tcp", "localhost:8001")
		if err != nil {
			testResult <- errorx.Wrap(err)
			fmt.Println(errorx.Wrap(err).Error())
			return
		}
		binder := NewClientBinder()
		var rs string
		binder.AddHandler(10086, &rs, func(modelPtr interface{}) {
			tmp := modelPtr.(*string)
			fmt.Println(*tmp)
		})
		var hello Hello
		binder.AddHandler(10087, &hello, func(modelPtr interface{}) {
			tmp := modelPtr.(*Hello)
			fmt.Println(tmp.Message)
		})
		// recv
		go func() {

			for {
				stream, e := FirstBlockOf(conn)
				if e == io.EOF {
					return
				}
				if e != nil {
					fmt.Println(errorx.Wrap(e).Error())
					return
				}
				binder.Handle(stream)
			}
		}()

		buf, e := PackJSON.Pack(1, "hello, I'm client")

		if e != nil {
			testResult <- errorx.Wrap(e)
			fmt.Println(errorx.Wrap(e).Error())
			return
		}
		conn.Write(buf)
	}()

	// server
	go func() {
		srv := NewTcpX(JsonMarshaller{})
		srv.OnMessage = nil
		srv.BeforeExit(func() {
			fmt.Println("exit")
		})
		// router middleware
		srv.AddHandler(1, func(c *Context) {
			middlewareOrder = append(middlewareOrder, 3)
		}, func(c *Context) {
			var hello = Hello{
				Message: "hello, I'm server",
			}
			c.JSON(10086, "hello, I'm server")
			c.ProtoBuf(10087, &hello, map[string]interface{}{
				"tcpx-marshal-name": "protobuf",
			})
		})

		go func() {
			time.Sleep(time.Second * 10)
			serverStart <- 1
		}()
		e := srv.ListenAndServeTCP("tcp", ":8002")
		if e != nil {
			testResult <- errorx.Wrap(e)
			fmt.Println(e.Error())
			return
		}
	}()

	e := <-testResult
	if e != nil {
		fmt.Println(e.Error())
		t.Fail()
	}
}

func TestClientBinder_HandleConn(t *testing.T) {
	var serverStart = make(chan int, 1)
	var testResult = make(chan error, 1)
	// middlewareOrder suggest the execute order of three kinds middleware [1,2,3]
	var middlewareOrder = make([]int, 0, 10)

	go func() {
		time.Sleep(30 * time.Second)
		testResult <- nil
	}()
	// client
	go func() {
		<-serverStart

		conn, err := net.Dial("tcp", "localhost:8001")
		if err != nil {
			testResult <- errorx.Wrap(err)
			fmt.Println(errorx.Wrap(err).Error())
			return
		}
		binder := NewClientBinder()
		var rs string
		binder.AddHandler(10086, &rs, func(modelPtr interface{}) {
			tmp := modelPtr.(*string)
			fmt.Println(*tmp)
		})
		var hello Hello
		binder.AddHandler(10087, &hello, func(modelPtr interface{}) {
			tmp := modelPtr.(*Hello)
			fmt.Println(tmp.Message)
		})
		binder.HandleConn(conn)

		buf, e := PackJSON.Pack(1, "hello, I'm client")

		if e != nil {
			testResult <- errorx.Wrap(e)
			fmt.Println(errorx.Wrap(e).Error())
			return
		}
		conn.Write(buf)
	}()

	// server
	go func() {
		srv := NewTcpX(JsonMarshaller{})
		srv.OnMessage = nil
		srv.BeforeExit(func() {
			fmt.Println("exit")
		})
		// router middleware
		srv.AddHandler(1, func(c *Context) {
			middlewareOrder = append(middlewareOrder, 3)
		}, func(c *Context) {
			var hello = Hello{
				Message: "hello, I'm server",
			}
			c.JSON(10086, "hello, I'm server")
			c.ProtoBuf(10087, &hello, map[string]interface{}{
				"tcpx-marshal-name": "protobuf",
			})
		})

		go func() {
			time.Sleep(time.Second * 10)
			serverStart <- 1
		}()
		e := srv.ListenAndServeTCP("tcp", ":8003")
		if e != nil {
			testResult <- errorx.Wrap(e)
			fmt.Println(e.Error())
			return
		}
	}()

	e := <-testResult
	if e != nil {
		fmt.Println(e.Error())
		t.Fail()
	}
}
