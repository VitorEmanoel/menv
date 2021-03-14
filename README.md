# Manager Environments  (menv)

Use this lib for load and validate environment variables.

- Support .env (or other name)

## Install

```bash
go get -u github.com/VitorEmanoel/menv
```

## Example

```golang
package main

import (
    "github.com/VitorEmanoel/menv"
    "log"
)

type Variables struct {
    Token   string  `env:"TOKEN" required:"true" default:"DEFAULT_TOKEN"`
}

func main() {
    var variables = Variables{}
    err := menv.LoadEnvironment(&variables)
    if err != nil {
        log.Fatal(err)
    }
    log.Println(variables.Token)
}
```
