package common

import (
	"fmt"
	"net"
	"time"
)

const winnersQueryRetryPeriod = 150 * time.Millisecond

func (c *Client) sendFinish(agencyID uint8) error {
	if err := c.protocol.SendFinish(agencyID); err != nil {
		return fmt.Errorf("send finish failure: %w", err)
	}

	ok, ackErr := c.protocol.ReadFinishAck()
	if err := c.handleSendResult(ok, ackErr); err != nil {
		return err
	}

	return nil
}

func (c *Client) requestWinners(agencyID uint8) (uint8, int, error) {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		return queryStatusFailure, 0, err
	}
	protocol := NewClientProtocol(conn)
	defer protocol.Close()

	if err := protocol.SendWinnersQuery(agencyID); err != nil {
		return queryStatusFailure, 0, err
	}

	status, winners, err := protocol.ReadWinnersResponse()
	if err != nil {
		return queryStatusFailure, 0, err
	}

	return status, len(winners), nil
}

func (c *Client) waitWinnersResult(agencyID uint8) error {
	for {
		if c.isShuttingDown() {
			return nil
		}

		status, winnersCount, err := c.requestWinners(agencyID)
		if err != nil {
			return fmt.Errorf("winners query failure: %w", err)
		}

		if status == queryStatusPending {
			time.Sleep(winnersQueryRetryPeriod)
			continue
		}

		if status != queryStatusSuccess {
			return fmt.Errorf("unexpected winners status: %d", status)
		}

		logWinnersQuerySuccess(winnersCount)
		return nil
	}
}
