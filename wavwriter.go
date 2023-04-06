package main

import (
	"errors"
	"io"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type WAVWriter struct {
	io.Closer
	file       *os.File
	encoder    *wav.Encoder
	sampleRate int
	channels   int
}

func NewWAVWriter(path string, sampleRate int, channels int) (*WAVWriter, error) {
	f, err := os.Create(path) //nolint:gosec
	if err != nil {
		return nil, err
	}
	return &WAVWriter{
		file:       f,
		encoder:    wav.NewEncoder(f, sampleRate, 16, channels, 1),
		sampleRate: sampleRate,
		channels:   channels,
	}, nil
}

func (w *WAVWriter) Close() error {
	return errors.Join(w.encoder.Close(), w.file.Close())
}

func (w *WAVWriter) WritePCMSamples(pcm []int16) error {
	pcmbuf := audio.PCMBuffer{
		Format: &audio.Format{
			NumChannels: int(w.channels),
			SampleRate:  int(w.sampleRate),
		},
		I16:            pcm,
		DataType:       audio.DataTypeI16,
		SourceBitDepth: 2,
	}
	return w.encoder.Write(pcmbuf.AsIntBuffer())
}
