package router

import (
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/rs/zerolog"
	"go-api/go-web-api/api/routings"
	"go-api/go-web-api/conf"
	"net/http"
)

func Register(logger zerolog.Logger) http.Handler {
	gRouter := conf.App.Router
	gRouter.NotFound(notFound)
	gRouter.Use(
		routingLogger(logger),
		errorHandler(logger),
		content.TypeNegotiator(content.JSON),
	)

	routings.Register()
	for _, route := range gRouter.Routes() {
		logger.Debug().Msgf("register route: \"%-6s -> %s\".", route.Method(), route.Path())
	}
	return gRouter
}
