package logging

// Level is the logger level
type Level int

const (
	// SQL level is the lowest logger level. It dumps Debug level + SQL queries.
	SQL Level = iota
	// Debug level dumps debug log traces and info logs.
	Debug
	// Info level dumps info logs and warnings.
	Info
	// Warn level dumps warnings.
	Warn
)

func (l Level) String() string {
	switch l {
	case SQL:
		return "sql"
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Warn:
		return "warn"
	}
	return "unknown"
}
