# slog -- Structured Logging with Context

[![GoDoc](https://godoc.org/github.com/spkg/slog/logfmt?status.svg)](https://godoc.org/github.com/spkg/slog)
[![license](http://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/spkg/slog/master/LICENSE.md)
[![Build Status](https://travis-ci.org/spkg/slog.svg?branch=master)](https://travis-ci.org/spkg/slog)
[![Coverage](http://gocover.io/_badge/github.com/spkg/slog)](http://gocover.io/github.com/spkg/slog)

## Another logging package? (Really?)

Yes. `slog` is another logging package. It's probably worth listing some of the more mature logging packages 
out there that may suit your purpose:

* [log: Go language standard library package](https://golang.org/pkg/log/)
* [glog: Leveled execution logs for Go](https://github.com/golang/glog)
* [logrus: Structured, pluggable logging for Go](https://github.com/Sirupsen/logrus)
* [loggo: Module level logging for Go](https://godoc.org/github.com/juju/loggo)

## Key-value pairs

Package `slog` does not encourage the use of `Printf`-like methods for formatting messages. Instead it encourages
the use of key-value pairs for logging properties associated with a log message. So, instead of of

```
 log.Printf("error: cannot open file %s: %s", filename, err.Error())
``` 

which would look like:

```
 error: cannot open file /etc/hosts: file not found
```

Logging key-value pairs looks like:

```
 slog.Error(ctx, "cannot open file",
     slog.WithValue("filename", filename),
     slog.WithError(err))
```

which results in a log message like:

```
 error: cannot open file filename=/etc/hosts error="file not found"
```

This idea has gained traction as a [best practice](http://dev.splunk.com/view/logging-best-practices/SP-CAAADP6),
and the result is both readable by humans, and can be [parsed quite easily](https://github.com/kr/logfmt).