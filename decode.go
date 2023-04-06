package main

import (
	"fmt"

	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	libopus "gopkg.in/hraban/opus.v2"
)

const outputBufferMilliseconds = 120

type OpusDecoder struct {
	decoder      *libopus.Decoder
	channels     int
	outputBuffer []int16
}

func NewOpusDecoder(sampleRate int, channels int) (*OpusDecoder, error) {
	decoder, err := libopus.NewDecoder(sampleRate, channels)
	if err != nil {
		return nil, err
	}
	return &OpusDecoder{
		decoder:      decoder,
		channels:     channels,
		outputBuffer: make([]int16, channels*sampleRate*outputBufferMilliseconds/1000),
	}, nil
}

func (d *OpusDecoder) Decode(pkt *rtp.Packet) ([]int16, error) {
	if len(pkt.Payload) == 0 {
		if pkt.Padding {
			// Skip padding-only packet.
			return nil, nil
		} else {
			return nil, fmt.Errorf("0 length non-padding packet")
		}
	}
	var opusPkt codecs.OpusPacket
	if _, err := opusPkt.Unmarshal(pkt.Payload); err != nil {
		return nil, err
	}
	samplesWritten, err := d.decoder.Decode(opusPkt.Payload, d.outputBuffer)
	if err != nil {
		return nil, err
	}

	return d.outputBuffer[:d.channels*samplesWritten], nil
}
