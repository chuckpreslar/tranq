package serializers

// Serializer interface provides the ability
// to serialize go objects into the format
// needed to comply with standards set by JSON API.
type Serializer interface {
	// Accept provides the Serializer with the
	// go object for serializtion.
	Accept(i interface{}) (map[string]interface{}, error)
}
