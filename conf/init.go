package conf

import (
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-xorm/xorm"
	"github.com/rs/zerolog"
)

type Apps struct {
	Conf   Config
	DB     *xorm.Engine
	Logger zerolog.Logger
	Redis  *RedisClient
	Router *routing.Router
}

var (
	App Apps
)

//初始化服务
func Init() {
	/*-----------------------------------------读取环境配置文件-----------------------------------------------*/
	var err error
	App.Conf, err = NewConfig()
	if err != nil {
		panic(err)
	}
	/*-----------------------------------------配置日志文件-----------------------------------------------*/
	{
		App.Logger = NewLogger(App.Conf.Debug)
		App.Logger.Info().Interface("Config", App.Conf).Msg("Logger ready.")
	}
	/*-----------------------------------------配置MySQL数据库-----------------------------------------------*/
	{
		App.DB, err = OpenDB(App.Conf.Mysql, App.Logger)
		if err != nil {
			panic(err)
		}
		err = App.DB.Sync2()
		if err != nil {
			panic(err)
		}
		App.Logger.Info().Msg("DB ready.")
	}
	/*-----------------------------------------Redis数据库-----------------------------------------------*/
	{
		App.Redis, err = OpenRedis(App.Conf.Redis, 10, App.Logger)
		if err != nil {
			panic(err)
		}
		App.Logger.Info().Msg("redis ready.")
	}
	/*-----------------------------------------初始化路由-----------------------------------------------*/
	App.Router = routing.New()
}
