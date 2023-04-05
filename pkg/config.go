package pkg

import "github.com/devlibx/gox-base/util"
import _ "github.com/go-sql-driver/mysql"

type MySQLConfig struct {
	ServerName   string `json:"server_name" yaml:"server_name"`
	Host         string `json:"host" yaml:"host"`
	Port         int    `json:"port" yaml:"port"`
	User         string `json:"user" yaml:"user"`
	Password     string `json:"password" yaml:"password"`
	Db           string `json:"db" yaml:"db"`
	TablePrefix  string
	TablePostfix string

	EnableSqlQueryLogging bool `json:"enable_sql_query_logging" yaml:"enable_sql_query_logging"`

	EnableSqlQueryMetricLogging bool `json:"enable_sql_query_metric_logging" yaml:"enable_sql_query_metric_logging"`
	MetricDumpIntervalSec       int  `json:"metric_dump_interval_sec" yaml:"metric_dump_interval_sec"`
	MetricResetAfterEveryNSec   int  `json:"metric_reset_after_every_n_sec" yaml:"metric_reset_after_every_n_sec"`
}

func (m *MySQLConfig) SetupDefaults() {
	if util.IsStringEmpty(m.Host) {
		m.Host = "localhost"
	}
	if m.Port <= 0 {
		m.Port = 3306
	}
	if util.IsStringEmpty(m.User) {
		m.User = "test"
	}
	if util.IsStringEmpty(m.Password) {
		m.Password = "test"
	}
	if util.IsStringEmpty(m.Db) {
		m.Db = "conversation"
	}
}
