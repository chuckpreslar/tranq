package tranq

import "github.com/chuckpreslar/tranq/configurators"

// Tranq ...
type Tranq struct {
	configurators.Configurator
}

// Serialize ...
func (t *Tranq) Serialize(i interface{}) (map[string]interface{}, error) {
	return t.NewSerializer().Accept(i)
}

// New ...
func New(c configurators.Configurator) *Tranq {
	var t = new(Tranq)
	t.Configurator = c
	return t
}
