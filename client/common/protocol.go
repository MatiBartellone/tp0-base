package common

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	queryStatusFailure = 0
	queryStatusSuccess = 1
	queryStatusPending = 2
)

type ClientProtocol struct {
	conn    net.Conn
	builder *ProtocolBuilder
}

func NewClientProtocol(conn net.Conn) *ClientProtocol {
	return &ClientProtocol{conn: conn, builder: NewProtocolBuilder()}
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

func (p *ClientProtocol) SendFinish(agencyID uint8) error {
	payload := p.builder.BuildFinishMessage(agencyID)
	return p.Send(payload)
}

func (p *ClientProtocol) SendWinnersQuery(agencyID uint8) error {
	payload := p.builder.BuildWinnersQueryMessage(agencyID)
	return p.Send(payload)
}

func (p *ClientProtocol) ReadAck() (bool, error) {
	ack, err := p.readU8()
	if err != nil {
		return false, err
	}
	return ack == protocolAckSuccessCode, nil
}

func (p *ClientProtocol) ReadBatchAck() (bool, error) {
	return p.ReadAck()
}

func (p *ClientProtocol) ReadFinishAck() (bool, error) {
	return p.ReadAck()
}

func (p *ClientProtocol) ReadWinnersResponse() (uint8, []uint32, error) {
	status, err := p.readU8()
	if err != nil {
		return queryStatusFailure, nil, err
	}

	count, err := p.readU16()
	if err != nil {
		return queryStatusFailure, nil, err
	}

	if count == 0 {
		return status, nil, nil
	}

	winners := make([]uint32, count)
	for i := range winners {
		winner, err := p.readU32()
		if err != nil {
			return queryStatusFailure, nil, err
		}
		winners[i] = winner
	}

	if status == queryStatusSuccess || status == queryStatusPending {
		return status, winners, nil
	}

	return status, winners, fmt.Errorf("unexpected winners response status: %d", status)
}

func (p *ClientProtocol) readU8() (uint8, error) {
	buf := make([]byte, protocolAckByteSize)
	if _, err := io.ReadFull(p.conn, buf); err != nil {
		return 0, err
	}
	return buf[protocolAckIndex], nil
}

func (p *ClientProtocol) readU16() (int, error) {
	buf := make([]byte, protocolWinnersCountByteSize)
	if _, err := io.ReadFull(p.conn, buf); err != nil {
		return 0, err
	}
	return int(binary.BigEndian.Uint16(buf)), nil
}

func (p *ClientProtocol) readU32() (uint32, error) {
	buf := make([]byte, protocolDNIByteSize)
	if _, err := io.ReadFull(p.conn, buf); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(buf), nil
}

func (p *ClientProtocol) Close() {
	if p.conn != nil {
		_ = p.conn.Close()
		p.conn = nil
	}
}
