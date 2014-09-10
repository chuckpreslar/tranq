package configurators

import (
	"github.com/chuckpreslar/tranq/serializers"
)

// Configurator interface provides the ability to
// create and customize types implementing the
// serializers.Serializer interface.
type Configurator interface {
	// NewSerializer returns an instance of the
	// serializers.Serializer interface.
	NewSerializer() serializers.Serializer
}
