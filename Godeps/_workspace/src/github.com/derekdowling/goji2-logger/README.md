# goji2-logger

Logging Middleware For [Goji 2](http://goji.io)

```go

import github.com/derekdowling/goji2-logger

// works with Logrus, stdlogger, most others as well
logger := New(os.Stderr, "", LstdFlags)
gojilogger.SetLogger(logger)

mux := goji.NewMux()
mux.UseC(gojilogger.Middleware)
```

Output looks like:

```bash
2015/12/16 16:55:21 Serving PATCH "/api/foos/1" from 127.0.0.1:57639
2015/12/16 16:55:21 Returning HTTP 200 after 53.044µs
```
