package dbx

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	fmtLogQuery        = `Query:          %s`
	fmtLogArgs         = `Arguments:      %#v`
	fmtLogRowsAffected = `Rows affected:  %d`
	fmtLogLastInsertID = `Last insert ID: %d`
	fmtLogError        = `Error:          %v`
	fmtLogTimeTaken    = `Time taken:     %0.5fs`
	fmtLogContext      = `Context:        %v`
)

var (
	reInvisibleChars       = regexp.MustCompile(`[\s\r\n\t]+`)
	reColumnCompareExclude = regexp.MustCompile(`[^a-zA-Z0-9]`)
)

// QueryStatus represents the status of a query after being executed.
type QueryStatus struct {
	Message      *string
	RowsAffected *int64
	LastInsertID *int64

	Query string
	Args  []interface{}

	Err error

	Start time.Time
	End   time.Time

	Context context.Context
}

// String returns a formatted log message.
func (q *QueryStatus) String() string {
	lines := make([]string, 0, 8)
	if q.Message != nil {
		line := *q.Message
		if q.Err != nil {
			line += "err: " + q.Err.Error()
		}
		return line
	}

	if query := q.Query; query != "" {
		query = reInvisibleChars.ReplaceAllString(query, ` `)
		query = strings.TrimSpace(query)
		lines = append(lines, fmt.Sprintf(fmtLogQuery, query))
	}

	if len(q.Args) > 0 {
		lines = append(lines, fmt.Sprintf(fmtLogArgs, q.Args))
	}

	if q.RowsAffected != nil {
		lines = append(lines, fmt.Sprintf(fmtLogRowsAffected, *q.RowsAffected))
	}
	if q.LastInsertID != nil {
		lines = append(lines, fmt.Sprintf(fmtLogLastInsertID, *q.LastInsertID))
	}

	if q.Err != nil {
		lines = append(lines, fmt.Sprintf(fmtLogError, q.Err))
	}

	lines = append(lines, fmt.Sprintf(fmtLogTimeTaken, float64(q.End.UnixNano()-q.Start.UnixNano())/float64(1e9)))

	if q.Context != nil {
		lines = append(lines, fmt.Sprintf(fmtLogContext, q.Context))
	}

	return strings.Join(lines, "\n")
}

type ILogger interface {
	Log(*QueryStatus)
}

// Settings defines methods to get or set configuration values.
type LoggerSettings interface {
	// SetLogging enables or disables logging.
	SetLogging(bool)
	// LoggingEnabled returns true if logging is enabled, false otherwise.
	LoggingEnabled() bool

	// SetLogger defines which logger to use.
	SetLogger(ILogger)
	// Returns the currently configured logger.
	Logger() ILogger
}

type defaultLogger struct {
	loggerSettings
}

func (lg *defaultLogger) Log(m *QueryStatus) {
	if m.Err != nil {
		s := fmt.Sprintf("\033[31m%s \033[0m  \n\t%s\n\n", "ERROR:", strings.Replace(m.String(), "\n", "\n\t", -1))
		fmt.Print(s)
	} else {
		fmt.Printf("\n\t%s\n\n", strings.Replace(m.String(), "\n", "\n\t", -1))
	}
}

var _ = ILogger(&defaultLogger{})

var defaultLoggerSettings = &loggerSettings{}

type loggerSettings struct {
	loggingEnabled bool
	queryLogger    ILogger
}

func (c *loggerSettings) Logger() ILogger {

	if c.queryLogger == nil {
		c.queryLogger = &defaultLogger{}
	}
	return c.queryLogger
}

func (c *loggerSettings) SetLogger(lg ILogger) {
	c.queryLogger = lg
}
func (c *loggerSettings) SetLogging(value bool) {
	c.loggingEnabled = value
}

func (c *loggerSettings) LoggingEnabled() bool {
	return c.loggingEnabled
}
