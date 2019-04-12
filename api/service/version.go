package service

import "go-api/go-web-api/conf"

var Version VersionService

type VersionService struct {
}

func (v VersionService) GetVersion() (version string) {
	return conf.App.Conf.Version
}
