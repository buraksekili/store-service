package amqp

// Event is an interface for the messaging events.
type Event interface {
	Name() string
}

// AddUserEvent represents an event which RabbitMQ listens
// for newly created users.
type AddUserEvent struct {
	ID       string `bson:"_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Name returns a name of the event.
func (ue *AddUserEvent) Name() string {
	return "add_user"
}
