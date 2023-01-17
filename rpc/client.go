package rpc

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type Client struct {
	Codec   Codec
	Mu      sync.Mutex
	Pending map[uint64]*Call
	Conn    net.Conn
	Closed  bool
}

func NewClient(addr string, deadline time.Duration) *Client {
	if conn, err := net.DialTimeout("tcp", addr, deadline); err == nil {

		json.NewEncoder(conn).Encode(Options{Indicate: 0, CodecType: GobCodec})
		var options Options
		json.NewDecoder(conn).Decode(&options)
		conn.Close()
		if options.Indicate == 0 {
			if conn, err = net.DialTimeout("tcp", fmt.Sprintf(
				"%s:%d",
				options.Single.Address,
				options.Single.Port), deadline); err == nil {
				json.NewEncoder(conn).Encode(Options{Indicate: 1, CodecType: GobCodec})
				if init, ok := InitCodecLists[options.CodecType]; ok {
					codec := init(conn)
					return &Client{
						Codec:   codec,
						Mu:      sync.Mutex{},
						Pending: make(map[uint64]*Call),
						Conn:    conn,
						Closed:  true,
					}
				}
			}
		}
	}

	return nil
}

func (c *Client) RegisterCall(call *Call) {

	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.Pending[call.Seq] = call

	c.Codec.Write(&Header{
		Sequence:      call.Seq,
		ServiceMethod: call.ServiceMethod}, call.Arguments)
}

func (c *Client) RemoveCall(sequence uint64) {

	c.Mu.Lock()
	defer c.Mu.Unlock()
	delete(c.Pending, sequence)
}

func (c *Client) Receive() {
	for {
		var h Header
		if err := c.Codec.ReadHeader(&h); err != nil {
			break
		}

		if h.Error != nil {
			c.Codec.ReadBody(nil)
			c.RemoveCall(h.Sequence)
			continue
		} else {
			call := c.Pending[h.Sequence]
			var reply int
			if err := c.Codec.ReadBody(&reply); err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				continue
			}
			call.Reply = reply
			c.RemoveCall(h.Sequence)

			fmt.Printf("Call sequence number: %d with result %v\n", call.Seq, call.Reply)
		}
	}
}
