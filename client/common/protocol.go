package common

import (
	"io"
	"net"
)

const (
	ackByteSize = 1
	ackSuccessCode = 1
	ackIndex = 0
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

func (p *ClientProtocol) SendBatch(payload []byte) error {
	return p.Send(payload)
}

func (p *ClientProtocol) ReadAck() (bool, error) {
	ack := make([]byte, ackByteSize)
	if _, err := io.ReadFull(p.conn, ack); err != nil {
		return false, err
	}
	return ack[ackIndex] == ackSuccessCode, nil
}

func (p *ClientProtocol) ReadBatchAck() (bool, error) {
	return p.ReadAck()
}

func (p *ClientProtocol) Close() {
	if p.conn != nil {
		_ = p.conn.Close()
		p.conn = nil
	}
}
