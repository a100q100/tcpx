package tcpx

import (
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"net"
	"strconv"
	"sync"
	"time"
)

// Call args
// Option is a context used in Call().It provide url connection cache and can config client configurations like timeout, keepAlive...
type Option struct {
	// Help unique request_id
	Salt string

	// key: network://host
	// Cache helps store alive connections
	Cache map[string]net.Conn
	l     *sync.RWMutex
	// key: network://host:randomStr
	// RequestMap helps define request_id for the same connection
	RequestMap map[string]chan []byte
	l2         *sync.RWMutex

	Network string
	Host    string

	Timeout    time.Duration
	Marshaller Marshaller

	KeepAlive bool
	AliveTime time.Duration
}

// init an empty option
func OptionTODO() *Option {
	return &Option{
		Cache:      make(map[string]net.Conn, 0),
		l:          &sync.RWMutex{},
		RequestMap: make(map[string]chan []byte, 0),
		l2:         &sync.RWMutex{},
	}
}

// Set option network and host
func (o *Option) SetNetworkHost(network, host string) *Option {
	o.Network = network
	o.Host = host
	return o
}

// Set option timout
func (o *Option) SetTimeout(timeout time.Duration) *Option {
	o.Timeout = timeout
	return o
}

// Set option keepAlive
func (o *Option) SetKeepAlive(alive bool, timeout time.Duration) *Option {
	o.KeepAlive = alive
	o.AliveTime = timeout
	return o
}

// Use option's non-empty value and set it into o
func (o *Option) Option(option Option) *Option {
	if option.Host != "" {
		o.Host = option.Host
	}
	if option.Network != "" {
		o.Network = option.Network
	}
	if option.Timeout != 0 {
		o.Timeout = option.Timeout
	}

	if option.KeepAlive != false {
		o.KeepAlive = option.KeepAlive
	}
	if option.Marshaller != nil {
		o.Marshaller = option.Marshaller
	}

	if option.AliveTime != 0 {
		o.AliveTime = option.AliveTime
	}
	return o
}

// Copy an option instance for different request
func (o *Option) Copy() *Option {
	return &Option{
		Network: o.Network,
		Host:    o.Host,
		Cache:   o.Cache,
		Timeout: o.Timeout,

		Marshaller: o.Marshaller,

		KeepAlive: o.KeepAlive,
		AliveTime: o.AliveTime,
	}
}

// Check an option has existing host connection cache
func (o *Option) HasCache() bool {
	o.l.RLock()
	defer o.l.RUnlock()
	return len(o.Cache) > 0
}

// Call require client send a request and server response once
//
// Deprecated: unstable, do not use.
func Call(request []byte, option *Option) ([]byte, error) {
	if option.Salt == "" {
		option.Salt = "tcpx"
	}
	connHash := fmt.Sprintf("%s://%s", option.Network, option.Host)
	requestHash := fmt.Sprintf("%s:%s:%s:%s", connHash, option.Salt, strconv.FormatInt(time.Now().UnixNano(), 10), RandomString(32))
	var e error
	var conn net.Conn

	// get conn from pool if exist, otherwise new tcp connection and put it into cache
	var ok bool
	func() {
		option.l.Lock()
		defer option.l.Unlock()
		conn, ok = option.Cache[connHash]
		if !ok {
			conn, e = net.Dial(option.Network, option.Host)
			if e != nil {
				fmt.Println(errorx.Wrap(e))
				return
			}
			option.Cache[connHash] = conn

			go func() {
				buf, e := FirstBlockOf(conn)
				if e != nil {
					fmt.Println(errorx.Wrap(e))
					return
				}
				option.l2.RLock()
				option.RequestMap[requestHash] <- buf
				option.l2.RUnlock()
				return
			}()
		}
	}()

	// init response chanel
	option.l2.Lock()
	option.RequestMap[requestHash] = make(chan []byte, 1)
	option.l2.Unlock()
	defer func() {
		option.l2.Lock()
		defer option.l2.Unlock()
		delete(option.RequestMap, requestHash)
	}()

	//// If keep connection alive, after alive duration, connection will be closed and remove from cache
	//{
	//	if option.KeepAlive == true {
	//		go func() {
	//			for {
	//				select {
	//				case <-time.After(option.AliveTime):
	//					option.l.Lock()
	//					delete(option.Cache, connHash)
	//					option.l.Unlock()
	//					conn.Close()
	//				}
	//			}
	//		}()
	//	}
	//}

	// start a goroutine to receive income response

	// write request bytes
	if _, e := conn.Write(request); e != nil {
		fmt.Println(errorx.Wrap(e))
		return nil, errorx.Wrap(e)
	}

	// If timeout == 0, stuck until result has value , otherwise after timeout interval, return time-out error
	var c chan []byte
	option.l2.RLock()
	c = option.RequestMap[requestHash]
	option.l2.RUnlock()
	if option.Timeout == 0 {
		return <-c, nil
	} else {
		select {
		case <-time.After(option.Timeout):
			return nil, fmt.Errorf("time out")
		case v := <-c:
			return v, nil
		}
	}
}
