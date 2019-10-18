package tcpx

import (
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"io"
	"net"
)

// Help go client divide message response
// Old handle style:
/*
    for {
        buf,e := tcpx.FirstBlockOf(conn)
        messageID, _ := tcpx.MessageIDOf(buf)
        bodyBytes,_ := tcpx.BodyBytesOf(buf)
        switch messageID {
            case 200:
            case 400:
            case 500:
            case 2:
            case 9:
                type Struct9 struct{
                    Message string `json:"message"`
                    ShopName string `json:"shop_name"`
                    ...
                }
                var struct9 Struct9
                json.Unmarshal(bodyBytes, &struct9)
                ...
        }
    }
 */
// Now client-bind provides better style:
/*
    binder := tcpx.NewBinder()
    var struct9 Struct9
    binder.AddHandler(9, &struct9,func(model interface{}){
        response, ok := model.(Struct9)
        if !ok {
            panic("type unmatch")
        }
    })
    binder.HandleConn(conn)
*/

type ClientBinder struct {
	ModelPtrMap map[int32]interface{}
	HandlerMap  map[int32]func(interface{})
}

func NewClientBinder() *ClientBinder {
	return &ClientBinder{
		ModelPtrMap: make(map[int32]interface{}),
		HandlerMap:  make(map[int32]func(interface{})),
	}
}
func (cb *ClientBinder) AddHandler(messageID int32, modelPtr interface{}, f func(modelPtr interface{})) {
	cb.ModelPtrMap[messageID] = modelPtr
	cb.HandlerMap[messageID] = f
}

// Handle raw tcpx block stream
func (cb ClientBinder) Handle(stream []byte) error {
	messageID, e := MessageIDOf(stream)
	if e != nil {
		return errorx.Wrap(e)
	}
	body, e := BodyBytesOf(stream)
	if e != nil {
		return errorx.Wrap(e)
	}
	header, e := HeaderOf(stream)
	if e != nil {
		return errorx.Wrap(e)
	}

	var marshaller Marshaller
	if header == nil || len(header) == 0 {
		marshaller = JsonMarshaller{}
	} else {
		if header["tcpx-marshal-name"] == nil {
			marshaller = JsonMarshaller{}
		} else {
			marshaller = GetMarshallerByMarshalNameDefaultJSON(header["tcpx-marshal-name"].(string))
		}
	}
	modelPtr, ok := cb.ModelPtrMap[messageID]
	if !ok {
		Logger.Println(fmt.Sprintf("ignored, client binder not found modelPtr for messagID '%d'", messageID))
		return nil
	}
	e = marshaller.Unmarshal(body, modelPtr)
	if e != nil {
		return errorx.Wrap(e)
	}
	handler, ok := cb.HandlerMap[messageID]
	if !ok {
		Logger.Println(fmt.Sprintf("ignored, client binder not found handler for messagID '%d'", messageID))
		return nil
	}
	handler(modelPtr)
	return nil
}

// Handle conn
func (cb ClientBinder) HandleConn(conn net.Conn) {
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
			cb.Handle(stream)
		}
	}()
}

