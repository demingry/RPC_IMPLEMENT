package rpc

import (
	"encoding/gob"
	"io"
)

type Gob struct {
	enc  gob.Encoder
	dec  gob.Decoder
	conn io.ReadWriteCloser
}

var _ Codec = (*Gob)(nil)

func NewGob(conn io.ReadWriteCloser) Codec {
	return &Gob{
		conn: conn,
		enc:  *gob.NewEncoder(conn),
		dec:  *gob.NewDecoder(conn),
	}
}

func (g *Gob) ReadHeader(h *Header) error {

	if err := g.dec.Decode(h); err != nil {
		return err
	}
	return nil
}

func (g *Gob) Write(h *Header, body interface{}) error {

	if err := g.enc.Encode(h); err != nil {
		return err
	}

	if err := g.enc.Encode(body); err != nil {
		return err
	}

	return nil
}

func (g *Gob) ReadBody(body interface{}) error {

	if err := g.dec.Decode(body); err != nil {
		return err
	}
	return nil
}
