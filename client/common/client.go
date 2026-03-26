package common

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

const signalChannelBufferSize = 1

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
	sigChan := make(chan os.Signal, signalChannelBufferSize)
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

func (c *Client) handleSendResult(ok bool, ackErr error) bool {
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
	return true
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	c.registerSignalHandler()

	if err := c.createClientSocket(); err != nil {
		return
	}
	defer c.closeConn()

	agencyID, err := c.parseAgencyID()
	if err != nil {
		logBetSendFailure(c.config.ID, err)
		return
	}

	if !c.sendDatasetBatches(agencyID) {
		return
	}

	if c.isShuttingDown() {
		return
	}

	if !c.sendFinish(agencyID) {
		return
	}

	c.closeConn()
	if !c.waitWinnersResult(agencyID) {
		return
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
