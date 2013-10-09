[![Build Status](https://travis-ci.org/olahol/redisio.png)](https://travis-ci.org/olahol/redisio)

# redisio

Painlessly read and write the Redis protocol in go.

# Install

    $ go get github.com/olahol/redisio

# Example

Write a Redis request to stdout.

```go
package main

import (
  "os"
  "github.com/olahol/redisio"
)

func main() {
  rd := redisio.NewWriter(os.Stdout)

  rd.WriteRequest([]string{"SET", "KEY", "VALUE"})

  rd.Flush()
}
```

# Documentation

http://godoc.org/github.com/olahol/redisio
