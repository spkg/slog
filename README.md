# slog
##Structured, Levelled Logging with Context

[![GoDoc](https://godoc.org/github.com/spkg/slog/logfmt?status.svg)](https://godoc.org/github.com/spkg/slog)
[![license](http://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/spkg/slog/master/LICENSE.md)
[![Build Status](https://travis-ci.org/spkg/slog.svg?branch=master)](https://travis-ci.org/spkg/slog)
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

```
 log.Printf("[error] cannot open file %s: %s", filename, err.Error())
``` 

which would look like:

```
 [error] cannot open file /etc/hosts: file not found
```

Logging key-value pairs looks like:

```
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

## Levelled

Like many other logging packages, `slog` requires the calling program to assign a level to each message
logged. The log levels available in the `slog` package are:

* **Debug** for messages that are of interest to software developers when they are debugging the application.
A debug message might involve quite low-level information, such as entering and leaving a function. 
* **Info** for messages that indicate an event of interest, but is not an error condition. An example might be
logging an info message when a new user signs up for a service. This might be useful for counting, but it
is not a cause for concern
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
[warning messages should not be eliminated](http://dave.cheney.net/2015/11/05/lets-talk-about-logging).

## With Context

Package `slog` makes heavy use of the `golang.org/x/net/context` package. If your application does not
use this context package, then you will probably want to look at one of the other logging packages, as
`slog` will not deliver much benefit to you.

The `golang.org/x/net/context` package is described well on the [Go Blog](https://blog.golang.org/context). To
quote that site, this package "makes it easy to pass request-scoped values, cancelation signals
and deadlines across API boundaries to all the goroutines involved in handling a request". 

*TODO: finish this section*

