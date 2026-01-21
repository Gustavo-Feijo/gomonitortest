package main

import (
	"log"
	"os"
)

func main() {
	if err := Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
