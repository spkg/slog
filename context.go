package slog

import "golang.org/x/net/context"

// logData is associated with a context and contains log
// values for the context.
type logData struct {
	Key   string
	Value interface{}
	Prev  *logData
}

type contextKey int

const (
	keyLogData contextKey = iota
)

// NewContext returns a new context that has one or more properties associated with it
// for logging purposes. The properties will be included in any log message logged with
// this context.
func NewContext(ctx context.Context, properties ...Property) context.Context {
	if len(properties) == 0 {
		return ctx
	}

	// extract log data from the context
	data, _ := ctx.Value(keyLogData).(*logData)

	// iterate backwards through the supplied properties:
	// this way they will be logged in the order that
	// they were supplied to this function
	for i := len(properties) - 1; i >= 0; i-- {
		p := properties[i]
		data = &logData{
			Key:   p.Key,
			Value: p.Value,
			Prev:  data,
		}
	}
	return context.WithValue(ctx, keyLogData, data)
}

func fromContext(ctx context.Context) *logData {
	if ctx == nil {
		return nil
	}
	data, ok := ctx.Value(keyLogData).(*logData)
	if !ok {
		return nil
	}
	return data
}
