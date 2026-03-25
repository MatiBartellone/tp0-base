package common

import (
	"encoding/binary"

	"github.com/op/go-logging"
)

type Bet struct {
	Agency    string `json:"agency"`
	FirstName string `json:"nombre"`
	LastName  string `json:"apellido"`
	Document  string `json:"documento"`
	BirthDate string `json:"nacimiento"`
	Number    string `json:"numero"`
}

func NewBet(config ClientConfig) (Bet, error) {
	return Bet{
		Agency:    config.ID,
		FirstName: config.BetInput.FirstName,
		LastName:  config.BetInput.LastName,
		Document:  config.BetInput.Document,
		BirthDate: config.BetInput.BirthDate,
		Number:    config.BetInput.Number,
	}, nil
}

func serializeString(value string, out []byte) []byte {
	b := []byte(value)
	lenBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBuf, uint16(len(b)))
	out = append(out, lenBuf...)
	out = append(out, b...)
	return out
}

func (b Bet) Serialize() []byte {
	buf := make([]byte, 0, 128)
	buf = serializeString(b.Agency, buf)
	buf = serializeString(b.FirstName, buf)
	buf = serializeString(b.LastName, buf)
	buf = serializeString(b.Document, buf)
	buf = serializeString(b.BirthDate, buf)
	buf = serializeString(b.Number, buf)
	return buf
}

func (b Bet) SentLogFields() (string, string) {
	return b.Document, b.Number
}

func (b Bet) LogSentSuccess(logger *logging.Logger) {
	dni, numero := b.SentLogFields()
	logger.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
		dni,
		numero,
	)
}
