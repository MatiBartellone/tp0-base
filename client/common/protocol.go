package common

import (
	"io"
	"net"
)

type ClientProtocol struct {
	conn net.Conn
}

func NewClientProtocol(conn net.Conn) *ClientProtocol {
	return &ClientProtocol{conn: conn}
}

func (p *ClientProtocol) Send(payload []byte) error {
	total := 0
	for total < len(payload) {
		n, err := p.conn.Write(payload[total:])
		if err != nil {
			return err
		}
		total += n
	}
	return nil
}

func (p *ClientProtocol) ReadAck() (bool, error) {
	ack := make([]byte, 1)
	if _, err := io.ReadFull(p.conn, ack); err != nil {
		return false, err
	}
	return ack[0] == 1, nil
}

func (p *ClientProtocol) Close() {
	if p.conn != nil {
		_ = p.conn.Close()
		p.conn = nil
	}
}
