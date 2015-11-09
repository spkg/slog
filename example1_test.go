package slog_test

import (
	"errors"

	log "github.com/spkg/slog"
	"golang.org/x/net/context"
)

func Example() {
	// everything logged needs a context
	ctx := context.Background()

	log.Info(ctx, "program started")
	if err := doFirstThing(ctx, 5, 4); err != nil {
		// ... error processing here ...
	}

	// YYYY-MM-DDTHH:MM:SS.FFFFFF info msg="program started"
	// YYYY-MM-DDTHH:MM:SS.FFFFFF error msg="cannot do third thing" error="error message goes here" a=5 b=4
}

// doFirstThing illustrates setting a new logging context.
// Any message logged with the new context will have values for
// a and b.
func doFirstThing(ctx context.Context, a, b int) error {
	// add some logging to the context
	ctx = log.NewContext(ctx,
		log.Property{"a", a},
		log.Property{"b", b})

	// ... perform some more processing and then

	return doSecondThing(ctx)
}

func doSecondThing(ctx context.Context) error {
	if err := doThirdThing(); err != nil {
		// log the message at the first point where the context is available
		return log.Error(ctx, "cannot do third thing", log.WithError(err))
	}
	return nil
}

func doThirdThing() error {
	// this function has no context so it does not log,
	// it just returns an error message
	return errors.New("error message goes here")
}
