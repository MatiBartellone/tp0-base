package common

import (
	"encoding/csv"
	"fmt"
	"os"
)

func openAgencyDataset(clientID string) (*os.File, error) {
	path := fmt.Sprintf(".data/agency-%s.csv", clientID)
	return os.Open(path)
}

func newDatasetReader(file *os.File) *csv.Reader {
	return csv.NewReader(file)
}
