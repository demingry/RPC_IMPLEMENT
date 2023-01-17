package rpc

import "io"

type Codec interface {
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type CodecType string

var (
	JsonCodec CodecType = "JSON"
	GobCodec  CodecType = "GOB"
)

type NewCodecFunc func(io.ReadWriteCloser) Codec

//编解码器实例初始化列表
var InitCodecLists map[CodecType]NewCodecFunc

func init() {
	InitCodecLists = make(map[CodecType]NewCodecFunc)
	InitCodecLists[GobCodec] = NewGob
}

type Header struct {
	ServiceMethod string
	Sequence      uint64
	BodySize      uint64
	Error         error
}
