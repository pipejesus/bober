package main

import (
	"errors"
	"os"
	"time"

	"github.com/go-audio/wav"
)

func getWavLength(filePath string) (time.Duration, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return time.Duration(0), err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	decoder := wav.NewDecoder(file)

	if !decoder.IsValidFile() {
		return time.Duration(0), errors.New("invalid WAV file")
	}

	dur, err := decoder.Duration()
	if err != nil {
		return time.Duration(0), err
	}

	return dur, nil
}
