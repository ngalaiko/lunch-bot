package main

import (
	"context"
	"log"

	"lunch/pkg/migrate"
)

func main() {
	if err := migrate.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
