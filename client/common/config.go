package common

import (
	"fmt"
	"os"
	"time"
)

type BetInput struct {
	FirstName string
	LastName  string
	Document  string
	BirthDate string
	Number    string
}

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	BetInput      BetInput
}

func requiredEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("missing required env var: %s", key)
	}
	return value, nil
}

func LoadBetInputFromEnv() (BetInput, error) {
	firstName, err := requiredEnv("NOMBRE")
	if err != nil {
		return BetInput{}, err
	}
	lastName, err := requiredEnv("APELLIDO")
	if err != nil {
		return BetInput{}, err
	}
	document, err := requiredEnv("DOCUMENTO")
	if err != nil {
		return BetInput{}, err
	}
	birthDate, err := requiredEnv("NACIMIENTO")
	if err != nil {
		return BetInput{}, err
	}
	number, err := requiredEnv("NUMERO")
	if err != nil {
		return BetInput{}, err
	}

	return BetInput{
		FirstName: firstName,
		LastName:  lastName,
		Document:  document,
		BirthDate: birthDate,
		Number:    number,
	}, nil
}
