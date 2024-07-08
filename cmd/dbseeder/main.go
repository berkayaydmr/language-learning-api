package main

import (
	"context"

	"github.com/berkayaydmr/language-learning-api/common"
	"github.com/berkayaydmr/language-learning-api/pkg/storage"
)

func main() {
	// storage olusturulacak
	// storage acilacak
	// storage tablolari olusturulacak
	// storage'a veri eklenecek
	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage := storage.New()

	err := storage.Open(context, common.DSN)
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
