package compress

import (
	"bytes"
	"compress/zlib"
)

type ZlibCompress struct{}

func (ZlibCompress) Zip(data []byte) []byte {

	buf := bytes.NewBuffer(nil)

	w := zlib.NewWriter(buf)
	defer w.Close()
	_, _ = w.Write(data)

	return buf.Bytes()
}

func (ZlibCompress) Unzip(data []byte) []byte {

	buf := bytes.NewBuffer(nil).Bytes()

	r, _ := zlib.NewReader(bytes.NewReader(data))
	_, _ = r.Read(buf)

	return buf
}
