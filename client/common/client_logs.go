package common

func logBetSendFailure(clientID string, err error) {
	log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | error: %v",
		clientID,
		err,
	)
}

func logAckReadFailure(clientID string, err error) {
	log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
		clientID,
		err,
	)
}

func logAckRejected(clientID string, ok bool) {
	log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | confirmation: %v",
		clientID,
		ok,
	)
}
