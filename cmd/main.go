package main

import (
	"context"
	"os"

	"invoicehub/http"
	"invoicehub/sqlite"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbpath := os.Getenv("FH_DB_PATH")
	if dbpath == "" {
		// defaults to a path valid for our Docker setup
		dbpath = "/app/data.sqlite3"
	}

	db, err := sqlite.SetupDB(dbpath)
	if err != nil {
		panic(err)
	}

	companyRepo := sqlite.NewCompanyRepository(db)

	router := http.NewRouter(companyRepo, nil)
	params := http.ServerParams{
		Handler: router.Handler(),
	}

	srv := http.New(&params)

	err = srv.Start(ctx)
	if err != nil {
		panic(err)
	}
}
