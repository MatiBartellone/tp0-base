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
	if agencyID < 0 || agencyID > protocolMaxU8Value {
		return 0, fmt.Errorf("invalid agency id: %d", agencyID)
	}

	return uint8(agencyID), nil
}

func (c *Client) sendBatch(batchPayload []byte) error {
	if err := c.protocol.SendBatch(batchPayload); err != nil {
		return fmt.Errorf("send batch failure: %w", err)
	}

	ok, ackErr := c.protocol.ReadBatchAck()
	if err := c.handleSendResult(ok, ackErr); err != nil {
		return err
	}

	return nil
}

func (c *Client) sendBuiltBatch(builder *BatchBuilder) error {
	payload, _, err := builder.Build()
	if err != nil {
		return fmt.Errorf("build batch failure: %w", err)
	}

	if err := c.sendBatch(payload); err != nil {
		return err
	}

	builder.Reset()
	return nil
}

func (c *Client) sendDatasetBatches(agencyID uint8) error {
	file, err := openAgencyDataset(c.config.ID)
	if err != nil {
		return fmt.Errorf("open dataset failure: %w", err)
	}
	defer file.Close()

	reader := newDatasetReader(file)
	builder := NewBatchBuilder(agencyID, c.config.BatchMaxAmount)

	for {
		if c.isShuttingDown() {
			return nil
		}

		record, readErr := reader.Read()
		if readErr == io.EOF {
			if !builder.IsEmpty() {
				if err := c.sendBuiltBatch(builder); err != nil {
					return err
				}
			}
			return nil
		}
		if readErr != nil {
			return fmt.Errorf("read dataset failure: %w", readErr)
		}

		bet, betErr := NewBetFromRecord(record)
		if betErr != nil {
			return fmt.Errorf("record decode failure: %w", betErr)
		}

		if builder.Add(bet) {
			continue
		}

		if err := c.sendBuiltBatch(builder); err != nil {
			return err
		}

		if !builder.Add(bet) {
			return fmt.Errorf("batch append failure: %w", io.ErrShortBuffer)
		}
	}
}
