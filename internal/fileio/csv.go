package fileio

import (
	"encoding/csv"
	"fmt"
	"os"
	"xvista/omise-challenge/internal/cipher"
)

func OpenEncryptedCSV(path string) (*csv.Reader, *os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}

	r, err := cipher.NewRot128Reader(file)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Rot128Reader: %w", err)
	}

	csvReader := csv.NewReader(r)
	_, err = csvReader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	return csvReader, file, nil
}
