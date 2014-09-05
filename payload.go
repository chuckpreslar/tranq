package tranq

import (
	"encoding/json"
)

// Payload represents a map of strings to interfaces.
type Payload map[string]interface{}

// Marshal calls the json packages `Marshal` method
// with the receiver Payload.
func (p Payload) Marshal() ([]byte, error) {
	return json.Marshal(p)
}
