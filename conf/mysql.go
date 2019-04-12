package conf

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/rs/zerolog"
)

func OpenDB(dsn string, logger zerolog.Logger) (*xorm.Engine, error) {
	e, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		return nil, err
	}
	e.SetMapper(core.GonicMapper{})

	e.SetLogger(newSimpleLogger(logger, core.LOG_DEBUG))
	e.ShowSQL(true)
	e.ShowExecTime(true)

	return e, nil
}

type simpleLogger struct {
	logger  zerolog.Logger
	level   core.LogLevel
	showSQL bool
}

func newSimpleLogger(logger zerolog.Logger, l core.LogLevel) *simpleLogger {
	return &simpleLogger{
		logger: logger,
		level:  l,
	}
}

func (s *simpleLogger) Error(v ...interface{}) {
	if s.level <= core.LOG_ERR {
		s.logger.Error().Msg(fmt.Sprint(s.format(v)...))
	}
	return
}

func (s *simpleLogger) Errorf(format string, v ...interface{}) {
	if s.level <= core.LOG_ERR {
		s.logger.Error().Msgf(format, s.format(v)...)
	}
	return
}

func (s *simpleLogger) Debug(v ...interface{}) {
	if s.level <= core.LOG_DEBUG {
		s.logger.Debug().Msg(fmt.Sprint(s.format(v)...))
	}
	return
}

func (s *simpleLogger) Debugf(format string, v ...interface{}) {
	if s.level <= core.LOG_DEBUG {
		s.logger.Debug().Msgf(format, s.format(v)...)
	}
	return
}

func (s *simpleLogger) Info(v ...interface{}) {
	if s.level <= core.LOG_INFO {
		s.logger.Info().Msg(fmt.Sprint(s.format(v)...))
	}
	return
}

func (s *simpleLogger) Infof(format string, v ...interface{}) {
	if s.level <= core.LOG_INFO {
		s.logger.Info().Msgf(format, s.format(v)...)
	}
	return
}

func (s *simpleLogger) Warn(v ...interface{}) {
	if s.level <= core.LOG_WARNING {
		s.logger.Warn().Msg(fmt.Sprint(s.format(v)...))
	}
	return
}

func (s *simpleLogger) Warnf(format string, v ...interface{}) {
	if s.level <= core.LOG_WARNING {
		s.logger.Warn().Msgf(format, s.format(v)...)
	}
	return
}

func (s *simpleLogger) Level() core.LogLevel {
	return s.level
}

func (s *simpleLogger) SetLevel(l core.LogLevel) {
	s.level = l
	return
}

func (s *simpleLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		s.showSQL = true
		return
	}
	s.showSQL = show[0]
}

func (s *simpleLogger) IsShowSQL() bool {
	return s.showSQL
}

func (s *simpleLogger) format(v []interface{}) []interface{} {
	tmpV := make([]interface{}, len(v))
	copy(tmpV, v)
	if len(tmpV) >= 2 {
		if slice, ok := tmpV[1].([]interface{}); ok {
			tmpSlice := make([]interface{}, len(slice))
			copy(tmpSlice, slice)
			for i, item := range tmpSlice {
				switch raw := item.(type) {
				case []byte:
					tmpSlice[i] = fmt.Sprintf("<%d bytes>", len(raw))
					tmpV[1] = tmpSlice
				}
			}
		}
	}
	return tmpV
}
