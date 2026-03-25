package common

import (
	"fmt"
	"strconv"
	"time"

	"github.com/op/go-logging"
)

const (
	firstNameIdx = 0
	lastNameIdx  = 1
	documentIdx  = 2
	birthDateIdx = 3
	numberIdx    = 4
	recordFieldCount = 5

	dateLayout = "2006-01-02"

	base10 = 10
	documentBitSize = 32
	numberBitSize   = 16

	birthYearByteSize  = 2
	birthMonthByteSize = 1
	birthDayByteSize   = 1
)

type Bet struct {
	FirstName string `json:"nombre"`
	LastName  string `json:"apellido"`
	Document  string `json:"documento"`
	BirthDate string `json:"nacimiento"`
	Number    string `json:"numero"`
}

func NewBet(firstName, lastName, document, birthDate, number string) (Bet, error) {
	return normalizeBet(firstName, lastName, document, birthDate, number)
}

func NewBetFromRecord(record []string) (Bet, error) {
	if len(record) < recordFieldCount {
		return Bet{}, fmt.Errorf("invalid record length: %d", len(record))
	}

	return normalizeBet(
		record[firstNameIdx],
		record[lastNameIdx],
		record[documentIdx],
		record[birthDateIdx],
		record[numberIdx],
	)
}

func normalizeBet(firstName, lastName, document, birthDate, number string) (Bet, error) {
	if _, err := time.Parse(dateLayout, birthDate); err != nil {
		return Bet{}, fmt.Errorf("invalid birth date %q", birthDate)
	}

	doc, err := strconv.ParseUint(document, base10, documentBitSize)
	if err != nil {
		return Bet{}, fmt.Errorf("invalid document %q", document)
	}

	parsedNumber, err := strconv.ParseUint(number, base10, numberBitSize)
	if err != nil {
		return Bet{}, fmt.Errorf("invalid number %q", number)
	}

	return Bet{
		FirstName: firstName,
		LastName:  lastName,
		Document:  strconv.FormatUint(doc, base10),
		BirthDate: birthDate,
		Number:    strconv.FormatUint(parsedNumber, base10),
	}, nil
}

func serializeBirthDate(value string, out []byte) ([]byte, error) {
	birthDate, err := time.Parse(dateLayout, value)
	if err != nil {
		return nil, err
	}

	out = serializeU16(uint16(birthDate.Year()), out)
	out = serializeU8(uint8(birthDate.Month()), out)
	out = serializeU8(uint8(birthDate.Day()), out)
	return out, nil
}

func (b Bet) Serialize() ([]byte, error) {
	doc, err := strconv.ParseUint(b.Document, base10, documentBitSize)
	if err != nil {
		return nil, err
	}

	number, err := strconv.ParseUint(b.Number, base10, numberBitSize)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 0, 96)
	buf = serializeString(b.FirstName, buf)
	buf = serializeString(b.LastName, buf)

	buf, err = serializeBirthDate(b.BirthDate, buf)
	if err != nil {
		return nil, err
	}

	buf = serializeU32(uint32(doc), buf)
	buf = serializeU16(uint16(number), buf)

	return buf, nil
}

func (b Bet) EncodedSize() int {
	return u16ByteSize + len(b.FirstName) +
		u16ByteSize + len(b.LastName) +
		birthYearByteSize + birthMonthByteSize + birthDayByteSize +
		u32ByteSize + u16ByteSize
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
