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

  rd.WriteRequest([]string{"GET", "KEY"})
}
```

# Documentation

http://godoc.org/github.com/olahol/redisio
