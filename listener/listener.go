package listener

type Listener interface {
	Listen(Handler) error
}

type Handler interface {
	Handle() error
}
