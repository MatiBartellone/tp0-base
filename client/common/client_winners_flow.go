package common

import (
	"fmt"
	"net"
	"time"
)

const winnersQueryRetryPeriod = 150 * time.Millisecond

func (c *Client) sendFinish(agencyID uint8) bool {
	if err := c.protocol.SendFinish(agencyID); err != nil {
		logBetSendFailure(c.config.ID, err)
		return false
	}

	ok, ackErr := c.protocol.ReadFinishAck()
	if !c.handleSendResult(ok, ackErr) {
		return false
	}

	return true
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

func (c *Client) waitWinnersResult(agencyID uint8) bool {
	for {
		if c.isShuttingDown() {
			return false
		}

		status, winnersCount, err := c.requestWinners(agencyID)
		if err != nil {
			logBetSendFailure(c.config.ID, err)
			return false
		}

		if status == queryStatusPending {
			time.Sleep(winnersQueryRetryPeriod)
			continue
		}

		if status != queryStatusSuccess {
			logBetSendFailure(c.config.ID, fmt.Errorf("unexpected winners status: %d", status))
			return false
		}

		log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", winnersCount)
		return true
	}
}
