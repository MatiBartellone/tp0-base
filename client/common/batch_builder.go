package common

import "fmt"

const (
	maxBatchBytes   = 8000
	headerBytes     = 3 // type(1) + count(1) + agency(1)
	typeBatch       = 1
	zeroBetsInBatch = 0
	minBatchBets    = 1
	maxU8Value      = 255
)

type BatchBuilder struct {
	agencyID     uint8
	maxBatchBets int
	currentSize  int
	currentBets  []Bet
}

func NewBatchBuilder(agencyID uint8, maxBatchBets int) *BatchBuilder {
	if maxBatchBets <= zeroBetsInBatch {
		maxBatchBets = minBatchBets
	}

	return &BatchBuilder{
		agencyID:     agencyID,
		maxBatchBets: maxBatchBets,
		currentSize:  headerBytes,
		currentBets:  make([]Bet, 0, maxBatchBets),
	}
}

func (b *BatchBuilder) Add(bet Bet) bool {
	if len(b.currentBets) >= b.maxBatchBets {
		return false
	}

	nextSize := b.currentSize + bet.EncodedSize()
	if nextSize > maxBatchBytes {
		return false
	}

	b.currentBets = append(b.currentBets, bet)
	b.currentSize = nextSize
	return true
}

func (b *BatchBuilder) IsEmpty() bool {
	return len(b.currentBets) == 0
}

func (b *BatchBuilder) Reset() {
	b.currentBets = b.currentBets[:0]
	b.currentSize = headerBytes
}

func (b *BatchBuilder) Build() ([]byte, int, error) {
	count := len(b.currentBets)
	if count == zeroBetsInBatch {
		return nil, 0, fmt.Errorf("cannot build empty batch")
	}
	if count > maxU8Value {
		return nil, 0, fmt.Errorf("batch size exceeds u8 range")
	}

	payload := make([]byte, 0, b.currentSize)
	payload = append(payload, byte(typeBatch))
	payload = append(payload, byte(count))
	payload = append(payload, b.agencyID)

	for _, bet := range b.currentBets {
		encoded, err := bet.Serialize()
		if err != nil {
			return nil, 0, err
		}
		payload = append(payload, encoded...)
	}

	return payload, count, nil
}
