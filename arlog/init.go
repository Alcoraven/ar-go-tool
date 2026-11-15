package arlog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type manager struct {
	sync.RWMutex
	mode     LogMode
	plan     LogPlan
	savePath string
	depth    int
	timeLoc  *time.Location
	f        *arLogFs

	showPrefix bool
}

func (m *manager) createFile(f *arLogFs, savePath string) error {
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return err
	}
	var name, schedule string
	switch m.plan {
	case LogPlanOnce:
		return f.openFile(filepath.Join(savePath, "out.log"))
	case LogPlanDay:
		name, schedule = "20060102", "0 0 * * *"
	case LogPlanMonth:
		name, schedule = "200601", "0 0 1 * *"
	default:
		return fmt.Errorf("unknown log plan: %d", m.plan)
	}
	open := func() error {
		return f.openFile(filepath.Join(savePath, time.Now().In(m.timeLoc).Format(name)+".log"))
	}
	if err := open(); err != nil {
		return err
	}
	c := cron.New(cron.WithLocation(m.timeLoc))
	if _, err := c.AddFunc(schedule, func() {
		if err := open(); err != nil {
			fmt.Fprintln(os.Stderr, "open log file error:", err)
		}
	}); err != nil {
		return err
	}
	c.Start()
	return nil
}

func (m *manager) writeLog(level LogLevel, v ...interface{}) {
	var (
		prefix  string
		msg     string
		msgList []string
	)

	if m.showPrefix {
		prefix = getPrefix(m.timeLoc, m.depth, level) + " "
	}
	msg = formatMsg(v...)

	for _, s := range strings.Split(msg, "\n") {
		msgList = append(msgList, fmt.Sprintf("%s%s", prefix, s))
	}

	if m.mode != LogModeProd {
		printStd(level, msgList...)
	}
	m.f.write(msgList...)
}

type PrefixConfig struct {
	ShowPrefix bool
}

type Config struct {
	Mode     LogMode        // log mode
	Plan     LogPlan        // log plan
	Depth    int            // log file depth
	SavePath string         // log save path
	TimeLoc  *time.Location // time location (default for server local location)
	*PrefixConfig
}

func Init(config *Config) (Log, error) {
	var (
		logMode  = LogModeDev
		logPlan  = LogPlanOnce
		savePath = DefaultSavePath
		depth    = 0
		timeLoc  = time.Local

		showPrefix = true
	)

	if config != nil {
		if config.Mode != LogModeDev {
			logMode = config.Mode
		}
		if config.Plan != LogPlanOnce {
			logPlan = config.Plan
		}
		if len(config.SavePath) > 0 {
			savePath = config.SavePath
			if savePath[len(savePath)-1] != '/' {
				savePath += "/"
			}
		}
		if config.Depth > 0 {
			depth = config.Depth
		}
		if config.TimeLoc != nil {
			timeLoc = config.TimeLoc
		}
		if config.PrefixConfig != nil {
			showPrefix = config.ShowPrefix
		}
	}

	var m = &manager{
		mode:       logMode,
		plan:       logPlan,
		depth:      depth,
		timeLoc:    timeLoc,
		f:          new(arLogFs),
		showPrefix: showPrefix,
	}
	if err := m.createFile(m.f, savePath); err != nil {
		return nil, err
	}

	return m, nil
}

type Log interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})

	ChangeMode(mode LogMode)
}

func (m *manager) Debug(v ...interface{}) {
	if m.mode != LogModeDebug {
		return
	}
	m.writeLog(LogLevelDebug, v...)
}

func (m *manager) Info(v ...interface{}) {
	m.writeLog(LogLevelInfo, v...)
}

func (m *manager) Warn(v ...interface{}) {
	m.writeLog(LogLevelWarn, v...)
}

func (m *manager) Error(v ...interface{}) {
	m.writeLog(LogLevelError, v...)
}

func (m *manager) Fatal(v ...interface{}) {
	m.writeLog(LogLevelFatal, v...)
}

func (m *manager) ChangeMode(mode LogMode) {
	m.Lock()
	defer m.Unlock()

	m.mode = mode
}
