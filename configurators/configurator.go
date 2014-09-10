package configurators

import (
	"github.com/chuckpreslar/tranq/serializers"
)

// Configurator ...
type Configurator interface {
	// Accept ...
	NewSerializer() serializers.Serializer
}
