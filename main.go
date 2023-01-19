package main

import (
	"rpc"
	"time"
)

type Args struct{ Num1, Num2 int }

type Rabbit struct{}

func (r *Rabbit) Add(args Args, re *int) error {

	*re = args.Num1 + args.Num2
	return nil
}

func main() {

	addr := make(chan string)

	server := rpc.Server{}
	rabbit := &Rabbit{}
	server.Register(rabbit)
	s := make([]*rpc.Single, 0)
	s = append(s, &rpc.Single{
		Address:         "[::]",
		Port:            8080,
		Max_retry_times: 10,
	})
	server.Discover(s)
	go server.StartServer(addr)

	client := rpc.NewClient(<-addr, 1*time.Second)
	defer client.Conn.Close()
	go client.Receive()

	for i := 0; i < 5; i++ {
		call := new(rpc.Call)
		call.Seq = uint64(i)
		call.ServiceMethod = "Rabbit.Add"
		call.Arguments = &Args{1, i}

		client.RegisterCall(call)
		time.Sleep(800 * time.Millisecond)
	}

	time.Sleep(3 * time.Second)

}
