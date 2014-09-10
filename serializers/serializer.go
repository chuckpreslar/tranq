package serializers

// Serializer ...
type Serializer interface {
	// Accept ...
	Accept(i interface{}) (map[string]interface{}, error)
}
