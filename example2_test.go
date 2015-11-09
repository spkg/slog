package slog_test

import (
	"log"
	"net/http"

	"github.com/spkg/slog"
	"golang.org/x/net/context"
)

type ExternalHandler struct {
	// ... details for accessing external handler ...
}

func (h *ExternalHandler) Handle(msgs []*slog.Message) {
	// ... send to msgs to external logging handler ...
}

func ExampleAddHandler() {
	slog.AddHandler(&ExternalHandler{})
}

func ExampleNewWriter(ctx context.Context) {
	// Creates a HTTP server whose error log will write to the
	// default slog.Logger. Any panics that are recovered will
	// have the details logged via slog.
	httpServer := &http.Server{
		Addr:     ":8080",
		ErrorLog: log.New(slog.NewWriter(ctx), "http", 0),
	}

	slog.Info(ctx, "web server started", slog.WithValue("addr", httpServer.Addr))
	if err := httpServer.ListenAndServe(); err != nil {
		slog.Error(ctx, "web server failed", slog.WithError(err))
	}

	// 2009-11-10T12:34:56.789 info msg="web server started" addr=:8080
	// 2009-11-10T12:35:57:987 error msg="http: panic serving 123.1:2.3:36145 runtime error: invalid memory address or nil pointer dereference"
}

func ExampleOption(ctx context.Context, n1, n2 int) error {
	if err := doSomethingWith(n1, n2); err != nil {
		return slog.Error(ctx, "cannot doSomething",
			slog.WithValue("n1", n1),
			slog.WithValue("n2", n2),
			slog.WithError(err))
	}

	// .. more processing and then ...

	return nil
}

func ExampleWithValue(ctx context.Context, n1, n2 int) error {
	if err := doSomethingWith(n1, n2); err != nil {
		return slog.Error(ctx, "doSomethingWith failed",
			slog.WithValue("n1", n1),
			slog.WithValue("n2", n2))
	}

	// ... more processing and then ...

	return nil
}

func ExampleWithError(ctx context.Context) error {
	if err := doSomething(); err != nil {
		return slog.Error(ctx, "doSomething failed",
			slog.WithError(err))
	}

	// ... more processing and then ...

	return nil
}

func ExampleNewContext(ctx context.Context, n1, n2 int) error {
	ctx = slog.NewContext(ctx,
		slog.Property{"n1", n1},
		slog.Property{"n2", n2})

	if err := doSomethingWith(n1, n2); err != nil {
		return slog.Error(ctx, "doSomethingWith failed",
			slog.WithError(err))
	}

	slog.Debug(ctx, "did something with")

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
