package pkg

import (
	"go.uber.org/zap"
	"strings"
	"time"
)

const queryKey = "__q__"
const txnKey = "__txn__"

type logInfo struct {
	name      string
	startTime int64
	query     string
	logger    *zap.Logger
}

func (l logInfo) dump(args ...interface{}) {
	query := regexToCleanQueryToDump.ReplaceAllString(l.query, "")
	query = strings.ReplaceAll(query, "\n", " ")
	l.logger.Info(l.name, zap.Int64("time", time.Now().UnixMilli()-l.startTime), zap.String("query", query), zap.Any("args", args))
}

func (l logInfo) done(args ...interface{}) {
	l.dump(args...)
}

func newLogInf(name string, query string, logger *zap.Logger) logInfo {
	return logInfo{startTime: time.Now().UnixMilli(), name: name, query: query, logger: logger}
}
