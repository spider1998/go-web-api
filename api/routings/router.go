package routings

import (
	"go-api/go-web-api/apiweb-api/api/handler"
	"go-api/go-web-api/conf"
)

func Register() {
	api := conf.App.Router.Group("")

	/*router.Use(......)*/
	{
		versionHandler := handler.NewVersionHandler()
		api.Get("/version", versionHandler.GetVersion)
	}
}
