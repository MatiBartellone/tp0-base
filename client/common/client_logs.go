package common

func logClientLoopFailure(clientID string, err error) {
	log.Errorf("action: loop_finished | result: fail | client_id: %v | error: %v",
		clientID,
		err,
	)
}
