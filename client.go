package tcpx

import (
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"net"
	"sync"
	"time"
)

// Call args
// Option is a context used in Call().It provide url connection cache and can config client configurations like timeout, keepAlive...
type Option struct {
	Cache map[string]net.Conn
	l     *sync.RWMutex

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
		Cache: make(map[string]net.Conn, 0),
		l:     &sync.RWMutex{},
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
	var result = make(chan []byte, 1)
	var connHash = fmt.Sprintf("%s://%s", option.Network, option.Host)
	var e error
	var conn net.Conn

	// get conn from pool
	if option.HasCache() {
		var ok bool
		option.l.RLock()
		conn, ok = option.Cache[connHash]
		option.l.RUnlock()
		if !ok {
			conn, e = net.Dial(option.Network, option.Host)
			if e != nil {
				return nil, e
			}
			option.l.Lock()
			option.Cache[connHash] = conn
			option.l.Unlock()
		}
	} else {
		conn, e = net.Dial(option.Network, option.Host)
		if e != nil {
			return nil, e
		}
		option.l.Lock()
		option.Cache[connHash] = conn
		option.l.Unlock()

	}
	if e != nil {
		return nil, e
	}

	if option.KeepAlive == true {
		go func() {
			for {
				select {
				case <-time.After(option.AliveTime):
					option.l.Lock()
					delete(option.Cache, connHash)
					option.l.Unlock()
					conn.Close()
				}
			}
		}()
	}
	go func() {
		var buf = make([]byte, 512)
		n, e := conn.Read(buf)
		if e != nil {
			fmt.Println(errorx.Wrap(e))
			return
		}
		result <- buf[:n]
		return
	}()
	conn.Write(request)
	if option.Timeout == 0 {
		v := <-result
		return v, nil
	} else {
		select {
		case <-time.After(option.Timeout):
			return nil, fmt.Errorf("time out")
		case v := <-result:
			return v, nil
		}
	}
}
