package compress

import (
	"bytes"
	"compress/gzip"
)

type GzipCompress struct{}

func (GzipCompress) Zip(data []byte) []byte {

	buf := bytes.NewBuffer(nil)
	w := gzip.NewWriter(buf)
	defer w.Close()

	_, _ = w.Write(data)
	_ = w.Flush()

	return buf.Bytes()
}

func (GzipCompress) Unzip(data []byte) []byte {

	buf := bytes.NewBuffer(nil).Bytes()
	r, _ := gzip.NewReader(bytes.NewReader(data))
	_, _ = r.Read(buf)

	return buf
}
