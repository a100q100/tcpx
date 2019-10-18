package tcpx
// Unstable, do not use
import (
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"io"
	"net"
	"sync"
	"time"
)

type Option struct {
	Marshaller Marshaller

	requestMap sync.Map
	offset     int
	l          *sync.RWMutex

	host, network string
	autoConnect   bool
	pingInterval  time.Duration
	conn          net.Conn

	reconnect chan struct{}
}

func OptionTODO(network string, host string) *Option {
	return &Option{
		requestMap: sync.Map{},
		offset:     0,
		l:          &sync.RWMutex{},
		host:       host,
		network:    network,

		reconnect: make(chan struct{}, 10),
	}
}
func (o *Option) SetMarshaller(m Marshaller) {
	if m == nil {
		m = JsonMarshaller{}
	}
	o.Marshaller = m
}

func (o *Option) AutoConnect(auto bool, pingInterval time.Duration) *Option {
	o.autoConnect = auto
	o.pingInterval = pingInterval
	return o
}

func (o *Option) InitConnection() error {
	conn, e := net.Dial(o.network, o.host)

	if e != nil {
		return errorx.Wrap(e)
	}
	o.conn = conn

	go func() {
		for {
			stream, e := FirstBlockOf(o.conn)
			if e != nil {
				if e == io.EOF {
					if o.autoConnect == true {
						<-o.reconnect
						continue
					} else {
						return
					}
				}
				fmt.Println(errorx.Wrap(e).Error())
			}
			go func(stream []byte) {
				requestId, e := RequestIDOf(stream)
				if e != nil {
					fmt.Println(errorx.Wrap(e).Error())
					return
				}
				if requestId == "" {
					return
				}
				tmp, ok := o.requestMap.Load(requestId)
				if !ok {
					fmt.Println("not found result chanel for requestID " + requestId)
					return
				}
				c, ok := tmp.(chan []byte)
				if !ok {
					fmt.Println("bad chan type for requestID" + requestId)
					return
				}
				c <- stream
			}(stream)
		}
	}()

	if o.autoConnect == true {
		go func() {
			ping := PackStuff(DEFAULT_PING_MESSAGEID)
			reConnect := func() (bool, error) {
				o.conn, e = net.Dial(o.network, o.host)
				if e != nil {
					return false, e
				}
				return true, nil
			}

			for {
				_, e := o.conn.Write(ping)
				if e != nil {
					// Logger.Println(errorx.Wrap(e))
					e := RetryHandlerWithInterval(-1, reConnect, 3, 3, 10, 10, 30, 30, 30, 60, 60, 5*60)
					if e == nil {
						o.reconnect <- struct{}{}
					}
				}

				time.Sleep(o.pingInterval)
				continue
			}
		}()
	}
	return nil
}
func Call(request Message, option *Option) ([]byte, error) {
	option.l.Lock()
	option.offset ++
	requestId := fmt.Sprintf("%s://%s/%d", option.network, option.host, option.offset)
	option.l.Unlock()

	request.Set("tcpx-request-id", requestId)

	var result = make(chan []byte, 1)
	option.requestMap.Store(requestId, result)

	buf, e := PackWithMarshaller(request, option.Marshaller)
	if e != nil {
		return nil, errorx.Wrap(e)
	}
	option.conn.Write(buf)

	select {
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("time out")
	case v := <-result:
		return v, nil
	}
}
