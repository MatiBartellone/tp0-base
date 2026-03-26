package common

import (
	"fmt"
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

func (c *Client) handleSendResult(ok bool, ackErr error) error {
	if ackErr != nil {
		if c.isShuttingDown() {
			return nil
		}
		return fmt.Errorf("ack read failure: %w", ackErr)
	}

	if !ok {
		return fmt.Errorf("ack rejected by server")
	}
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	var loopErr error
	defer func() {
		if c.isShuttingDown() {
			return
		}

		if loopErr != nil {
			logClientLoopFailure(c.config.ID, loopErr)
			return
		}

		log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
	}()

	c.registerSignalHandler()

	if err := c.createClientSocket(); err != nil {
		loopErr = fmt.Errorf("connect failure: %w", err)
		return
	}
	defer c.closeConn()

	agencyID, err := c.parseAgencyID()
	if err != nil {
		loopErr = fmt.Errorf("agency parse failure: %w", err)
		return
	}

	if err := c.sendDatasetBatches(agencyID); err != nil {
		loopErr = err
		return
	}

	if c.isShuttingDown() {
		return
	}

	if err := c.sendFinish(agencyID); err != nil {
		loopErr = err
		return
	}

	c.closeConn()
	if err := c.waitWinnersResult(agencyID); err != nil {
		loopErr = err
		return
	}
}
