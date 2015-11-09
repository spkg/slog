package slog_test

import (
	log "github.com/spkg/slog"
	"golang.org/x/net/context"
)

func ExampleOption(ctx context.Context, n1, n2 int) error {
	if err := doSomethingWith(n1, n2); err != nil {
		return log.Error(ctx, "cannot doSomething",
			log.WithValue("n1", n1),
			log.WithValue("n2", n2),
			log.WithError(err))
	}

	// .. more processing and then ...

	return nil
}

func ExampleWithValue(ctx context.Context, n1, n2 int) error {
	if err := doSomethingWith(n1, n2); err != nil {
		return log.Error(ctx, "doSomethingWith failed",
			log.WithValue("n1", n1),
			log.WithValue("n2", n2))
	}

	// ... more processing and then ...

	return nil
}

func ExampleWithError(ctx context.Context) error {
	if err := doSomething(); err != nil {
		return log.Error(ctx, "doSomething failed",
			log.WithError(err))
	}

	// ... more processing and then ...

	return nil
}

func ExampleNewContext(ctx context.Context, n1, n2 int) error {
	ctx = log.NewContext(ctx,
		log.Property{"n1", n1},
		log.Property{"n2", n2})

	if err := doSomethingWith(n1, n2); err != nil {
		return log.Error(ctx, "doSomethingWith failed",
			log.WithError(err))
	}

	log.Debug(ctx, "did something with")

	// ... more processing and then ...

	return nil
}

func doSomethingWith(n1 int, n2 int) error {
	return nil
}

func doSomething() error {
	return nil
}

func doAnotherThing() error {
	return nil
}

func doOneMoreThing() error {
	return nil
}
