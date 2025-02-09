package main

import (
	"github.com/ramk42/mini-url/internal/infra/transport/api"
	"github.com/ramk42/mini-url/pkg/logger"
)

const ApplicationName = "mini-url"

func main() {
	logger.Init(ApplicationName)
	api.ListenAndServe()
}
