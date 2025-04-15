package arlog

const (
	DefaultSavePath = "./logs/"
)

type LogMode int

const (
	LogModeDev LogMode = iota
	LogModeDebug
	LogModeProd
)

type LogPlan int

const (
	LogPlanOnce LogPlan = iota
	LogPlanDay
	LogPlanMonth
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)
