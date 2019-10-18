package main

import (
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"net"
	"tcpx"
)

type Struct200 struct {
	Message string `json:"message"`
}
type Struct400 struct {
	Message string `json:"message"`
}
type Struct500 struct {
	Message string `json:"message"`
}

type StructTextSending struct {
	FromUser int64  `json:"from_user"`
	ToUser   int64  `json:"to_user"`
	Content  string `json:"content"`
}
type UserLogin struct {
	UserId   int64  `json:"user_id"`
	Password string `json:"password"`
}

type TextMessage struct {
	FromUser int64  `json:"from_user"`
	ToUser   int64  `json:"to_user"`
	Content  string `json:"content"`
}

func main() {
	fmt.Println("client1 start")

	// client message handler
	binder := tcpx.NewClientBinder()
	{
		binder.AddHandler(200, &Struct200{}, func(modelPtr interface{}) {
			fmt.Println(modelPtr.(*Struct200).Message)
		})
		binder.AddHandler(400, &Struct400{}, func(modelPtr interface{}) {
			fmt.Println(modelPtr.(*Struct400).Message)
		})
		binder.AddHandler(500, &Struct500{}, func(modelPtr interface{}) {
			fmt.Println(modelPtr.(*Struct500).Message)
		})
		binder.AddHandler(10, &StructTextSending{}, func(modelPtr interface{}) {
			fmt.Println(modelPtr.(*StructTextSending).Content)
		})
	}

	// connect to user pool server
	poolConn, e := net.Dial("tcp", "localhost:6661")
	{
		if e != nil {
			panic(e)
		}

		binder.HandleConn(poolConn)

		onlineBuf, e := tcpx.PackWithMarshallerName(
			tcpx.Message{
				MessageID: 1,
				Body: UserLogin{
					UserId:   17492019,
					Password: "xxx",
				},
			},
			tcpx.JSON,
		)
		if e != nil {
			panic(e)
		}
		if _, e := poolConn.Write(onlineBuf); e != nil {
			fmt.Println(errorx.Wrap(e).Error())
			return
		}
	}

	// connect to center server
	centerConn, e := net.Dial("tcp", "localhost:8999")
	{
		if e != nil {
			panic(e)
		}

		binder.HandleConn(centerConn)

		textMessageBuf, e := tcpx.PackWithMarshallerName(
			tcpx.Message{
				MessageID: 101,
				Body: TextMessage{
					ToUser:   17492020,
					FromUser: 17492019,
					Content:  "hello，I'm a message from 17492019",
				},
			},
			tcpx.JSON,
		)
		if e != nil {
			panic(e)
		}
		if _, e := centerConn.Write(textMessageBuf); e != nil {
			fmt.Println(errorx.Wrap(e).Error())
			return
		}
	}

	select {}
}
