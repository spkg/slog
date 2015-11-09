// Package slog provides structured, context-aware logging. It is intended
// for use by applications that make use of the golang.org/x/net/context
// package. (See https://blog.golang.org/context for more information).
//
// The idea is that the application can build up information on the
// current request, transaction, or operation and store it in the
// current context. When an event occurs that needs to be logged, the
// built-up context is included in the log message.
//  file, err := os.Open(filename)
//  if err != nil {
//      // message logged will include any logging context stored in ctx
//      slog.Error(ctx, "cannot open file",
//          slog.WithValue("filename", filename),
//          slog.WithError(err))
//  }
//
//  // Output: (assuming userip has been set in the context)
//  // 2009-11-10T12:34:56.789 error msg="cannot open file" filename=/etc/hosts error="file does not exist" userip=123.231.111.222
//
// Out of the box, this package logs to stdout in logfmt format. (https://brandur.org/logfmt).
// Other formats are planned (including pretty TTY output), and a handler mechanism exists to
// integrate with external logging providers.
//
// See the examples for more details. A more comprehensive guide is
// available at https://github.com/spkg/slog
package slog
