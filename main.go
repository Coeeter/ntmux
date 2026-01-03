package main

import (
	"context"
	"log"

	"github.com/coeeter/ntmux/cmd"
)

func main() {
	if err := cmd.RootCmd.ExecuteContext(context.Background()); err != nil {
		log.Fatal(err)
	}
}
