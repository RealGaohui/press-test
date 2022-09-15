package main

import (
	"os"
	Controller "press-test/controller"
)

func main() {
	if err := Controller.Prepare(); err != nil {
		os.Exit(1)
	}
}
