package utils

import (
	"encoding/json"
	"go-rebuild/internal/model"
)

func BuildPacket(eventType string, m any) ([]byte, error) {
	body, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	packet := model.Envelope{
		Type:    eventType,
		Payload: body,
	}

	packetByte, err := json.Marshal(packet)
	if err != nil {
		return nil, err
	}

	return packetByte, nil
}