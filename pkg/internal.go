package pkg

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"go.uber.org/zap"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var regexToCleanQueryToDump = regexp.MustCompile(`^[^\n]+\n`)
var regexFileQueryName = regexp.MustCompile(`^--\s*name:\s*(\S+)\s*:.*\n`)
var startMetricDumpSyncOnce = sync.Once{}

const queryKey = "__q__"
const txnKey = "__txn__"

var histogram = map[string]metrics.Histogram{}
var mutex = sync.Mutex{}

type logInfo struct {
	name       string
	startTime  int64
	timeTaken  int64
	query      string
	cleanQuery string
	logger     *zap.Logger
	hist       metrics.Histogram

	enableSqlQueryLogging       bool
	enableSqlQueryMetricLogging bool
}

func (l logInfo) dump(args ...interface{}) {
	if l.enableSqlQueryLogging {
		l.logger.Info(l.name, zap.Int64("time", l.timeTaken), zap.String("query", l.cleanQuery), zap.Any("args", args))
	}
}

func (l logInfo) done(args ...interface{}) {
	l.timeTaken = time.Now().UnixMilli() - l.startTime

	if l.enableSqlQueryMetricLogging {
		mutex.Lock()
		defer mutex.Unlock()

		q := l.cleanQuery
		if h, ok := histogram[q].(metrics.Histogram); !ok {
			s := metrics.NewExpDecaySample(1028, 0.015)
			hist := metrics.NewHistogram(s)
			if err := metrics.Register(q, hist); err != nil {
				fmt.Println(err)
			} else {
				l.hist = hist
				histogram[q] = hist
			}
		} else {
			l.hist = h
		}
		l.hist.Update(l.timeTaken)
	}
	l.dump(args...)
}

func newLogInf(name string, query string, logger *zap.Logger, enableSqlQueryLogging bool, enableSqlQueryMetricLogging bool) logInfo {

	// Start metric dumping - start only once
	startMetricDumpSyncOnce.Do(func() {
		go metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
	})

	//  re := regexp.MustCompile(`^--\s*name:\s*(\S+)\s*:.*\n`)
	//    match := re.FindStringSubmatch(input)
	//
	//    if len(match) > 1 {
	//        name := match[1]
	//        fmt.Println(name)
	cleanQuery := regexToCleanQueryToDump.ReplaceAllString(query, "")
	cleanQuery = strings.ReplaceAll(query, "\n", " ")
	return logInfo{startTime: time.Now().UnixMilli(), name: name, query: query, cleanQuery: cleanQuery, logger: logger,
		enableSqlQueryLogging: enableSqlQueryLogging, enableSqlQueryMetricLogging: enableSqlQueryMetricLogging}
}
