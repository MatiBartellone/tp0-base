package common

const (
	maxBatchBytes   = 8000
	zeroBetsInBatch = 0
	minBatchBets    = 1
)

type BatchBuilder struct {
	agencyID        uint8
	maxBatchBets    int
	currentSize     int
	currentBets     []Bet
	protocolBuilder *ProtocolBuilder
}

func NewBatchBuilder(agencyID uint8, maxBatchBets int) *BatchBuilder {
	if maxBatchBets <= zeroBetsInBatch {
		maxBatchBets = minBatchBets
	}

	return &BatchBuilder{
		agencyID:        agencyID,
		maxBatchBets:    maxBatchBets,
		currentSize:     protocolBatchHeaderByteSize,
		currentBets:     make([]Bet, 0, maxBatchBets),
		protocolBuilder: NewProtocolBuilder(),
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
	b.currentSize = protocolBatchHeaderByteSize
}

func (b *BatchBuilder) Build() ([]byte, int, error) {
	count := len(b.currentBets)
	payload, err := b.protocolBuilder.BuildBatchMessage(b.agencyID, b.currentBets, b.currentSize)
	if err != nil {
		return nil, 0, err
	}

	return payload, count, nil
}
