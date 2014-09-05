package tranq

import (
	"encoding/json"
)

// Payload ...
type Payload map[string]interface{}

// Marshal ...
func (p Payload) Marshal() ([]byte, error) {
	return json.Marshal(p)
}
