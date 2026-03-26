package common

import "fmt"

const (
	protocolMessageTypeBatch uint8 = 1

	protocolBatchHeaderByteSize = 3 // type(1) + count(1) + agency(1)
	protocolMaxU8Value          = 255
)

type ProtocolBuilder struct{}

func NewProtocolBuilder() *ProtocolBuilder {
	return &ProtocolBuilder{}
}

func (b *ProtocolBuilder) BuildBatchHeader(batchCount, agencyID uint8) ([]byte, error) {
	if batchCount == 0 {
		return nil, fmt.Errorf("batch count must be greater than zero")
	}

	return []byte{protocolMessageTypeBatch, batchCount, agencyID}, nil
}

func (b *ProtocolBuilder) BuildBatchMessage(agencyID uint8, bets []Bet, payloadCapacity int) ([]byte, error) {
	count := len(bets)
	if count == 0 {
		return nil, fmt.Errorf("cannot build empty batch")
	}
	if count > protocolMaxU8Value {
		return nil, fmt.Errorf("batch size exceeds u8 range")
	}

	header, err := b.BuildBatchHeader(uint8(count), agencyID)
	if err != nil {
		return nil, err
	}

	payload := make([]byte, 0, payloadCapacity)
	payload = append(payload, header...)

	for _, bet := range bets {
		encoded, err := bet.Serialize()
		if err != nil {
			return nil, err
		}
		payload = append(payload, encoded...)
	}

	return payload, nil
}
