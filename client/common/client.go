package common

import (
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
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

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	c.registerSignalHandler()

	if err := c.createClientSocket(); err != nil {
		return
	}
	defer c.closeConn()

	agencyID, err := strconv.Atoi(c.config.ID)
	if err != nil {
		logBetSendFailure(c.config.ID, err)
		return
	}

	file, err := openAgencyDataset(c.config.ID)
	if err != nil {
		logBetSendFailure(c.config.ID, err)
		return
	}
	defer file.Close()

	reader := newDatasetReader(file)
	builder := NewBatchBuilder(uint8(agencyID), c.config.BatchMaxAmount)

	for {
		if c.isShuttingDown() {
			return
		}

		record, readErr := reader.Read()
		if readErr == io.EOF {
			if !builder.IsEmpty() && !c.sendBuiltBatch(builder) {
				return
			}
			break
		}
		if readErr != nil {
			logBetSendFailure(c.config.ID, readErr)
			return
		}

		bet, betErr := NewBetFromRecord(record)
		if betErr != nil {
			logBetSendFailure(c.config.ID, betErr)
			return
		}

		if builder.Add(bet) {
			continue
		}

		if !c.sendBuiltBatch(builder) {
			return
		}

		if !builder.Add(bet) {
			logBetSendFailure(c.config.ID, io.ErrShortBuffer)
			return
		}
	}

	if c.isShuttingDown() {
		return
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
