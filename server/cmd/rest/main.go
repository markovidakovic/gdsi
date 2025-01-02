package main

import (
	"fmt"
	"log"

	"github.com/markovidakovic/gdsi/server/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("cfg: %+v\n", cfg)
}
