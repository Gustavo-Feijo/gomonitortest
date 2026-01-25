package main

import (
	"context"
	"log"
	"os"
)

func main() {
	if err := Run(context.Background()); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
