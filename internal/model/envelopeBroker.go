package model

import "encoding/json"

type EnvelopeBroker struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
