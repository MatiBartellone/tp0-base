package common

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// Client Entity that encapsulates how
type Client struct {
	config       ClientConfig
	protocol     *ClientProtocol
	shutdownChan chan struct{}
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config:       config,
		shutdownChan: make(chan struct{}),
	}
	return client
}

func (c *Client) Stop() {
	select {
	case <-c.shutdownChan:
		return
	default:
		close(c.shutdownChan)
		c.closeConn()
		log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
	}
}

func (c *Client) closeConn() {
	if c.protocol != nil {
		c.protocol.Close()
		c.protocol = nil
	}
}

// createClientSocket initializes client socket and protocol.
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	c.protocol = NewClientProtocol(conn)
	return nil
}

func (c *Client) registerSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Infof("action: shutdown | result: in_progress | client_id: %v", c.config.ID)
		c.Stop()
	}()
}

func (c *Client) isShuttingDown() bool {
	select {
	case <-c.shutdownChan:
		return true
	default:
		return false
	}
}

func (c *Client) initBet() (Bet, bool) {
	bet, err := NewBet(c.config)
	if err != nil {
		logBetSendFailure(c.config.ID, err)
		return Bet{}, false
	}

	return bet, true
}

func (c *Client) sendBet(bet Bet) (bool, error, error) {
	if err := c.createClientSocket(); err != nil {
		return false, err, nil
	}

	if err := c.protocol.Send(bet.Serialize()); err != nil {
		c.closeConn()
		return false, err, nil
	}

	ok, ackErr := c.protocol.ReadAck()
	c.closeConn()
	return ok, nil, ackErr
}

func (c *Client) handleSendResult(bet Bet, ok bool, ackErr error) bool {
	if ackErr != nil {
		if c.isShuttingDown() {
			return false
		}
		logAckReadFailure(c.config.ID, ackErr)
		return false
	}

	if !ok {
		logAckRejected(c.config.ID, ok)
		return false
	}

	bet.LogSentSuccess(log)
	return true
}

func (c *Client) waitNextIteration() bool {
	select {
	case <-time.After(c.config.LoopPeriod):
		return true
	case <-c.shutdownChan:
		return false
	}
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	c.registerSignalHandler()
	bet, ok := c.initBet()
	if !ok {
		return
	}

	for i := 0; i < c.config.LoopAmount; i++ {
		if c.isShuttingDown() {
			return
		}

		ackOK, sendErr, ackErr := c.sendBet(bet)
		if sendErr != nil {
			logBetSendFailure(c.config.ID, sendErr)
			return
		}

		if !c.handleSendResult(bet, ackOK, ackErr) {
			return
		}

		if !c.waitNextIteration() {
			return
		}
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
