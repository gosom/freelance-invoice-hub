package main

import (
	"context"

	"invoicehub/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	router := http.NewRouter(nil, nil)
	params := http.ServerParams{
		Handler: router.Handler(),
	}

	srv := http.New(&params)

	err := srv.Start(ctx)
	if err != nil {
		panic(err)
	}
}
