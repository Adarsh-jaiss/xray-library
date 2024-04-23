package logger

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/thesaas-company/xray/types"
)

type Logger struct {
	logs types.ISQL
}

func NewLogger(logs types.ISQL) *Logger {
	return &Logger{
		logs: logs,
	}
}

func (l *Logger) Schema(table string) (types.Table, error) {
	defer func(start time.Time) {
		// Log the execution time
		logrus.WithFields(logrus.Fields{
			"table_name":           table,
			"Query_Execution_time": time.Since(start),
		}).Info("Schema retrieval completed")
	}(time.Now())

	result, err := l.logs.Schema(table)
	if err != nil {
		// Log the error
		logrus.WithFields(logrus.Fields{
			"table_name": table,
			"error":      err.Error(),
		}).Error("Schema retrieval failed")
	}

	return result, err
}

func (l *Logger) Execute(query string) ([]byte, error) {

	defer func(start time.Time) {
		// Log the execution time
		logrus.WithFields(logrus.Fields{
			"query":                query,
			"Query_Execution_time": time.Since(start),
		}).Info("Query execution completed")
	}(time.Now())

	result, err := l.logs.Execute(query)
	if err != nil {
		// Log the error
		logrus.WithFields(logrus.Fields{
			"query": query,
			"error": err.Error(),
		}).Error("Query execution failed")
	}

	return result, err
}

func (l *Logger) Tables(DatabaseName string) ([]string, error) {
	defer func(start time.Time) {
		// Log the execution time
		logrus.WithFields(logrus.Fields{
			"Database_Name":        DatabaseName,
			"Query_Execution_time": time.Since(start),
		}).Info("Tables retrieval completed")
	}(time.Now())

	result, err := l.logs.Tables(DatabaseName)
	if err != nil {
		// Log the error
		logrus.WithFields(logrus.Fields{
			"Database_Name": DatabaseName,
			"error":         err.Error(),
		}).Error("Tables retrieval failed")
	}

	return result, err
}
