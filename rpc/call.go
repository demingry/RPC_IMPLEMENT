package rpc

type Call struct {
	Seq           uint64
	ServiceMethod string
	Arguments     interface{}
	Reply         interface{}
	Error         error
	Done          chan struct{}
}
