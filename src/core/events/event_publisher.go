package events

// EventPublisher defines the contract for publishing domain events.
//
// Implement this interface in your WebSocket hub (external project) and
// inject it into any use-case that needs real-time broadcasting.
//
// Example:
//
//	type WsHub struct { ... }
//	func (h *WsHub) Publish(event interface{}) error { ... }
type EventPublisher interface {
	Publish(event interface{}) error
}

// NoopPublisher is the default no-operation publisher.
// It is used when no WebSocket hub has been registered yet so the
// application starts and runs without panicking.
type NoopPublisher struct{}

func (n *NoopPublisher) Publish(_ interface{}) error { return nil }
