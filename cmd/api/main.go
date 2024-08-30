package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type api struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":8000", "port to run the web server on")
	flag.Parse()

	api := &api{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	api.logger.Info("server running at port", "addr", *addr)

	err := http.ListenAndServe(*addr, api.routes())
	api.logger.Error(err.Error())
	os.Exit(1)
}
