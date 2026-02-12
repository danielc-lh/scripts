package capoeira

// Transport provides methods to send and receive messages between locations.
type Transport interface {
	Send(from, to string, data interface{})
	Receive(from, at string) interface{}
	Locations() []string
}
