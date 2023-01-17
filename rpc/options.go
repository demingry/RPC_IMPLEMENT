package rpc

type Options struct {
	Indicate  int //indicate 0: get backend server; indicate 1: common cummunication
	CodecType CodecType
	Single    Single
}
