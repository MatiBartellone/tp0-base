package common

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	ackByteSize    = 1
	ackSuccessCode = 1
	ackIndex       = 0

	typeBatchMessage        = 1
	typeFinishMessage       = 2
	typeWinnersQueryMessage = 3

	queryStatusFailure = 0
	queryStatusSuccess = 1
	queryStatusPending = 2

	winnersCountByteSize = 2
	dniByteSize          = 4
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

func (p *ClientProtocol) SendFinish(agencyID uint8) error {
	payload := []byte{typeFinishMessage, agencyID}
	return p.Send(payload)
}

func (p *ClientProtocol) SendWinnersQuery(agencyID uint8) error {
	payload := []byte{typeWinnersQueryMessage, agencyID}
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

func (p *ClientProtocol) ReadFinishAck() (bool, error) {
	return p.ReadAck()
}

func (p *ClientProtocol) ReadWinnersResponse() (uint8, []uint32, error) {
	status := make([]byte, ackByteSize)
	if _, err := io.ReadFull(p.conn, status); err != nil {
		return queryStatusFailure, nil, err
	}

	countRaw := make([]byte, winnersCountByteSize)
	if _, err := io.ReadFull(p.conn, countRaw); err != nil {
		return queryStatusFailure, nil, err
	}

	count := int(binary.BigEndian.Uint16(countRaw))
	if count == 0 {
		return status[ackIndex], nil, nil
	}

	winners := make([]uint32, count)
	buf := make([]byte, dniByteSize)
	for i := 0; i < count; i++ {
		if _, err := io.ReadFull(p.conn, buf); err != nil {
			return queryStatusFailure, nil, err
		}
		winners[i] = binary.BigEndian.Uint32(buf)
	}

	if status[ackIndex] == queryStatusSuccess || status[ackIndex] == queryStatusPending {
		return status[ackIndex], winners, nil
	}

	return status[ackIndex], winners, fmt.Errorf("unexpected winners response status: %d", status[ackIndex])
}

func (p *ClientProtocol) Close() {
	if p.conn != nil {
		_ = p.conn.Close()
		p.conn = nil
	}
}
