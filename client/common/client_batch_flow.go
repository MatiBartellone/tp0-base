package common

import (
	"fmt"
	"io"
	"strconv"
)

func (c *Client) parseAgencyID() (uint8, error) {
	agencyID, err := strconv.Atoi(c.config.ID)
	if err != nil {
		return 0, err
	}
	if agencyID < 0 || agencyID > maxU8Value {
		return 0, fmt.Errorf("invalid agency id: %d", agencyID)
	}

	return uint8(agencyID), nil
}

func (c *Client) sendBatch(batchPayload []byte) (bool, error, error) {
	if err := c.protocol.SendBatch(batchPayload); err != nil {
		return false, err, nil
	}

	ok, ackErr := c.protocol.ReadBatchAck()
	return ok, nil, ackErr
}

func (c *Client) sendBuiltBatch(builder *BatchBuilder) bool {
	payload, _, err := builder.Build()
	if err != nil {
		logBetSendFailure(c.config.ID, err)
		return false
	}

	ok, sendErr, ackErr := c.sendBatch(payload)
	if sendErr != nil {
		logBetSendFailure(c.config.ID, sendErr)
		return false
	}

	if !c.handleSendResult(ok, ackErr) {
		return false
	}

	builder.Reset()
	return true
}

func (c *Client) sendDatasetBatches(agencyID uint8) bool {
	file, err := openAgencyDataset(c.config.ID)
	if err != nil {
		logBetSendFailure(c.config.ID, err)
		return false
	}
	defer file.Close()

	reader := newDatasetReader(file)
	builder := NewBatchBuilder(agencyID, c.config.BatchMaxAmount)

	for {
		if c.isShuttingDown() {
			return false
		}

		record, readErr := reader.Read()
		if readErr == io.EOF {
			if !builder.IsEmpty() && !c.sendBuiltBatch(builder) {
				return false
			}
			return true
		}
		if readErr != nil {
			logBetSendFailure(c.config.ID, readErr)
			return false
		}

		bet, betErr := NewBetFromRecord(record)
		if betErr != nil {
			logBetSendFailure(c.config.ID, betErr)
			return false
		}

		if builder.Add(bet) {
			continue
		}

		if !c.sendBuiltBatch(builder) {
			return false
		}

		if !builder.Add(bet) {
			logBetSendFailure(c.config.ID, io.ErrShortBuffer)
			return false
		}
	}
}
