package tranq

import "github.com/chuckpreslar/tranq/configurators"

// Tranq stores an instnace of the configurators.Configurator
// interface for creating and configuring serialization.Serializer
// instances.
type Tranq struct {
	configurators.Configurator
}

// Serialize uses the embedded configurators.Configurator
// instance to create a new serialization.Serializer
// instance and start serialization.
func (t *Tranq) Serialize(i interface{}) (map[string]interface{}, error) {
	return t.NewSerializer().Accept(i)
}

// New returns a new instance of the Tranq type.
func New(c configurators.Configurator) *Tranq {
	var t = new(Tranq)
	t.Configurator = c
	return t
}
