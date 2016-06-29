# go-stdlogger

The Go Standard Logging Interface. Plain and simple.

```go
// Logger describes a logger interface that is compatible with the standard
// log.Logger but also logrus and others. As not to limit which loggers can and
// can't be used with the API.
//
// This interface is from https://godoc.org/github.com/Sirupsen/logrus#StdLogger
type Logger interface {
    Print(...interface{})
    Printf(string, ...interface{})
    Println(...interface{})

    Fatal(...interface{})
    Fatalf(string, ...interface{})
    Fatalln(...interface{})

    Panic(...interface{})
    Panicf(string, ...interface{})
    Panicln(...interface{})
}
```
