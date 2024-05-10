// Package logger provides a logging wrapper for the ISQL interface.
package logger

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/thesaas-company/xray/types"
)

// Logger is a struct that implements the ISQL interface and adds logging functionality.
type Logger struct {
	logs types.ISQL // The underlying ISQL interface for database operations.
}

// NewLogger creates a new Logger instance with the provided ISQL implementation.
func NewLogger(logs types.ISQL) *Logger {
	return &Logger{
		logs: logs,
	}
}

// Schema retrieves the schema for the specified table.
// It logs the execution time and any errors that occur during the retrieval process.
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

// Execute executes the given SQL query.
// It logs the execution time and any errors that occur during the execution process.
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

// Tables retrieves the list of tables for the specified database.
// It logs the execution time and any errors that occur during the retrieval process.
func (l *Logger) Tables(databaseName string) ([]string, error) {
	defer func(start time.Time) {
		// Log the execution time
		logrus.WithFields(logrus.Fields{
			"Database_Name":        databaseName,
			"Query_Execution_time": time.Since(start),
		}).Info("Tables retrieval completed")
	}(time.Now())

	result, err := l.logs.Tables(databaseName)
	if err != nil {
		// Log the error
		logrus.WithFields(logrus.Fields{
			"Database_Name": databaseName,
			"error":         err.Error(),
		}).Error("Tables retrieval failed")
	}

	return result, err
}

func (l *Logger) GenerateCreateTableQuery(table types.Table) string {
	defer func(start time.Time) {
		// Log the execution time
		logrus.WithFields(logrus.Fields{
			"table_name":           table.Name,
			"Query_Execution_time": time.Since(start),
		}).Info("Create table query generation completed")
	}(time.Now())

	result := l.logs.GenerateCreateTableQuery(table)
	return result
}
