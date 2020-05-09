# Media

Quickly serve your sharable public directory

## Installation
`go get github.com/supanadit/media`

## Quick Start

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/supanadit/media"
)

func main() {
    // Create instance of Gin Engine
	g := gin.Default()
    // Serve and create shared directory
	_ = media.Gin(g).SetDestination("./upload").Create()
    // Serve gin at port :8080
	_ = g.Run(":8080")
}
```

## Note
Sorry, currently its only support [**Gin**](https://github.com/gin-gonic/gin) web framework, soon it will support [**Echo**](https://github.com/labstack/echo) as well