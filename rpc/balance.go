package rpc

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type Colony struct {
	mu     sync.Mutex
	index  int
	Single []*Single
}

func (c *Colony) Update(object []*Single) {

	c.Single = make([]*Single, len(object))
	total := copy(c.Single, object)

	fmt.Fprintf(os.Stdout, "Total updated server: %d\n", total)
}

func (c *Colony) HeartBeat(max time.Duration) {

	for {
		for i, o := range c.Single {
			if o.Max_retry_times > 10 {
				fmt.Fprintf(os.Stderr, "[!]some backend server failed\n")
				c.mu.Lock()
				c.Single = append(c.Single[:i], c.Single[i+1:]...)
				c.mu.Lock()
				continue
			}
			success := make(chan struct{}, 1)
			go func() {
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", o.Address, o.Port))
				success <- struct{}{}
				if err == nil {
					defer conn.Close()
				}
			}()
			select {
			case <-time.After(max):
				c.Single[i].Max_retry_times += 1
			case <-success:
				continue
			}
		}
	}

}

type Single struct {
	Address         string
	Port            uint16
	Max_retry_times int
}
