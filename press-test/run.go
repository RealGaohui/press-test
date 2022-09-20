package main

import (
	"context"
	"os"
	cfg "press-test/config"
	Controller "press-test/controller"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	tenantContext := context.WithValue(ctx, "tenant", cfg.Client)
	if err := Controller.Prepare(tenantContext); err != nil {
		os.Exit(1)
	}
}
