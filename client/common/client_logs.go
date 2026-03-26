package common

func logClientLoopFailure(clientID string, err error) {
	log.Errorf("action: loop_finished | result: fail | client_id: %v | error: %v",
		clientID,
		err,
	)
}

func logClientLoopSuccess(clientID string) {
	log.Infof("action: loop_finished | result: success | client_id: %v", clientID)
}

func logClientShutdownInProgress(clientID string) {
	log.Infof("action: shutdown | result: in_progress | client_id: %v", clientID)
}

func logClientShutdownSuccess(clientID string) {
	log.Infof("action: shutdown | result: success | client_id: %v", clientID)
}

func logWinnersQuerySuccess(winnersCount int) {
	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", winnersCount)
}
