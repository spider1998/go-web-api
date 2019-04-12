package engine

import (
	"context"
	"go-api/go-web-api/conf"
	"go-api/go-web-api/router"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	std *Engine
)

type Engine struct {
	server *http.Server
	close  chan struct{}
	wg     sync.WaitGroup
}

func Get() *Engine {
	return std
}

func NewStdInstance() *Engine {
	conf.Init()
	std = new(Engine)
	std.close = make(chan struct{})
	std.server = &http.Server{Addr: conf.App.Conf.HTTPAddr}
	return std
}

func (e *Engine) Run() {
	go e.registerSignal()

	e.wg.Add(1)
	go e.serveHTTP()
	e.wg.Wait()
}

func (e *Engine) serveHTTP() {
	defer e.wg.Done()

	e.server.Handler = router.Register(
		conf.App.Logger,
	)

	conf.App.Logger.Info().Str("addr", conf.App.Conf.HTTPAddr).Msg("listen and serve http service.")

	err := e.server.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			conf.App.Logger.Error().Err(err).Msg("an error was returned while listen and serve engine.")
			return
		}
	}
	conf.App.Logger.Info().Msg("engine shutdown successfully.")
}

func (e *Engine) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	return e.server.Shutdown(ctx)
}

func (e *Engine) registerSignal() {
	conf.App.Logger.Info().Msg("register signal handler.")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case sig := <-ch:
		signal.Ignore(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)
		conf.App.Logger.Info().Str("signal", sig.String()).Msg("received signal, try to shutdown engine.")
		close(ch)
		close(e.close)
		err := e.shutdown()
		if err != nil {
			conf.App.Logger.Error().Err(err).Msg("fail to shutdown engine.")
		}
	}
}
