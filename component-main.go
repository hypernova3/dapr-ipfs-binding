package main

import (
	components "github.com/dapr-sandbox/components-go-sdk"
	bindings "github.com/dapr-sandbox/components-go-sdk/bindings/v1"
	"github.com/dapr/kit/logger"

	"github.com/hypernova3/dapr-ipfs-binding/component"
)

var log = logger.NewLogger("ipfs-pluggable")

func main() {
	components.Register("ipfs-pluggable", components.WithOutputBinding(func() bindings.OutputBinding {
		return component.NewIPFSBinding(log)
	}))
	components.MustRun()
}
