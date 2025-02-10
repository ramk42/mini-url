// Package application provides the entry point to the application.
package application

import (
	"github.com/ramk42/mini-url/internal/infra/transport/api"
	"github.com/ramk42/mini-url/pkg/logger"
)

const ApplicationName = "mini-url"

func Run() {
	logger.Init(ApplicationName)
	api.ListenAndServe()
}
