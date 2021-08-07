package measure

import (
	"time"

	"github.com/hirosuzuki/go-sql-logger/logger"
)

type MeasureRec struct {
	Query     string
	StartTime int64
}

func Start(query string) *MeasureRec {
	m := MeasureRec{
		Query:     query,
		StartTime: time.Now().UnixNano(),
	}
	return &m
}

func (m *MeasureRec) Stop() {
	timeDelta := time.Now().UnixNano() - m.StartTime
	logger.Log(m.StartTime, timeDelta, "", m.Query)
}
