package handler

import (
	"github.com/go-ozzo/ozzo-routing"
	"go-api/go-web-api/apiweb-api/api/service"
)

func NewVersionHandler() VersionHandler {
	return VersionHandler{}
}

type VersionHandler struct{}

/*----------------------------------------测试service-----------------------------------------------------*/

//版本
func (v VersionHandler) GetVersion(c *routing.Context) error {
	version := service.Version.GetVersion()
	return c.Write(version)
}
