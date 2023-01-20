package compress

type Compress interface {
	Zip([]byte) []byte
	Unzip([]byte) []byte
}

const (
	Raw int = iota
	Gzip
	Zlib
)

func NewCompressor(selector int) Compress {

	switch selector {
	case Gzip:
		return GzipCompress{}
	case Zlib:
		return ZlibCompress{}
	default:
		return nil
	}

}
