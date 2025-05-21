package broker

type Broker interface {
	Connect() error
	Disconnect() error
	Publish(topic string, msg []byte) error
}
