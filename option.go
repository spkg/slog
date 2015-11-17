package slog

// An Option is a function option that can be applied when logging a message.
// See the example for how they are used. Options is based on Dave Cheney's article
// "Functional options for friendly APIs" (http://goo.gl/l2KaW3)
// that can be applied to a Message.
type Option func(*Message)

// WithError sets the error associated with the log message.
func WithError(err error) Option {
	return func(m *Message) {
		m.Err = err
	}
}

// WithValue sets a property with a name and a value.
func WithValue(name string, value interface{}) Option {
	return func(m *Message) {
		m.Properties = append(m.Properties, Property{name, value})
	}
}

// WithCode associates an arbitrary code with the Message that is logged.
// This is useful for associating a pre-agreed error code with the message.
func WithCode(code string) Option {
	return func(m *Message) {
		m.code = code
	}
}

// WithStatusCode sets a status code associated with the log message. This is
// useful for associating a HTTP status code, but the status can be any
// number that makes sense for the application.
func WithStatusCode(status int) Option {
	return func(m *Message) {
		m.status = status
	}
}
