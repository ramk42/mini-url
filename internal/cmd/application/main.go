package main

import (
	"github.com/ramk42/mini-url/internal/infra/logging"
	"github.com/ramk42/mini-url/internal/infra/transport/api"
)

func main() {
	logging.InitLogger()
	api.ListenAndServe()
}
