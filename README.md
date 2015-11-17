# slog
##Structured, Leveled Logging with Context

[![GoDoc](https://godoc.org/github.com/spkg/slog/logfmt?status.svg)](https://godoc.org/github.com/spkg/slog)
[![license](http://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/spkg/slog/master/LICENSE.md)
[![Build Status (Linux)](https://travis-ci.org/spkg/slog.svg?branch=master)](https://travis-ci.org/spkg/slog)
[![Build status (Windows)](https://ci.appveyor.com/api/projects/status/febtndgk5w4uv90h?svg=true)](https://ci.appveyor.com/project/jjeffery/slog)
[![Coverage](http://gocover.io/_badge/github.com/spkg/slog)](http://gocover.io/github.com/spkg/slog)

## Another logging package? Really?

Yes. `slog` is another logging package. It's probably worth listing some of the more mature logging packages 
out there that may suit your purpose:

* [log: Go language standard library package](https://golang.org/pkg/log/)
* [glog: Leveled execution logs for Go](https://github.com/golang/glog)
* [logrus: Structured, pluggable logging for Go](https://github.com/Sirupsen/logrus)
* [loggo: Module level logging for Go](https://godoc.org/github.com/juju/loggo)

## Structured

Package `slog` does not provide the use of `Printf`-like methods for formatting messages. Instead it encourages
the use of key-value pairs for logging properties associated with a log message. So, instead of of

```Go
 log.Printf("[error] cannot open file %s: %s", filename, err.Error())
``` 

which would look like:

```
 [error] cannot open file /etc/hosts: file not found
```

Logging key-value pairs looks like:

```Go
 slog.Error(ctx, "cannot open file",
     slog.WithValue("filename", filename),
     slog.WithError(err))
```

which results in a log message like:

```
 error msg="cannot open file" filename=/etc/hosts error="file not found"
```

This idea has gained traction as a [best practice](http://dev.splunk.com/view/logging-best-practices/SP-CAAADP6),
and the result is both readable by humans, and can be [parsed quite easily](https://github.com/kr/logfmt).

## Leveled

Like many other logging packages, `slog` requires the calling program to assign a level to each message
logged. The log levels available in the `slog` package are:

* **Debug** for messages that are of interest to software developers when they are debugging the application.
A debug message might involve quite low-level information, such as entering and leaving a function. 
* **Info** for messages that indicate an event of interest, but is not an error condition. An example might be
logging an info message when a new user signs up for a service. This might be useful for counting, but it
is not a cause for concern for the dev-ops team.
* **Warn** for messages that indicate a condition that is not necessarily a cause for concern by the dev-ops
team in its own right, but could be a cause for concern if it were to happen on a regular basis. An example
of a warning message might be an attempt to login with an invalid username/password combination. Not a problem
if it happens occasionally, but the dev-ops team might be interested if it were happening on a very regular
basis from one particular IP address.
* **Error** for messages that indicate a condition that may require immediate attention from the dev-ops team.
Error messages indicate some sort of failure that the application program may not be able to recover from
without human intervention.

The guidelines for the levels above are quite general. There is room for interpretation, and the level chosen
for any particular message can depend on the application, and the agreed standards of the development and
operations teams for that application.

It might be worth noting that there is a growing opinion that fewer levels might be better than many
levels. Dave Cheney, for example, promotes an argument that 
[warning messages should be eliminated](http://dave.cheney.net/2015/11/05/lets-talk-about-logging).

## With Context

Package `slog` makes heavy use of the `golang.org/x/net/context` package. If your application does not
use this context package, then you will probably want to look at one of the other logging packages, as
`slog` will not deliver much benefit to you. If you are unfamiliar with the `golang.org/x/net/context` 
package there is an excellent article on the [Go Blog](https://blog.golang.org/context). 

The `golang.org/x/net/context` package is useful when writing servers that handle requests. Common
examples are HTTP servers, RPC servers and batch processors. As each request is processed, multiple
goroutines may be started to assist with the processing of the request. By convention each function 
involved in the processing of the request receives as its first parameter a `ctx` variable of type
`context.Context`. The context makes it easy to pass values associated with the request, and `slog`
makes use of this by adding log properties to the request.

```Go
func Login(ctx context.Context, username, password string) (*User, error) {
    // create a new context with log properties
    ctx = slog.NewContext(ctx,
        slog.Property{"operation", "Login"}, 
        slog.Property{"username", username})

    // ... pass request onto database access functions ...
    user, err := db.FindUserByUsername(username)
    if err != nil {
        // will log `error msg="cannot find user" operation="Login" username="fnurk"`
        return nil, slog.Error(ctx, "cannot find user", slog.WithError(err))
    }

    // ... more processing ...

    return user, nil
 }
```

In the above example, the `Login` function attaches some properties to the context, so if at a 
later time an error condition is logged, the properties in the context are logged with the message.

## When to log an error message

When using `slog` to log error messages, there are a few simple rules of thumb for 
logging error messages. 

* If a function with a context calls a function without a context, then log any error
and return the message logged as the error.

	```Go
	func FuncWithContext(ctx context.Context, int someArg) error {
	    // calling a function that does not accept a context, could
	    // be some external library
	    if err := DoThatOneThing(someArg); err != nil {
	        // log a message and return that message as the error
	        return slog.Error(ctx, "cannot do that one thing",
	            slog.WithError(err)) 
	    }
	
	    // ... do more processing ...
	    return nil
	}
	```

* If a function with a context calls another function with a context, then there is no need 
log to log an error if the only processing to be performed is to pass the error back to
the caller.

	```Go
	func FuncWithContext(ctx context.Context, int someArg) error {
	    // calling another function with context: that function
        // will log an error if it encounters it and return
	    if err := DoOneThingWithContext(ctx, someArg); err != nil {
	        // don't log a message: the DoOneThingWithContext function
            // has already logged and all we are doing is passing the
            // error back to our caller
	        return nil, err 
	    }
	
	    // ... do more processing ...
	    return nil
	}
	```  

* If a function with a context calls another function with a context and receives an
error response, then if there is some non-trivial error handling then there may be
scope for additional logging.

	```Go
	func FuncWithContext(ctx context.Context, int someArg) error {
	    // calling another function with context: that function
        // will log an error if it encounters it and return
	    if err := DoOneThingWithContext(ctx, someArg); err != nil {
            slog.Info(ctx, "cleaning up")
            DoSomeCleanup(ctx, someArg)
	        return nil, err 
	    }
	
	    // ... do more processing ...
	    return nil
	}
	```  
 

## Messages are errors

As seen in the above examples, the `slog.Error`, `slog.Warn`, `slog.Info` and `slog.Debug` functions
all return a non-nil `*slog.Message`. This non-nil pointer implements the `error` interface, and
can be returned as an error value.

## Messages can have a status code

In the common case of a HTTP server, it may be useful to pass back a suggested HTTP status code
when logging an error:

```Go
if user, err := FindUser(username); err != nil {
	return slog.Error(ctx, "cannot find user", slog.WithError(err))
} else if user == nil {
	// log message and include a hint at a suitable HTTP status code
	return slog.Warn(ctx, "user not found",
		slog.WithStatusCode(http.StatusNotFound))
}

// ... continue processing user ...
```

The HTTP middleware can then make use of the status code later if necessary

```Go
// statusCodeFromError chooses a HTTP status code based on an error.
func statusCodeFromError(err error) int {
	// default to internal error
	statusCode := http.StatusInternalServerError

	type statusCoder interface {
		StatusCode() int
	}

	if errWithStatusCode, ok := err.(statusCoder); ok {
		if sc := errWithStatusCode.StatusCode(); sc > 0 {
			statusCode = sc
		}
	}

	return statusCode
}
```

## Messages can have an error code

There are times when it may be useful to pass back a code to inform the requesting party that a
specific error condition has occurred.

```Go
// optimistic locking exception has occurred
return slog.Info(ctx, "optimistic locking error",
	slog.WithCode("OptimisticLockingError"))
```

TODO: we have played around with an `errors`-like package with functions `StatusCode(error) int` and 
`Code(error) string`, but haven't got around to publishing it yet. It keeps changing with every
project we use it on.
