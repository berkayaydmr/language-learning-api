package main

import (
	"context"

	"github.com/berkayaydmr/language-learning-api/pkg/storage"
)

func main() {
	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage := storage.New()

	err := storage.Open(context, "../server/words.db")
	if err != nil {
		panic(err)
	}

	err = storage.CreateTables(context)
	if err != nil {
		panic(err)
	}

	err = storage.SeedData(context)
	if err != nil {
		panic(err)
	}

	err = storage.Close()
	if err != nil {
		panic(err)
	}
}
